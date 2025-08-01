package entity

import (
	"time"

	"gorm.io/gorm"
)

// User represents the user entity in the database
type User struct {
	ID           int    `gorm:"primaryKey"`
	Username     string `gorm:"type:varchar(50);uniqueIndex;not null"`
	Email        string `gorm:"type:varchar(100);uniqueIndex;not null"`
	PasswordHash string `gorm:"type:varchar(255);not null"`
	Role         string `gorm:"type:varchar(20);not null;default:'user'"`
	IsVerified   bool   `gorm:"not null;default:false"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

// TableName specifies the table name for User
func (User) TableName() string {
	return "users"
}
