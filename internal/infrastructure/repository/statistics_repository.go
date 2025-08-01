package repository

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"jcourse_go/internal/domain/statistics"
	"jcourse_go/internal/infrastructure/entity"
)

type statisticsRepository struct {
	db *gorm.DB
}

func NewStatisticsRepository(db *gorm.DB) statistics.StatisticsRepository {
	return &statisticsRepository{db: db}
}

func (r *statisticsRepository) GetSystemStatistics(ctx context.Context) (*statistics.SystemStatistics, error) {
	stats := &statistics.SystemStatistics{}

	// Get DAU (Daily Active Users)
	err := r.db.Raw("SELECT COUNT(DISTINCT user_id) FROM user_sessions WHERE DATE(created_at) = CURRENT_DATE").Scan(&stats.DAU).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get DAU: %w", err)
	}

	// Get MAU (Monthly Active Users)
	err = r.db.Raw("SELECT COUNT(DISTINCT user_id) FROM user_sessions WHERE created_at >= DATE_TRUNC('month', CURRENT_DATE)").Scan(&stats.MAU).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get MAU: %w", err)
	}

	// Get Daily New Reviews
	err = r.db.Raw("SELECT COUNT(*) FROM reviews WHERE DATE(created_at) = CURRENT_DATE").Scan(&stats.DailyNewReviews).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get daily new reviews: %w", err)
	}

	// Get Total Courses
	var totalCourses int64
	err = r.db.Model(&entity.Course{}).Count(&totalCourses).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get total courses: %w", err)
	}
	stats.TotalCourses = int(totalCourses)

	// Get Total Reviews
	var totalReviews int64
	err = r.db.Model(&entity.Review{}).Count(&totalReviews).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get total reviews: %w", err)
	}
	stats.TotalReviews = int(totalReviews)

	// Get Courses with Reviews
	err = r.db.Raw("SELECT COUNT(DISTINCT course_id) FROM reviews").Scan(&stats.CoursesWithReviews).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get courses with reviews: %w", err)
	}

	return stats, nil
}
