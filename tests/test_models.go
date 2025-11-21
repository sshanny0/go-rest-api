package main

import (
	"fmt"
	"rest-api/internal/model"
)

func main() {
	// CREATE USER
	user := model.User{
		Username: "testuser",
		Email:    "test@mail.com",
		Password: "hashedPassword",
		Fullname: "Testing User",
	}

	todo := model.Todo{
		Title:       "Test Todo",
		Description: "For testing todo purposes",
		Status:      "Pending",
		Priority:    "Medium",
		UserID:      1,
	}

	fmt.Printf("User struct created: %s (%s)!\n", user.Username, user.Email)
	fmt.Printf("Todo struct created: %s\n", todo.Title)
	fmt.Println("Models work properly!")
}
