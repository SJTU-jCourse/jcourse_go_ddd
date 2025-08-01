package entity

import (
	"time"

	"gorm.io/gorm"
)

// UserPointRecord represents the user point record entity in the database
type UserPointRecord struct {
	ID          int    `gorm:"primaryKey"`
	UserID      int    `gorm:"not null"`
	Point       int    `gorm:"not null"`
	Action      string `gorm:"type:varchar(50);not null"`
	Description string `gorm:"type:text"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`

	// Relations
	User User `gorm:"foreignKey:UserID"`
}

// TableName specifies the table name for UserPointRecord
func (UserPointRecord) TableName() string {
	return "user_point_records"
}
