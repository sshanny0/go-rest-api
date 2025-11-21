package main

import (
	"fmt"
	"rest-api/internal/handler"
	"rest-api/internal/middleware"
	"rest-api/internal/route"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Apply global middleware
	router.Use(middleware.LoggerMiddleware())
	router.Use(middleware.CORSMiddleware())

	// Create dummy handlers (won't work without real services)
	userHandler := &handler.UserHandler{}
	todoHandler := &handler.TodoHandler{}
	healthHandler := &handler.HealthHandler{}

	// Setup routes
	route.SetupRoutes(router, userHandler, healthHandler, todoHandler)

	// List all routes
	fmt.Println("📍 Registered Routes:")
	fmt.Println("====================")
	routes := router.Routes()
	for _, r := range routes {
		fmt.Printf("%-7s %s\n", r.Method, r.Path)
	}

	fmt.Println("\n✅ Route configuration loaded successfully!")
	fmt.Printf("📊 Total routes: %d\n", len(routes))

}
