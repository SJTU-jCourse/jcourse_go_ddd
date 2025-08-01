package entity

import (
	"time"

	"gorm.io/gorm"
)

// ReviewRevision represents the review revision entity in the database
type ReviewRevision struct {
	ID        int    `gorm:"primaryKey"`
	ReviewID  int    `gorm:"not null"`
	Content   string `gorm:"type:text;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	// Relations
	Review Review `gorm:"foreignKey:ReviewID"`
}

// TableName specifies the table name for ReviewRevision
func (ReviewRevision) TableName() string {
	return "review_revisions"
}
