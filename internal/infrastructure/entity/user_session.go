package entity

import (
	"time"

	"gorm.io/gorm"
)

// UserSession represents the user session entity in the database
type UserSession struct {
	ID        int    `gorm:"primaryKey"`
	UserID    int    `gorm:"not null"`
	Token     string `gorm:"type:varchar(255);not null"`
	ExpiresAt time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	// Relations
	User User `gorm:"foreignKey:UserID"`
}

// TableName specifies the table name for UserSession
func (UserSession) TableName() string {
	return "user_sessions"
}
