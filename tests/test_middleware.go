package main

import (
	"fmt"
	"rest-api/internal/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Apply middleware
	router.Use(middleware.LoggerMiddleware())
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.ErrorHandler())

	// Test route
	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Middleware loaded successfully"})
	})

	fmt.Println("✅ Middleware functions loaded successfully")
	fmt.Println("✅ Router configured with middleware chain")
	fmt.Println("✅ Ready for testing with actual handlers")
}
