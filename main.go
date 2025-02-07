package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/K4rian/kfrs/cmd"
	"github.com/K4rian/kfrs/internal/server"
)

func main() {
	rootCmd := cmd.BuildRootCommand()

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	redirectServer := server.NewKFHTTPRedirectServer(
		cmd.Host,
		cmd.Port,
		cmd.Directory,
		cmd.MaxRequests,
		cmd.BanTime,
		ctx,
	)

	log.Printf("> Starting the HTTP Redirect Server on %s:%d...\n", cmd.Host, cmd.Port)

	if err := redirectServer.Listen(); err != nil {
		log.Fatalf("Failed to start the HTTP Redirect Server: %v", err)
	}

	log.Printf("> HTTP Redirect Server serving '%s' on %s\n", redirectServer.RootDirectory(), redirectServer.Host())

	<-signalChan

	log.Println("\nShutting down...")

	cancel()
	redirectServer.Stop()

	log.Println("The HTTP Redirect Server has been stopped.")
}
