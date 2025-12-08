package models

import (
	"time"

	"gorm.io/gorm"
)

type CartItem struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Quantity  int            `gorm:"default:1" json:"quantity"`
	UserID    uint           `gorm:"not null" json:"userId"`
	ProductID uint           `gorm:"not null" json:"productId"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	User    User    `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Product Product `gorm:"foreignKey:ProductID" json:"product,omitempty"`
}

func (CartItem) TableName() string {
	return "cart_items"
}
