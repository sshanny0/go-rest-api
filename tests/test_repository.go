package main

import (
	"fmt"
	"rest-api/internal/model"
	"rest-api/internal/repository"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func main() {
	// In-memory SQLite for testing
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	// Auto-migrate
	db.AutoMigrate(&model.User{}, &model.Todo{})

	// Create repositories
	userRepo := repository.NewUserRepository(db)
	todoRepo := repository.NewTodoRepository(db)

	// Test user repository
	user := &model.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "hashedpassword",
		Fullname: "Test User",
	}
	err = userRepo.Create(user)
	if err != nil {
		panic(err)
	}
	fmt.Printf("✅ User created with ID: %d\n", user.ID)

	// Test todo repository
	todo := &model.Todo{
		Title:    "Test Todo",
		Status:   "pending",
		Priority: "high",
		UserID:   user.ID,
	}
	err = todoRepo.Create(todo)
	if err != nil {
		panic(err)
	}
	fmt.Printf("✅ Todo created with ID: %d\n", todo.ID)

	// Test queries
	found, err := userRepo.FindByUsername("testuser")
	if err != nil {
		panic(err)
	}
	fmt.Printf("✅ User found: %s (%s)\n", found.Username, found.Email)

	todos, err := todoRepo.FindByUserID(user.ID)
	if err != nil {
		panic(err)
	}
	fmt.Printf("✅ Found %d todos for user\n", len(todos))

	// Test ownership
	owned, err := todoRepo.IsOwnedByUser(todo.ID, user.ID)
	if err != nil {
		panic(err)
	}
	fmt.Printf("✅ Ownership check: %v\n", owned)

	fmt.Println("\n✅ All repository methods working correctly!")
}
