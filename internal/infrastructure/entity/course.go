package entity

import (
	"time"

	"gorm.io/gorm"
)

// Course represents the course entity in the database
type Course struct {
	ID            int     `gorm:"primaryKey"`
	Name          string  `gorm:"type:varchar(255);not null"`
	Code          string  `gorm:"type:varchar(50);not null"`
	MainTeacherID int     `gorm:"not null"`
	Department    string  `gorm:"type:varchar(100)"`
	Credits       float64 `gorm:"type:decimal(3,1)"`
	Description   string  `gorm:"type:text"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"`

	// Relations
	MainTeacher User `gorm:"foreignKey:MainTeacherID"`
}

// TableName specifies the table name for Course
func (Course) TableName() string {
	return "courses"
}
