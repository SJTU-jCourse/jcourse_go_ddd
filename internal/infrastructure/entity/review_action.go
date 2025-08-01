package entity

import (
	"time"

	"gorm.io/gorm"
)

// ReviewAction represents the review action entity in the database
type ReviewAction struct {
	ID          int    `gorm:"primaryKey"`
	ReviewID    int    `gorm:"not null"`
	UserID      int    `gorm:"not null"`
	Action      string `gorm:"type:varchar(50);not null"`
	Description string `gorm:"type:text"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`

	// Relations
	Review Review `gorm:"foreignKey:ReviewID"`
	User   User   `gorm:"foreignKey:UserID"`
}

// TableName specifies the table name for ReviewAction
func (ReviewAction) TableName() string {
	return "review_actions"
}
