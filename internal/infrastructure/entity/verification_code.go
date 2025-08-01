package entity

import (
	"time"

	"gorm.io/gorm"
)

// VerificationCode represents the verification code entity in the database
type VerificationCode struct {
	ID        int       `gorm:"primaryKey"`
	Email     string    `gorm:"type:varchar(100);uniqueIndex;not null"`
	Code      string    `gorm:"type:varchar(10);not null"`
	ExpiresAt time.Time `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// TableName specifies the table name for VerificationCode
func (VerificationCode) TableName() string {
	return "verification_codes"
}
