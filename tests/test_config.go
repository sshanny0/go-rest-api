package main

import (
	"fmt"
	"rest-api/internal/config"
)

func main() {
	cfg := config.LoadConfig()

	fmt.Printf("Load Config\n")
	fmt.Printf("DB Host: %s\n", cfg.DBHost)
	fmt.Printf("  DB Port: %s\n", cfg.DBPort)
	fmt.Printf("  DB Name: %s\n", cfg.DBName)
	fmt.Printf("  Server Port: %s\n", cfg.ServerPort)
	fmt.Printf("  Gin Mode: %s\n", cfg.GinMode)
}
