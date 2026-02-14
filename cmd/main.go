package main

import (
	"log"
	"net/http"
	"subscription-service/configs"
	_ "subscription-service/docs"
	"subscription-service/internal/db"
	"subscription-service/internal/handlers"
	"subscription-service/internal/logger"
	"subscription-service/internal/middleware"
	"subscription-service/internal/repository"
	"subscription-service/internal/services"
	"subscription-service/pkg/health"
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
	defer pool.Close()

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
	log.Printf("Sever is started on port %s", conf.App.Port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
