package models

import (
	"time"

	"gorm.io/gorm"
)

type OrderItem struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Quantity  int            `gorm:"not null" json:"quantity"`
	Price     float64        `gorm:"type:decimal(10,2);not null" json:"price"` // Giá tại thời điểm đặt hàng
	Total     float64        `gorm:"type:decimal(10,2);not null" json:"total"`  // quantity * price
	OrderID   uint           `gorm:"not null" json:"orderId"`
	ProductID uint           `gorm:"not null" json:"productId"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Order   Order   `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE" json:"order,omitempty"`
	Product Product `gorm:"foreignKey:ProductID" json:"product,omitempty"`
}

func (OrderItem) TableName() string {
	return "order_items"
}

