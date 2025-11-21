package main

import (
	"encoding/json"
	"fmt"
	"rest-api/internal/dto"
)

func main() {
	// Test RegisterRequest
	reg := dto.RegisterRequest{
		Username: "TEST",
		Email:    "test@example.com",
		Password: "password123",
		FullName: "Test User",
	}

	// Test CreateTodoRequest
	todo := dto.CreateTodoRequest{
		Title:    "Test Todo",
		Status:   "pending",
		Priority: "high",
	}

	// Convert to JSON
	regJSON, _ := json.MarshalIndent(reg, "", "  ")
	todoJSON, _ := json.MarshalIndent(todo, "", "  ")

	fmt.Println("RegisterRequest:")
	fmt.Println(string(regJSON))
	fmt.Println("\nCreateTodoRequest:")
	fmt.Println(string(todoJSON))
}
