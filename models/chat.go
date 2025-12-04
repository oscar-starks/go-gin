package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Room struct {
	ID        string         `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	Members   []User         `json:"members" gorm:"many2many:room_members"`
}

// BeforeCreate will set a UUID rather than numeric ID
func (r *Room) BeforeCreate(tx *gorm.DB) error {
	if r.ID == "" {
		r.ID = uuid.New().String()
	}
	return nil
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

type ChatRequest struct {
	ID         string         `json:"id" gorm:"primarykey"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`
	SenderID   uint           `json:"sender_id" gorm:"not null"`
	ReceiverID uint           `json:"receiver_id" gorm:"not null"`
	Sender     User           `json:"sender" gorm:"foreignKey:SenderID"`
	Receiver   User           `json:"receiver" gorm:"foreignKey:ReceiverID"`
	Status     string         `json:"status" gorm:"type:varchar(20);default:'pending'"`
}

// BeforeCreate will set a UUID rather than numeric ID
func (cr *ChatRequest) BeforeCreate(tx *gorm.DB) error {
	if cr.ID == "" {
		cr.ID = uuid.New().String()
	}
	return nil
}

type AcceptOrRejectChatRequest struct {
	Accept bool `json:"accept"`
}

type Message struct {
	ID        string         `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	RoomID    string         `json:"room_id" gorm:"not null;index"`
	Room      Room           `json:"room" gorm:"foreignKey:RoomID"`
	SenderID  uint           `json:"sender_id" gorm:"not null"`
	Content   string         `json:"content" gorm:"type:text;not null"`
}

type RoomWithLastMessage struct {
	Room
	LastMessage *Message `json:"last_message"`
}
