package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID       uint   `gorm:"primaryKey"`
	Username string `gorm:"unique;not null;size:50;index"`
	Email    string `gorm:"unique;not null;size:50;index"`
	Password string `gorm:"not null"`
	Fullname string `gorm:"size:100"`
	// Reset password token and expiry
	ResetPasswordToken  string     `gorm:"size:255;index"`
	ResetPasswordExpiry *time.Time `gorm:"index"`
	Todos               []Todo     `gorm:"foreignKey:UserID"`
	CreatedAt           time.Time
	UpdatedAt           time.Time
	DeletedaT           gorm.DeletedAt `gorm:"index"`
}

func (User) TableName() string {
	return "users"
}
