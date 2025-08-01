package statistics

import "context"

type SystemStatistics struct {
	DAU                int `json:"dau"`
	MAU                int `json:"mau"`
	DailyNewReviews    int `json:"daily_new_reviews"`
	TotalCourses       int `json:"total_courses"`
	TotalReviews       int `json:"total_reviews"`
	CoursesWithReviews int `json:"courses_with_reviews"`
}

type StatisticsRepository interface {
	GetSystemStatistics(ctx context.Context) (*SystemStatistics, error)
}
