package statistics

import (
	"context"
	"time"
)

type DailyStatistics struct {
	Date               time.Time `json:"date"`
	DAU                int       `json:"dau"`
	DNU                int       `json:"dnu"`
	MAU                int       `json:"mau"`
	DailyNewReviews    int       `json:"daily_new_reviews"`
	TotalReviews       int       `json:"total_reviews"`
	TotalCourses       int       `json:"total_courses"`
	CoursesWithReviews int       `json:"courses_with_reviews"`
}

type StatisticsRepository interface {
	GetCurrentStatistics(ctx context.Context) (*DailyStatistics, error)
	GetDailyStatistics(ctx context.Context, date time.Time) (*DailyStatistics, error)
	GetDailyStatisticsRange(ctx context.Context, startDate, endDate time.Time) ([]*DailyStatistics, error)
	SaveDailyStatistics(ctx context.Context, stats *DailyStatistics) error
	GetLatestDailyStatistics(ctx context.Context) (*DailyStatistics, error)
}
