package entity

import (
	"time"
)

// SystemStatistics represents the system statistics entity in the database
type SystemStatistics struct {
	ID                 int       `gorm:"primaryKey"`
	DAU                int       `gorm:"not null;default:0;comment:'Daily Active Users'"`
	MAU                int       `gorm:"not null;default:0;comment:'Monthly Active Users'"`
	DailyNewReviews    int       `gorm:"not null;default:0;comment:'Daily New Reviews'"`
	TotalCourses       int       `gorm:"not null;default:0;comment:'Total Courses'"`
	TotalReviews       int       `gorm:"not null;default:0;comment:'Total Reviews'"`
	CoursesWithReviews int       `gorm:"not null;default:0;comment:'Courses with Reviews'"`
	RecordedAt         time.Time `gorm:"not null;default:CURRENT_TIMESTAMP;comment:'When this statistic was recorded'"`
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

// TableName specifies the table name for SystemStatistics
func (SystemStatistics) TableName() string {
	return "system_statistics"
}
