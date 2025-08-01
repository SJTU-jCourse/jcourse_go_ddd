package entity

import (
	"time"

	"gorm.io/gorm"
)

// Review represents the review entity in the database
type Review struct {
	ID        int    `gorm:"primaryKey"`
	UserID    int    `gorm:"not null"`
	CourseID  int    `gorm:"not null"`
	Rating    int    `gorm:"not null;check:rating >= 1 AND rating <= 5"`
	Semester  string `gorm:"type:varchar(20);not null"`
	Content   string `gorm:"type:text;not null"`
	Category  string `gorm:"type:varchar(50);not null"`
	IsPublic  bool   `gorm:"not null;default:true"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	// Relations
	User   User   `gorm:"foreignKey:UserID"`
	Course Course `gorm:"foreignKey:CourseID"`
}

// TableName specifies the table name for Review
func (Review) TableName() string {
	return "reviews"
}
