package models

import (
	"time"

	"gorm.io/gorm"
)

type Category struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	Name          string         `gorm:"not null" json:"name"` // Tên tiếng Việt (unique index được tạo thủ công với WHERE deleted_at IS NULL)
	NameEn        *string        `json:"nameEn"`               // Tên tiếng Anh
	Description   *string        `json:"description"`          // Mô tả tiếng Việt
	DescriptionEn *string        `json:"descriptionEn"`        // Mô tả tiếng Anh
	Image         *string        `json:"image"`
	IsActive      bool           `gorm:"default:true" json:"isActive"`
	CreatedAt     time.Time      `json:"createdAt"`
	UpdatedAt     time.Time      `json:"updatedAt"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	// Một category có thể là parent của nhiều categories (children)
	ParentRelations []CategoryChild `gorm:"foreignKey:ParentID" json:"parentRelations,omitempty"`
	// Một category có thể là child của một category (parent)
	ChildRelations []CategoryChild `gorm:"foreignKey:ChildID" json:"childRelations,omitempty"`
	Products       []Product       `gorm:"foreignKey:CategoryID" json:"products,omitempty"`
}

func (Category) TableName() string {
	return "categories"
}
