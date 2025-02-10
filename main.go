package main

import (
	"context"
	"fmt"
	stdlog "log"
	"os"
	"os/signal"
	"syscall"

	"github.com/K4rian/kfrs/cmd"
	"github.com/K4rian/kfrs/internal/config"
	"github.com/K4rian/kfrs/internal/log"
	"github.com/K4rian/kfrs/internal/server"
)

func main() {
	rootCmd := cmd.BuildRootCommand()
	if err := rootCmd.Execute(); err != nil {
		stdlog.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	conf := config.Get()
	server := server.NewKFHTTPRedirectServer(
		conf.Host,
		conf.Port,
		conf.ServeDir,
		conf.MaxRequests,
		conf.BanTime,
		ctx,
	)

	log.Logger.Info(fmt.Sprintf("Starting the HTTP Redirect Server on %s:%d", conf.Host, conf.Port))

	if err := server.Listen(); err != nil {
		log.Logger.Error("Failed to start the HTTP Redirect Server", "error", err)
	}
	log.Logger.Info("HTTP Redirect Server started", "rootDir", server.RootDirectory(), "address", server.Address())

	<-signalChan

	log.Logger.Info("Shutting down...")
	cancel()
	server.Stop()

	log.Logger.Info("The HTTP Redirect Server has been stopped.")
}
