package viewobject

import "jcourse_go/internal/domain/statistics"

type SystemStatisticsVO struct {
	DAU                int `json:"dau"`
	MAU                int `json:"mau"`
	DailyNewReviews    int `json:"daily_new_reviews"`
	TotalCourses       int `json:"total_courses"`
	TotalReviews       int `json:"total_reviews"`
	CoursesWithReviews int `json:"courses_with_reviews"`
}

func NewSystemStatisticsVO(stats *statistics.SystemStatistics) SystemStatisticsVO {
	return SystemStatisticsVO{
		DAU:                stats.DAU,
		MAU:                stats.MAU,
		DailyNewReviews:    stats.DailyNewReviews,
		TotalCourses:       stats.TotalCourses,
		TotalReviews:       stats.TotalReviews,
		CoursesWithReviews: stats.CoursesWithReviews,
	}
}
