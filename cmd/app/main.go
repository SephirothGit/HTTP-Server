package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/SephirothGit/Backend-service/internal/handler"
	"github.com/SephirothGit/Backend-service/internal/repository"
	"github.com/SephirothGit/Backend-service/internal/server"
	"github.com/SephirothGit/Backend-service/internal/service"
)

func main() {
	repo := repository.NewInMemoryOrderRepository()
	svc := service.NewOrderService(repo)
	orderHandler := handler.NewOrderHandler(svc)

	router := server.NewRouter(server.RouterDeps{
		OrderHandler: orderHandler,
	})
	wrapped := server.LoggingMiddleware(router)

	srv := server.NewServer(":8080", wrapped)

	go func() {
		if err := srv.Start(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed %v", err)
		}
	}()

	log.Println("Server started on :8080")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed %v", err)
	}

	log.Println("Server stopped")
}