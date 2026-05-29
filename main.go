package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	backgroundMode bool
	startEnabled   bool
	shutdownChan   = make(chan struct{})
)

func init() {
	flag.BoolVar(&backgroundMode, "background", false, "Run in background with no console window")
	flag.BoolVar(&startEnabled, "enabled", false, "Start with proxy enabled")
}

func main() {
	flag.Parse()

	if backgroundMode {
		runBackgroundInstance()
		return
	}

	// Normal interactive mode: Check if already running in background
	conn, err := net.DialTimeout("tcp", "127.0.0.1:8080", 150*time.Millisecond)
	if err == nil {
		conn.Close()
		// Send request to restore (terminate background instance)
		client := &http.Client{Timeout: 1 * time.Second}
		resp, err := client.Get("http://127.0.0.1:8080/restore")
		if err == nil {
			resp.Body.Close()
			fmt.Println("\033[32;1m>>> Restoring Focus Proxy control to this window... <<<\033[0m")
			time.Sleep(600 * time.Millisecond) // Let the port free
		} else {
			log.Fatalf("Port 8080 is already in use by another application.")
		}
	}

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

func runBackgroundInstance() {
	if err := initHostsFile(); err != nil {
		log.Fatalf("Failed to initialize configuration: %v", err)
	}

	server := &http.Server{
		Addr:    ":8080",
		Handler: http.HandlerFunc(handleProxyRequest),
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	if startEnabled {
		enableSystemProxy()
	} else {
		disableSystemProxy()
	}

	// Wait for restore signal
	<-shutdownChan

	// Clean exit
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	_ = server.Shutdown(ctx)
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
