package models

import (
	"time"

	"gorm.io/gorm"
)

type Wishlist struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	UserID    uint           `gorm:"not null;uniqueIndex:idx_user_product" json:"userId"`
	ProductID uint           `gorm:"not null;uniqueIndex:idx_user_product" json:"productId"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	User    User    `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Product Product `gorm:"foreignKey:ProductID" json:"product,omitempty"`
}

func (Wishlist) TableName() string {
	return "wishlists"
}

