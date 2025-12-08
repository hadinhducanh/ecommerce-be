package models

import (
	"time"

	"gorm.io/gorm"
)

type OrderStatus string

const (
	OrderStatusPending    OrderStatus = "pending"
	OrderStatusConfirmed  OrderStatus = "confirmed"
	OrderStatusProcessing OrderStatus = "processing"
	OrderStatusShipping   OrderStatus = "shipping"
	OrderStatusDelivered  OrderStatus = "delivered"
	OrderStatusCancelled  OrderStatus = "cancelled"
	OrderStatusRefunded   OrderStatus = "refunded"
)

type Order struct {
	ID                uint           `gorm:"primaryKey" json:"id"`
	OrderNumber       string         `gorm:"uniqueIndex;not null" json:"orderNumber"` // Mã đơn hàng
	TotalAmount       float64        `gorm:"type:decimal(10,2);not null" json:"totalAmount"`
	ShippingFee       float64        `gorm:"type:decimal(10,2);default:0" json:"shippingFee"`
	Discount          float64        `gorm:"type:decimal(10,2);default:0" json:"discount"`
	Status            OrderStatus    `gorm:"type:varchar(50);default:'pending'" json:"status"`
	Notes             *string        `json:"notes"` // Ghi chú của khách hàng
	UserID            uint           `gorm:"not null" json:"userId"`
	ShippingAddressID uint           `gorm:"not null" json:"shippingAddressId"`
	CreatedAt         time.Time      `json:"createdAt"`
	UpdatedAt         time.Time      `json:"updatedAt"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	User            User        `gorm:"foreignKey:UserID" json:"user,omitempty"`
	ShippingAddress Address     `gorm:"foreignKey:ShippingAddressID" json:"shippingAddress,omitempty"`
	Items           []OrderItem `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE" json:"items,omitempty"`
	Payments        []Payment   `gorm:"foreignKey:OrderID" json:"payments,omitempty"`
}

func (Order) TableName() string {
	return "orders"
}
