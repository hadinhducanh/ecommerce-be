package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	Email           string         `gorm:"uniqueIndex;not null" json:"email"`
	Password        string         `gorm:"not null" json:"-"`
	Name            string         `gorm:"not null" json:"name"`
	Role            string         `gorm:"type:varchar(50);default:'customer'" json:"role"` // 'guest', 'customer', 'admin'
	Phone           *string        `json:"phone"`
	Avatar          *string        `json:"avatar"`
	Address         *string        `json:"address"`
	Gender          *string        `json:"gender"` // 'male', 'female', 'other'
	IsActive        bool           `gorm:"default:true" json:"isActive"`
	RefreshToken    *string        `json:"-"`
	OTP             *string        `json:"-"`
	OTPExpiresAt    *time.Time     `json:"-"`
	LastOTPSentAt   *time.Time     `json:"-"`
	IsEmailVerified bool           `gorm:"default:false" json:"isEmailVerified"`
	IsFirstLogin    bool           `gorm:"default:false" json:"isFirstLogin"`
	CreatedAt       time.Time      `json:"createdAt"`
	UpdatedAt       time.Time      `json:"updatedAt"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Addresses []Address  `gorm:"foreignKey:UserID" json:"addresses,omitempty"`
	Orders    []Order    `gorm:"foreignKey:UserID" json:"orders,omitempty"`
	CartItems []CartItem `gorm:"foreignKey:UserID" json:"cartItems,omitempty"`
	Payments  []Payment  `gorm:"foreignKey:UserID" json:"payments,omitempty"`
	Wishlists []Wishlist `gorm:"foreignKey:UserID" json:"wishlists,omitempty"`
	Reviews   []Review   `gorm:"foreignKey:UserID" json:"reviews,omitempty"`
	// Chat và Notification sẽ được xử lý bởi microservices riêng (NoSQL)
	// CustomerChats []Chat     `gorm:"foreignKey:CustomerID" json:"customerChats,omitempty"`
	// AdminChats    []Chat     `gorm:"foreignKey:AdminID" json:"adminChats,omitempty"`
}

func (User) TableName() string {
	return "users"
}
