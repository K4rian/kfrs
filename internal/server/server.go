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
)

type KFHTTPRedirectServer struct {
	ip               string
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
	ip string,
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
			Addr:    fmt.Sprintf("%s:%d", ip, port),
			Handler: mux,
		},
		ip:               ip,
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

func (h *KFHTTPRedirectServer) Host() string {
	return fmt.Sprintf("%s:%d", h.ip, h.port)
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
		h.cancel()
	}

	shutdownCtx, cancelShutdown := context.WithTimeout(h.ctx, 5*time.Second)
	defer cancelShutdown()
	return h.server.Shutdown(shutdownCtx)
}

func (h *KFHTTPRedirectServer) handleRequest(fileServer http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr

		// Check if the IP is blocked
		if h.isIPBlocked(ip) {
			http.Error(w, ErrForbidden, http.StatusForbidden)
			return
		}

		// Track requests for the IP
		h.trackIPRequest(ip)

		// Only GET method is allowed
		if r.Method != http.MethodGet {
			sendHTTPError(w, http.StatusMethodNotAllowed, ErrMethodNotAllowed)
			return
		}

		// Resolve and clean the requested file path
		reqPath := r.URL.Path
		filePath := filepath.Join(h.rootDir, filepath.Clean(reqPath))

		// Ensure the file path is within the root directory
		rootDirAbs, err := filepath.Abs(h.rootDir)
		if err != nil {
			sendHTTPError(w, http.StatusInternalServerError, fmt.Sprintf("failed to resolve root directory: %v", err))
			return
		}

		absFilePath, err := filepath.Abs(filePath)
		if err != nil {
			sendHTTPError(w, http.StatusInternalServerError, fmt.Sprintf("failed to resolve file path: %v", err))
			return
		}

		rel, err := filepath.Rel(rootDirAbs, absFilePath)
		if err != nil || strings.HasPrefix(rel, "..") {
			// If the file path is outside of the root directory, send 403
			sendHTTPError(w, http.StatusForbidden, ErrForbidden)
			return
		}

		// If the path is "/" or empty, try to serve the index.html
		if reqPath == "/" || reqPath == "" {
			if _, err := os.Stat(h.defaultFile); err == nil {
				http.ServeFile(w, r, h.defaultFile)
				return
			}
			// Return 500 if index.html doesn't exist
			sendHTTPError(w, http.StatusInternalServerError, ErrInternalServerError)
			return
		}

		// Check that requested file ends with .uz2
		if !strings.HasSuffix(reqPath, ".uz2") {
			sendHTTPError(w, http.StatusForbidden, ErrForbidden)
			return
		}

		// Check if the file exists
		info, err := os.Stat(absFilePath)
		if os.IsNotExist(err) {
			sendHTTPError(w, http.StatusNotFound, ErrNotFound)
			return
		} else if err != nil {
			sendHTTPError(w, http.StatusInternalServerError, ErrInternalServerError)
			return
		} else if info.IsDir() {
			sendHTTPError(w, http.StatusForbidden, ErrForbidden)
			return
		}
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
		return
	}

	// Clean up old request times (older than 60 seconds)
	h.ipRequestCount[ip] = filterRecentRequests(h.ipRequestCount[ip], now)

	// Check if the IP has exceeded the request threshold
	if len(h.ipRequestCount[ip]) > h.maxRequestsPerIP {
		// Block the IP
		h.ipBlockList[ip] = now.Add(time.Duration(h.banTime) * time.Minute)
	}

	// Record the request time
	h.ipRequestCount[ip] = append(h.ipRequestCount[ip], now)
}

func (h *KFHTTPRedirectServer) isIPBlocked(ip string) bool {
	h.mu.Lock()
	defer h.mu.Unlock()

	// Check if the IP is in the blocklist and if the block period has not expired
	unblockTime, blocked := h.ipBlockList[ip]
	return blocked && time.Now().Before(unblockTime)
}
