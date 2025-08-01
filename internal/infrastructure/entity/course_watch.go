package entity

import (
	"time"

	"gorm.io/gorm"
)

// CourseWatch represents the course watch entity in the database
type CourseWatch struct {
	ID        int `gorm:"primaryKey"`
	UserID    int `gorm:"not null"`
	CourseID  int `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	// Relations
	User   User   `gorm:"foreignKey:UserID"`
	Course Course `gorm:"foreignKey:CourseID"`
}

// TableName specifies the table name for CourseWatch
func (CourseWatch) TableName() string {
	return "course_watches"
}
