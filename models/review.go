package models

import (
	"time"

	"gorm.io/gorm"
	"github.com/lib/pq"
)

type Review struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	Rating     int            `gorm:"not null" json:"rating"` // 1-5 sao
	Comment    *string        `gorm:"type:text" json:"comment"`
	Images     pq.StringArray `gorm:"type:text[]" json:"images,omitempty"` // Hình ảnh đánh giá
	IsVerified bool           `gorm:"default:false" json:"isVerified"`     // Đã mua hàng
	UserID     uint           `gorm:"not null" json:"userId"`
	ProductID  uint           `gorm:"not null" json:"productId"`
	OrderID    *uint          `json:"orderId"` // Link với đơn hàng để verify
	CreatedAt  time.Time      `json:"createdAt"`
	UpdatedAt  time.Time      `json:"updatedAt"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	User    User    `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Product Product `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	Order   *Order  `gorm:"foreignKey:OrderID" json:"order,omitempty"`
}

func (Review) TableName() string {
	return "reviews"
}

