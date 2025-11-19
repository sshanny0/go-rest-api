package route

import (
	"rest-api/internal/handler"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(
	router *gin.Engine,
	userHandler *handler.UserHandler,
	healthHandler *handler.HealthHandler,
	todoHandler *handler.TodoHandler,
) {
	// Check health
	router.GET("/health", healthHandler.HealthCheck)

	//API V1 Group
	v1 := router.Group("/api/v1")
	{
		// Auth Routes (Public)
		auth := v1.Group("/auth")
		{
			auth.POST("/register", userHandler.Register)
			auth.POST("/login", userHandler.Login)
			auth.POST("/reset-password", userHandler.ResetPasswordRequest)
			auth.POST("/reset-password/confirm", userHandler.ResetPasswordConfirm)
		}

		// user routes (protected)
		users := v1.Group("/users")
		{
			users.GET("/profile", userHandler.GetProfile)
			users.PUT("profile", userHandler.UpdateProfile)
		}

		// Todo
		todos := v1.Group("todos")
		{
			todos.GET("", todoHandler.GetAll)
			todos.GET("/:id", todoHandler.GetByID)
			todos.POST("", todoHandler.Create)
			todos.PUT(":id", todoHandler.Update)
			todos.DELETE(":id", todoHandler.Delete)
		}
	}
}
