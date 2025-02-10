package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/K4rian/kfrs/internal/log"
)

type KFHTTPRedirectServer struct {
	host             string
	port             int
	rootDir          string
	server           *http.Server
	maxRequestsPerIP int
	banTime          int
	ctx              context.Context
	cancel           context.CancelFunc
	mu               sync.Mutex
	ipRequestCount   map[string][]time.Time
	ipBlockList      map[string]time.Time
	defaultFile      string
}

func NewKFHTTPRedirectServer(
	host string,
	port int,
	rootDir string,
	maxRequestsPerIP int,
	banTime int,
	ctx context.Context,
) *KFHTTPRedirectServer {
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir(rootDir))
	server := &KFHTTPRedirectServer{
		server: &http.Server{
			Addr:    fmt.Sprintf("%s:%d", host, port),
			Handler: mux,
		},
		host:             host,
		port:             port,
		rootDir:          rootDir,
		maxRequestsPerIP: maxRequestsPerIP,
		banTime:          banTime,
		ipRequestCount:   make(map[string][]time.Time),
		ipBlockList:      make(map[string]time.Time),
		defaultFile:      filepath.Join(rootDir, "index.html"),
	}
	server.ctx, server.cancel = context.WithCancel(ctx)

	mux.HandleFunc("/", server.handleRequest(fileServer))
	return server
}

func (h *KFHTTPRedirectServer) Address() string {
	return fmt.Sprintf("%s:%d", h.host, h.port)
}

func (h *KFHTTPRedirectServer) RootDirectory() string {
	return h.rootDir
}

func (h *KFHTTPRedirectServer) Listen() error {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.server == nil {
		return fmt.Errorf("already stopped")
	}

	done := make(chan error, 1)

	go func() {
		log.Logger.Debug("Server starting to listen...", "address", h.Address())
		if err := h.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			done <- err
		}
	}()

	// Don't wait more than 1 second
	select {
	case err := <-done:
		return err
	case <-time.After(1 * time.Second):
		return nil
	}
}

func (h *KFHTTPRedirectServer) Stop() error {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.cancel != nil {
		log.Logger.Debug("Cancelling server context")
		h.cancel()
	}

	shutdownCtx, cancelShutdown := context.WithTimeout(h.ctx, 5*time.Second)
	defer cancelShutdown()

	log.Logger.Debug("Shutting down the server...", "timeout", "5 seconds")
	if err := h.server.Shutdown(shutdownCtx); err != nil {
		return err
	}
	log.Logger.Debug("Server shutdown complete")
	return nil
}

func (h *KFHTTPRedirectServer) handleRequest(fileServer http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		log.Logger.Debug("Request received", "ip", ip, "method", r.Method, "path", r.URL.Path)

		// Check if the IP is blocked
		if h.isIPBlocked(ip) {
			log.Logger.Warn("IP blocked", "ip", ip)
			http.Error(w, ErrForbidden, http.StatusForbidden)
			return
		}

		// Track requests for the IP
		h.trackIPRequest(ip)

		// Only GET method is allowed
		if r.Method != http.MethodGet {
			log.Logger.Warn("Method not allowed", "method", r.Method, "path", r.URL.Path)
			sendHTTPError(w, http.StatusMethodNotAllowed, ErrMethodNotAllowed)
			return
		}

		// Resolve and clean the requested file path
		reqPath := r.URL.Path
		filePath := filepath.Join(h.rootDir, filepath.Clean(reqPath))

		log.Logger.Debug("Resolving file path", "requestedPath", reqPath, "resolvedPath", filePath)

		// Ensure the file path is within the root directory
		rootDirAbs, err := filepath.Abs(h.rootDir)
		if err != nil {
			log.Logger.Error("Failed to resolve root directory", "error", err)
			sendHTTPError(w, http.StatusInternalServerError, fmt.Sprintf("failed to resolve root directory: %v", err))
			return
		}

		absFilePath, err := filepath.Abs(filePath)
		if err != nil {
			log.Logger.Error("Failed to resolve file path", "error", err)
			sendHTTPError(w, http.StatusInternalServerError, fmt.Sprintf("failed to resolve file path: %v", err))
			return
		}

		rel, err := filepath.Rel(rootDirAbs, absFilePath)
		if err != nil || strings.HasPrefix(rel, "..") {
			log.Logger.Warn("File path outside root directory", "requestedPath", reqPath)
			sendHTTPError(w, http.StatusForbidden, ErrForbidden)
			return
		}

		// If the path is "/" or empty, try to serve the index.html
		if reqPath == "/" || reqPath == "" {
			if _, err := os.Stat(h.defaultFile); err == nil {
				log.Logger.Debug("Serving default index file", "file", h.defaultFile)
				http.ServeFile(w, r, h.defaultFile)
				return
			}
			log.Logger.Error("Index.html not found", "error", "file not found")
			sendHTTPError(w, http.StatusInternalServerError, ErrInternalServerError)
			return
		}

		// Check that requested file ends with .uz2
		if !strings.HasSuffix(reqPath, ".uz2") {
			log.Logger.Warn("Forbidden file extension", "requestedPath", reqPath)
			sendHTTPError(w, http.StatusForbidden, ErrForbidden)
			return
		}

		// Check if the file exists
		info, err := os.Stat(absFilePath)
		if os.IsNotExist(err) {
			log.Logger.Warn("File not found", "requestedPath", reqPath)
			sendHTTPError(w, http.StatusNotFound, ErrNotFound)
			return
		} else if err != nil {
			log.Logger.Error("Error checking file status", "error", err)
			sendHTTPError(w, http.StatusInternalServerError, ErrInternalServerError)
			return
		} else if info.IsDir() {
			log.Logger.Warn("Requested path is a directory", "requestedPath", reqPath)
			sendHTTPError(w, http.StatusForbidden, ErrForbidden)
			return
		}
		log.Logger.Debug("Serving file", "filePath", absFilePath)
		fileServer.ServeHTTP(w, r)
	}
}

func (h *KFHTTPRedirectServer) trackIPRequest(ip string) {
	now := time.Now()

	// Limit the tracked request times to the last 60 seconds
	h.mu.Lock()
	defer h.mu.Unlock()

	// If this IP has been blocked, do nothing
	if unblockTime, ok := h.ipBlockList[ip]; ok && time.Now().Before(unblockTime) {
		log.Logger.Debug("IP still blocked", "ip", ip, "unblockTime", unblockTime)
		return
	}

	// Clean up old request times (older than 60 seconds)
	h.ipRequestCount[ip] = filterRecentRequests(h.ipRequestCount[ip], now)

	// Check if the IP has exceeded the request threshold
	if len(h.ipRequestCount[ip]) > h.maxRequestsPerIP {
		log.Logger.Warn("IP exceeded request threshold, blocking", "ip", ip)
		h.ipBlockList[ip] = now.Add(time.Duration(h.banTime) * time.Minute)
	}

	// Record the request time
	h.ipRequestCount[ip] = append(h.ipRequestCount[ip], now)
	log.Logger.Debug("Tracking IP request", "ip", ip, "requestCount", len(h.ipRequestCount[ip]))
}

func (h *KFHTTPRedirectServer) isIPBlocked(ip string) bool {
	h.mu.Lock()
	defer h.mu.Unlock()

	// Check if the IP is in the blocklist and if the block period has not expired
	unblockTime, blocked := h.ipBlockList[ip]
	return blocked && time.Now().Before(unblockTime)
}
