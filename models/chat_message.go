package models

import (
	"time"

	"gorm.io/gorm"
)

type MessageType string

const (
	MessageTypeText   MessageType = "text"
	MessageTypeImage  MessageType = "image"
	MessageTypeFile   MessageType = "file"
	MessageTypeSystem MessageType = "system"
)

type ChatMessage struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Content   string         `gorm:"type:text;not null" json:"content"`
	Type      MessageType    `gorm:"type:varchar(50);default:'text'" json:"type"`
	IsRead    bool           `gorm:"default:false" json:"isRead"`
	ChatID    uint           `gorm:"not null" json:"chatId"`
	SenderID  uint           `gorm:"not null" json:"senderId"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Chat   Chat `gorm:"foreignKey:ChatID;constraint:OnDelete:CASCADE" json:"chat,omitempty"`
	Sender User `gorm:"foreignKey:SenderID" json:"sender,omitempty"`
}

func (ChatMessage) TableName() string {
	return "chat_messages"
}

