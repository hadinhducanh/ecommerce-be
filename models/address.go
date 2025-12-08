package models

import (
	"time"

	"gorm.io/gorm"
)

type Address struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	FullName  string         `gorm:"not null" json:"fullName"`
	Phone     string         `gorm:"not null" json:"phone"`
	Address   string         `gorm:"not null" json:"address"` // Địa chỉ chi tiết
	Ward      *string        `json:"ward"`                    // Phường/Xã
	District  *string        `json:"district"`                // Quận/Huyện
	City      *string        `json:"city"`                    // Tỉnh/Thành phố
	IsDefault bool           `gorm:"default:false" json:"isDefault"`
	UserID    uint           `gorm:"not null" json:"userId"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	User   User    `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Orders []Order `gorm:"foreignKey:ShippingAddressID" json:"orders,omitempty"`
}

func (Address) TableName() string {
	return "addresses"
}
