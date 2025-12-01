package models

import (
	"time"

	"gorm.io/gorm"
)

type Room struct {
	ID        string         `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	Members   []User         `json:"members" gorm:"many2many:room_members"`
}

type RoomCheckResult struct {
	Exists bool `json:"exists"`
	RoomID any  `json:"room_id,omitempty"`
}

type RoomCreateResult struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
	Room    *Room  `json:"room,omitempty"`
}
