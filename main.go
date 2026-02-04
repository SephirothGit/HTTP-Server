package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	repo := repository.NewInMemoryOrderRepository()
	svc := service.NewOrderService(repo)
	orderHandler := handler.NewOrderHandler(svc)

	router := server.NewRouter(orderHandler)
	wrapped := server.LoggingMiddleware(router)

	srv := server.NewServer(":8080", wrapped)

	go func() {
		if err != srv.Start(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background, 5 * time.Second)
	defer cancel()

	if err := srv.Shutdown(); err != nil {
		log.Fatalf("Server shutdown failed")
	}

	log.Println("Server exited correctly")
}
