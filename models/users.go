package models

import (
	"encoding/json"
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

type Notification struct {
	ID        uint            `json:"id" gorm:"primarykey"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
	DeletedAt gorm.DeletedAt  `json:"-" gorm:"index"`
	UserID    uint            `json:"user_id" gorm:"not null;index"`
	Message   string          `json:"message" gorm:"size:255;not null"`
	Metadata  json.RawMessage `json:"metadata" gorm:"type:jsonb"`
	Read      bool            `json:"read" gorm:"default:false"`
}
