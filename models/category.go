package models

import (
	"time"

	"gorm.io/gorm"
)

type Category struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	Name          string         `gorm:"not null" json:"name"` // Tên tiếng Việt
	NameEn        *string        `json:"nameEn"`               // Tên tiếng Anh
	Description   *string        `json:"description"`          // Mô tả tiếng Việt
	DescriptionEn *string        `json:"descriptionEn"`        // Mô tả tiếng Anh
	Image         *string        `json:"image"`
	IsActive      bool           `gorm:"default:true" json:"isActive"`
	CreatedAt     time.Time      `json:"createdAt"`
	UpdatedAt     time.Time      `json:"updatedAt"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Products []Product `gorm:"foreignKey:CategoryID" json:"products,omitempty"`
}

func (Category) TableName() string {
	return "categories"
}
