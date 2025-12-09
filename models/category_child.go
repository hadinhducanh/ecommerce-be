package models

import (
	"time"

	"gorm.io/gorm"
)

// CategoryChild lưu quan hệ parent-child giữa các categories
// Một category có thể có nhiều children, một category chỉ có thể là child của một parent
type CategoryChild struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	ParentID  uint           `gorm:"not null;index" json:"parentId"` // ID của category cha
	ChildID   uint           `gorm:"not null;index" json:"childId"`  // ID của category con
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Parent *Category `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Child  *Category `gorm:"foreignKey:ChildID" json:"child,omitempty"`
}

func (CategoryChild) TableName() string {
	return "category_children"
}

// Unique constraint: một child chỉ có thể thuộc về một parent
// Được đảm bảo bởi unique index trên child_id
