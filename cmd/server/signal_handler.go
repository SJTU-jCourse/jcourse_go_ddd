package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	ShutdownTimeout       = 30 * time.Second
	WorkerCleanupDelay    = 5 * time.Second
	SignalChannelCapacity = 1
)

func SetupSignalHandler() (chan os.Signal, context.Context, context.CancelFunc) {
	sigChan := make(chan os.Signal, SignalChannelCapacity)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	
	ctx, cancel := context.WithCancel(context.Background())
	return sigChan, ctx, cancel
}

func WaitForShutdown(server *Server, sigChan chan os.Signal, cancel context.CancelFunc) {
	// Wait for shutdown signal
	<-sigChan
	log.Println("Received shutdown signal, gracefully stopping server...")

	// Cancel context to stop all background workers
	cancel()

	// Create shutdown context with timeout
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), ShutdownTimeout)
	defer shutdownCancel()

	// Gracefully shutdown the server
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	} else {
		log.Println("Server stopped gracefully")
	}

	// Give workers time to clean up
	time.Sleep(WorkerCleanupDelay)
	log.Println("All services stopped")
}