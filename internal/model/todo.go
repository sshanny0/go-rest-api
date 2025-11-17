package model

import (
	"time"

	"gorm.io/gorm"
)

type Todo struct {
	ID          uint   `gorm:"primaryKey"`
	Title       string `gorm:"size:200;not null"`
	Description string `gorm:"type:text"`
	Status      string `gorm:"type:varchar(20);default:'pending';index"`
	Priority    string `gorm:"type:varchar(10);default:'medium'"`
	DueDate     *time.Time
	User        User `gorm:"not null;index"`
	UserID      uint `gorm:"foreignKey:UserID"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

func (Todo) TableName() string {
	return "todos"
}
