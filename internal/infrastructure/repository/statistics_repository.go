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
	var dau int64
	err := r.db.Model(&entity.UserSession{}).
		Where("DATE(created_at) = CURRENT_DATE").
		Distinct("user_id").
		Count(&dau).Error
	stats.DAU = int(dau)
	if err != nil {
		return nil, fmt.Errorf("failed to get DAU: %w", err)
	}

	// Get MAU (Monthly Active Users)
	var mau int64
	err = r.db.Model(&entity.UserSession{}).
		Where("created_at >= DATE_TRUNC('month', CURRENT_DATE)").
		Distinct("user_id").
		Count(&mau).Error
	stats.MAU = int(mau)
	if err != nil {
		return nil, fmt.Errorf("failed to get MAU: %w", err)
	}

	// Get Daily New Reviews
	var dailyReviews int64
	err = r.db.Model(&entity.Review{}).
		Where("DATE(created_at) = CURRENT_DATE").
		Count(&dailyReviews).Error
	stats.DailyNewReviews = int(dailyReviews)
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
	var coursesWithReviews int64
	err = r.db.Model(&entity.Review{}).
		Distinct("course_id").
		Count(&coursesWithReviews).Error
	stats.CoursesWithReviews = int(coursesWithReviews)
	if err != nil {
		return nil, fmt.Errorf("failed to get courses with reviews: %w", err)
	}

	return stats, nil
}
