package main

import (
	"context"
	"curse-count/handler"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		<-sigChan
		cancel()
	}()

	if err := handler.Handler(ctx); err != nil {
		log.Printf("Failed to serve: %v\n", err)
	}
}
