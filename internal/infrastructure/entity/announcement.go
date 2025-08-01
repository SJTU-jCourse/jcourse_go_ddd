package entity

import (
	"time"

	"gorm.io/gorm"
)

// Announcement represents the announcement entity in the database
type Announcement struct {
	ID          int       `gorm:"primaryKey"`
	Title       string    `gorm:"type:varchar(255);not null"`
	Content     string    `gorm:"type:text;not null"`
	IsPublished bool      `gorm:"not null;default:false"`
	PublishedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

// TableName specifies the table name for Announcement
func (Announcement) TableName() string {
	return "announcements"
}
