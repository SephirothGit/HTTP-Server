package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"fmt"
	"github.com/golang-jwt/jwt/v5"

	"github.com/SephirothGit/Backend-service/internal/handler"
	"github.com/SephirothGit/Backend-service/internal/repository"
	"github.com/SephirothGit/Backend-service/internal/server"
	"github.com/SephirothGit/Backend-service/internal/service"
)

func main() {
	
	if err := server.InitLogger(); err != nil {
		log.Fatalf("failed to init logger: %v", err)
	}
	defer server.Log.Sync()

	repo := repository.NewInMemoryOrderRepository()
	svc := service.NewOrderService(repo)
	orderHandler := handler.NewOrderHandler(svc)

	router := server.NewRouter(server.RouterDeps{
		OrderHandler: orderHandler,
	})
	wrapped := server.LoggingMiddleware(server.TimeoutMiddleware(5 * time.Second)(router))
	wrapped = server.Timeout504Middleware(wrapped)

	srv := server.NewServer(server.Config{
		Addr: ":8080",
		ReadTimeout: 5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout: 60 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
	}, wrapped)


	go func() {
		log.Println("Server started on :8080")
		if err := srv.Start(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown signal received")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed %v", err)
	}

	// JWT
	func generateToken() {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"user_id": "123",
			"exp": time.Now.Add(time.Hour).Unix(),
		})

		tokenStr, _ := token.SignedString([]byte("supersecretkey"))
		fmt.Println(tokenStr)
	}

	log.Println("Server stopped gracefully")
}