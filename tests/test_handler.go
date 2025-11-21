package main

import (
	"fmt"
	"net/http/httptest"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	"rest-api/internal/handler"
	"rest-api/internal/model"
	"rest-api/internal/repository"
	"rest-api/internal/service"
)

func main() {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Setup in-memory DB and migrate
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	if err := db.AutoMigrate(&model.User{}, &model.Todo{}); err != nil {
		panic(err)
	}

	// Wire repository -> service -> handler
	userRepo := repository.NewUserRepository(db)
	authService := service.NewAuthService(userRepo)
	userHandler := handler.NewUserHandler(authService)

	router.POST("/register", userHandler.Register)

	// Create request
	reqBody := `{
		"username": "testuser",
		"email": "test@example.com",
		"password": "password123",
		"full_name": "Test User"
	}`
	req := httptest.NewRequest("POST", "/register", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	// Record response
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Check
	if w.Code != 201 {
		panic(fmt.Sprintf("expected 201, got %d - body: %s", w.Code, w.Body.String()))
	}
	if !strings.Contains(w.Body.String(), "testuser") {
		panic("response does not contain username")
	}

	fmt.Println("✅ User handler registration route working correctly")
}
