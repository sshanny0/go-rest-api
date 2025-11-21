package main

import (
	"fmt"
	"log"
	"rest-api/internal/config"
	"rest-api/internal/handler"
	"rest-api/internal/middleware"
	"rest-api/internal/repository"
	"rest-api/internal/route"
	"rest-api/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()
	log.Println("Configuration loaded")

	// Set Gin mode
	gin.SetMode(cfg.GinMode)

	// Initialize database connection
	db, _ := config.NewDatabase(cfg)
	log.Println("Database connected")

	// Manual Dependency Injection Pattern
	// Layer 1: Initialize Repositories (Data Access Layer)
	userRepository := repository.NewUserRepository(db)
	todoRepository := repository.NewTodoRepository(db)
	log.Println("Repositories initialized")

	// Layer 2: Initialize Services (Business Logic Layer)
	authService := service.NewAuthService(userRepository)
	todoService := service.NewTodoService(todoRepository)
	log.Println("Services initialized")

	// Layer 3: Initialize Handlers (HTTP Layer)
	userHandler := handler.NewUserHandler(authService)
	todoHandler := handler.NewTodoHandler(todoService)
	healthHandler := handler.NewHealthHandler(db)
	log.Println("Handlers initialized")

	// Initialize Gin router
	router := gin.Default()

	// Apply global middleware
	router.Use(middleware.LoggerMiddleware())
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.ErrorHandler())
	log.Println("Middleware applied")

	// Setup routes
	route.SetupRoutes(router, userHandler, healthHandler, todoHandler)
	log.Println("Routes configured")

	// Start server
	serverAddr := ":" + cfg.ServerPort
	log.Printf("Server starting on %s", serverAddr)
	log.Printf("Environment: %s", cfg.GinMode)
	fmt.Printf("%q\n", cfg.ServerPort)

	if err := router.Run(serverAddr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

}
