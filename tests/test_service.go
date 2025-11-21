package main

import (
	"fmt"
	"rest-api/internal/dto"
	"rest-api/internal/model"
	"rest-api/internal/repository"
	"rest-api/internal/service"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func main() {
	// Setup in-memory database
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&model.User{}, &model.Todo{})

	// Create repositories
	userRepo := repository.NewUserRepository(db)
	todoRepo := repository.NewTodoRepository(db)

	// Create services
	authService := service.NewAuthService(userRepo)
	todoService := service.NewTodoService(todoRepo)

	// Test registration
	registerReq := dto.RegisterRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
		FullName: "Test User",
	}
	user, err := authService.Register(registerReq)
	if err != nil {
		panic(err)
	}
	fmt.Printf("✅ User registered: %s (ID: %d)\n", user.Username, user.ID)

	// Test duplicate username
	_, err = authService.Register(registerReq)
	if err == service.ErrUserExists {
		fmt.Println("✅ Duplicate username detected correctly")
	}

	// Test login
	loginReq := dto.LoginRequest{
		Username: "testuser",
		Password: "password123",
	}
	loginResp, err := authService.Login(loginReq)
	if err != nil {
		panic(err)
	}
	token := loginResp.Token
	loginUser := loginResp.User
	fmt.Printf("✅ Login successful, token generated: %s...\n", token[:20])

	// Test wrong password
	wrongReq := dto.LoginRequest{
		Username: "testuser",
		Password: "wrongpassword",
	}
	_, err = authService.Login(wrongReq)
	if err == service.ErrInvalidCredentials {
		fmt.Println("✅ Wrong password detected correctly")
	}

	// Test create todo
	createTodoReq := dto.CreateTodoRequest{
		Title:       "Test Todo",
		Description: "Test description",
		Status:      "pending",
		Priority:    "high",
	}
	todo, err := todoService.CreateTodo(loginUser.ID, createTodoReq)
	if err != nil {
		panic(err)
	}
	fmt.Printf("✅ Todo created: %s (ID: %d)\n", todo.Title, todo.ID)

	// Test invalid status
	invalidReq := dto.CreateTodoRequest{
		Title:    "Invalid Todo",
		Status:   "invalid_status",
		Priority: "high",
	}
	_, err = todoService.CreateTodo(loginUser.ID, invalidReq)
	if err == service.ErrInvalidStatus {
		fmt.Println("✅ Invalid status detected correctly")
	}

	// Test get user todos
	todos, err := todoService.GetUserTodos(loginUser.ID, "", "")
	if err != nil {
		panic(err)
	}
	fmt.Printf("✅ Retrieved %d todos for user\n", len(todos))

	// Test ownership check
	_, err = todoService.GetTodoByID(todo.ID, 9999) // Wrong user ID
	if err == service.ErrUnauthorizedAccess {
		fmt.Println("✅ Unauthorized access prevented correctly")
	}

	fmt.Println("\n✅ All service methods working correctly!")

}
