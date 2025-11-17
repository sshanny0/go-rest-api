package config

import (
	"fmt"
	"log"
	"rest-api/internal/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

func NewDatabase(cfg *Config) (*gorm.DB, error) {
	// inisiasi koneksi
	dsn := fmt.Sprintf("host = %s user = %s password = %s dbname = %s port = %s sslmode = disable, TimeZone = Asia/Jakarta",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort)

	// koneksi database
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	log.Println("Succesfully connected")

	// auto migrate model later
	if err := db.AutoMigrate(&model.User{}, &model.Todo{}); err != nil {
		return nil, fmt.Errorf("failed to migrate the database: %w", err)
	}

	return db, nil
}
