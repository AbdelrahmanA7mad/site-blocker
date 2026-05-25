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
	log.Println("Focus Proxy Starting...")

	server := &http.Server{
		Addr:    ":8080",
		Handler: http.HandlerFunc(routeRequest),
	}

	go func() {
		log.Println("Running on", proxyAddr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	enableSystemProxy()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("Stopping...")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_ = server.Shutdown(ctx)

	disableSystemProxy()
	log.Println("Clean exit")
}
