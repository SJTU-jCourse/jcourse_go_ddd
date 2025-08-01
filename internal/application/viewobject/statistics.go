package viewobject

import (
	"jcourse_go/internal/domain/statistics"
)

type DailyStatisticsVO struct {
	Date               string `json:"date"`
	DAU                int    `json:"dau"`
	DNU                int    `json:"dnu"`
	MAU                int    `json:"mau"`
	DailyNewReviews    int    `json:"daily_new_reviews"`
	TotalReviews       int    `json:"total_reviews"`
	TotalCourses       int    `json:"total_courses"`
	CoursesWithReviews int    `json:"courses_with_reviews"`
}

func NewDailyStatisticsVO(stats *statistics.DailyStatistics) DailyStatisticsVO {
	return DailyStatisticsVO{
		Date:               stats.Date.Format("2006-01-02"),
		DAU:                stats.DAU,
		DNU:                stats.DNU,
		MAU:                stats.MAU,
		DailyNewReviews:    stats.DailyNewReviews,
		TotalReviews:       stats.TotalReviews,
		TotalCourses:       stats.TotalCourses,
		CoursesWithReviews: stats.CoursesWithReviews,
	}
}

func NewDailyStatisticsVOList(stats []*statistics.DailyStatistics) []DailyStatisticsVO {
	result := make([]DailyStatisticsVO, len(stats))
	for i, stat := range stats {
		result[i] = NewDailyStatisticsVO(stat)
	}
	return result
}
