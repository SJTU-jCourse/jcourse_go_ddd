package entity

import (
	"time"

	"gorm.io/gorm"
)

// UserEnrolledCourse represents the user enrolled course entity in the database
type UserEnrolledCourse struct {
	ID        int    `gorm:"primaryKey"`
	UserID    int    `gorm:"not null"`
	CourseID  int    `gorm:"not null"`
	Semester  string `gorm:"type:varchar(20);not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	// Relations
	User   User   `gorm:"foreignKey:UserID"`
	Course Course `gorm:"foreignKey:CourseID"`
}

// TableName specifies the table name for UserEnrolledCourse
func (UserEnrolledCourse) TableName() string {
	return "user_enrolled_courses"
}
