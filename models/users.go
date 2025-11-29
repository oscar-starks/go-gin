package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	Name      string         `json:"name" gorm:"size:100;not null"`
	Email     string         `json:"email" gorm:"size:100;not null;uniqueIndex"`
	Password  string         `json:"-" gorm:"size:255;not null"`
	Age       int            `json:"age"`
}

// type Product struct {
// 	ID          uint           `json:"id" gorm:"primarykey"`
// 	CreatedAt   time.Time      `json:"created_at"`
// 	UpdatedAt   time.Time      `json:"updated_at"`
// 	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
// 	Name        string         `json:"name" gorm:"size:100;not null"`
// 	Description string         `json:"description" gorm:"type:text"`
// 	Price       float64        `json:"price" gorm:"type:decimal(10,2);not null"`
// 	UserID      uint           `json:"user_id"`
// 	User        User           `json:"user" gorm:"foreignKey:UserID"`
// }
