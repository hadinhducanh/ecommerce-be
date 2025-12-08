package models

import (
	"time"

	"gorm.io/gorm"
)

type PaymentMethod string

const (
	PaymentMethodCOD         PaymentMethod = "cod"
	PaymentMethodBankTransfer PaymentMethod = "bank_transfer"
	PaymentMethodCreditCard   PaymentMethod = "credit_card"
	PaymentMethodEWallet     PaymentMethod = "e_wallet"
)

type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "pending"
	PaymentStatusProcessing PaymentStatus = "processing"
	PaymentStatusCompleted PaymentStatus = "completed"
	PaymentStatusFailed   PaymentStatus = "failed"
	PaymentStatusRefunded PaymentStatus = "refunded"
)

type Payment struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	TransactionID  string         `gorm:"uniqueIndex;not null" json:"transactionId"` // Mã giao dịch
	Amount         float64        `gorm:"type:decimal(10,2);not null" json:"amount"`
	Method         PaymentMethod  `gorm:"type:varchar(50);default:'cod'" json:"method"`
	Status         PaymentStatus  `gorm:"type:varchar(50);default:'pending'" json:"status"`
	PaymentDetails *string        `gorm:"type:text" json:"paymentDetails"` // JSON string hoặc text
	Notes          *string        `json:"notes"`
	UserID         uint           `gorm:"not null" json:"userId"`
	OrderID        uint           `gorm:"not null" json:"orderId"`
	CreatedAt      time.Time      `json:"createdAt"`
	UpdatedAt      time.Time      `json:"updatedAt"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	User  User  `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Order Order `gorm:"foreignKey:OrderID" json:"order,omitempty"`
}

func (Payment) TableName() string {
	return "payments"
}

