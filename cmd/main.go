package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"subscription-service/configs"
	_ "subscription-service/docs"
	"subscription-service/internal/db"
	"subscription-service/internal/handlers"
	"subscription-service/internal/logger"
	"subscription-service/internal/middleware"
	"subscription-service/internal/repository"
	"subscription-service/internal/services"
	"subscription-service/pkg/health"
	"syscall"
	"time"

	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Susbsription Service API
// @version 1.0
// @description API for managing subscriptions
// @host localhost:8081
// @BasePath /

func main() {
	// Config
	conf := configs.LoadConfig()

	// Logger
	logger.Init(conf.App.LogLevel)

	// DB connection
	pool, err := db.NewPool(conf.Db.Dsn)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	router := http.NewServeMux()
	//Repositories
	subsRepo := repository.NewSubsRepo(pool)
	//Services
	subsService := services.NewSubsService(subsRepo)
	//Handler
	handlers.RegisterSubsRoutes(router, subsService)
	router.HandleFunc("/healt", health.Handler)

	//Middleware

	handler := middleware.Logging(router)

	//Swagger
	router.Handle("/swagger/", httpSwagger.WrapHandler)

	server := &http.Server{
		Addr:         ":" + conf.App.Port,
		Handler:      handler,
		ReadTimeout:  conf.App.ReadTimeout,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	go func() {
		log.Printf("Server started on port %s", conf.App.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("server error: %v", err)
		}
	}()
	//graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	defer pool.Close()
	log.Println("Server exited properly")
}
