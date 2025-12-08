package models

import (
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Product struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	Name          string         `gorm:"not null;index" json:"name"` // Tên tiếng Việt
	NameEn        *string        `json:"nameEn"`                     // Tên tiếng Anh
	Description   *string        `gorm:"type:text" json:"description"`
	DescriptionEn *string        `gorm:"type:text" json:"descriptionEn"`
	Price         float64        `gorm:"type:decimal(10,2);not null;index" json:"price"`
	Stock         int            `gorm:"default:0" json:"stock"`
	Image         *string        `json:"image"`
	Images        pq.StringArray `gorm:"type:text[]" json:"images,omitempty"` // Nhiều hình ảnh
	Sold          int            `gorm:"default:0" json:"sold"`
	Rating        float64        `gorm:"default:0" json:"rating"`      // Điểm đánh giá trung bình (0-5)
	ReviewCount   int            `gorm:"default:0" json:"reviewCount"` // Số lượng đánh giá
	IsActive      bool           `gorm:"default:true" json:"isActive"`
	SKU           *string        `json:"sku"` // Stock Keeping Unit
	CategoryID    uint           `gorm:"not null;index" json:"categoryId"`
	CreatedAt     time.Time      `json:"createdAt"`
	UpdatedAt     time.Time      `json:"updatedAt"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Category   Category    `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	OrderItems []OrderItem `gorm:"foreignKey:ProductID" json:"orderItems,omitempty"`
	CartItems  []CartItem  `gorm:"foreignKey:ProductID" json:"cartItems,omitempty"`
	Wishlists  []Wishlist `gorm:"foreignKey:ProductID" json:"wishlists,omitempty"`
	Reviews    []Review    `gorm:"foreignKey:ProductID" json:"reviews,omitempty"`
}

func (Product) TableName() string {
	return "products"
}
