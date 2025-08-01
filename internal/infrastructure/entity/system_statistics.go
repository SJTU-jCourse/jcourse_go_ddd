package entity

import (
	"time"
)

// DailyStatistics represents daily statistics stored in the database
type DailyStatistics struct {
	ID                 int       `gorm:"primaryKey"`
	Date               time.Time `gorm:"not null;uniqueIndex;comment:'Date of the statistics'"`
	DAU                int       `gorm:"not null;default:0;comment:'Daily Active Users'"`
	DNU                int       `gorm:"not null;default:0;comment:'Daily New Users'"`
	MAU                int       `gorm:"not null;default:0;comment:'Monthly Active Users'"`
	DailyNewReviews    int       `gorm:"not null;default:0;comment:'Daily New Reviews'"`
	TotalReviews       int       `gorm:"not null;default:0;comment:'Total Reviews'"`
	TotalCourses       int       `gorm:"not null;default:0;comment:'Total Courses'"`
	CoursesWithReviews int       `gorm:"not null;default:0;comment:'Courses with Reviews'"`
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

// TableName specifies the table name for DailyStatistics
func (DailyStatistics) TableName() string {
	return "daily_statistics"
}
