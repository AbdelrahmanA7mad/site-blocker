package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Initialize configuration
	if err := initHostsFile(); err != nil {
		log.Fatalf("Failed to initialize configuration: %v", err)
	}

	log.Println("Focus Proxy Starting...")

	server := &http.Server{
		Addr:    ":8080",
		Handler: http.HandlerFunc(handleProxyRequest),
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	enableSystemProxy()

	// Listen for interrupt signals in the background to ensure cleanup on abrupt exit
	go func() {
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
		<-stop
		cleanupAndExit(server)
	}()

	// Start the interactive Terminal UI menu
	runTerminalUI()

	// Clean exit when runTerminalUI returns
	cleanupAndExit(server)
}

func cleanupAndExit(server *http.Server) {
	log.Println("\nStopping server and restoring network settings...")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_ = server.Shutdown(ctx)

	disableSystemProxy()
	log.Println("Clean exit. Goodbye!")
	os.Exit(0)
}
