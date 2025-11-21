package main

import (
	"fmt"
	"rest-api/internal/utils"
)

func main() {
	password := "mySecretPassword123"

	// Hash password
	fmt.Println("Original password:", password)
	hash, err := utils.HashPassword(password)
	if err != nil {
		fmt.Println("Error hashing:", err)
		return
	}
	fmt.Println("Hashed password:", hash)

	// Verify correct password
	if utils.CheckPassword(password, hash) {
		fmt.Println("✅ Password verification: SUCCESS")
	} else {
		fmt.Println("❌ Password verification: FAILED")
	}

	// Verify wrong password
	if utils.CheckPassword("wrongPassword", hash) {
		fmt.Println("❌ Wrong password accepted (BAD!)")
	} else {
		fmt.Println("✅ Wrong password rejected: SUCCESS")
	}

	// Test that same password gives different hash
	hash2, _ := utils.HashPassword(password)
	fmt.Println("\nSecond hash:", hash2)
	fmt.Println("Hashes different?", hash != hash2)
	fmt.Println("Both valid?", utils.CheckPassword(password, hash2))
}
