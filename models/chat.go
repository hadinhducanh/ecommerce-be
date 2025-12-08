package models

import (
	"time"

	"gorm.io/gorm"
)

type ChatStatus string

const (
	ChatStatusOpen    ChatStatus = "open"
	ChatStatusClosed  ChatStatus = "closed"
	ChatStatusPending ChatStatus = "pending"
)

type Chat struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	Subject    string         `gorm:"not null" json:"subject"` // Chủ đề chat
	Status     ChatStatus     `gorm:"type:varchar(50);default:'open'" json:"status"`
	CustomerID uint           `gorm:"not null" json:"customerId"`
	AdminID    *uint          `json:"adminId"` // Admin phụ trách (nullable)
	CreatedAt  time.Time      `json:"createdAt"`
	UpdatedAt  time.Time      `json:"updatedAt"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Customer User         `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
	Admin    *User        `gorm:"foreignKey:AdminID" json:"admin,omitempty"`
	Messages []ChatMessage `gorm:"foreignKey:ChatID;constraint:OnDelete:CASCADE" json:"messages,omitempty"`
}

func (Chat) TableName() string {
	return "chats"
}

