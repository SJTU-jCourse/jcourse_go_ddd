package repository

import (
	"context"
	"fmt"
	"time"

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

func (r *statisticsRepository) GetCurrentStatistics(ctx context.Context) (*statistics.DailyStatistics, error) {
	// Try to get today's statistics first
	today := time.Now().Truncate(24 * time.Hour)
	var dailyStats entity.DailyStatistics
	err := r.db.Where("DATE(date) = ?", today.Format("2006-01-02")).First(&dailyStats).Error

	if err == nil {
		// Return today's statistics if available
		return &statistics.DailyStatistics{
			Date:               dailyStats.Date,
			DAU:                dailyStats.DAU,
			DNU:                dailyStats.DNU,
			MAU:                dailyStats.MAU,
			DailyNewReviews:    dailyStats.DailyNewReviews,
			TotalReviews:       dailyStats.TotalReviews,
			TotalCourses:       dailyStats.TotalCourses,
			CoursesWithReviews: dailyStats.CoursesWithReviews,
		}, nil
	}

	// If no statistics for today, calculate real-time statistics
	stats := &statistics.DailyStatistics{
		Date: today,
	}

	// Get DAU (Daily Active Users)
	var dau int64
	err = r.db.Model(&entity.UserSession{}).
		Where("DATE(created_at) = CURRENT_DATE").
		Distinct("user_id").
		Count(&dau).Error
	stats.DAU = int(dau)
	if err != nil {
		return nil, fmt.Errorf("failed to get DAU: %w", err)
	}

	// Get DNU (Daily New Users)
	var dnu int64
	err = r.db.Model(&entity.User{}).
		Where("DATE(created_at) = CURRENT_DATE").
		Count(&dnu).Error
	stats.DNU = int(dnu)
	if err != nil {
		return nil, fmt.Errorf("failed to get DNU: %w", err)
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

func (r *statisticsRepository) GetDailyStatistics(ctx context.Context, date time.Time) (*statistics.DailyStatistics, error) {
	var dailyStats entity.DailyStatistics
	err := r.db.Where("DATE(date) = ?", date.Format("2006-01-02")).First(&dailyStats).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get daily statistics: %w", err)
	}

	return &statistics.DailyStatistics{
		Date:               dailyStats.Date,
		DAU:                dailyStats.DAU,
		DNU:                dailyStats.DNU,
		MAU:                dailyStats.MAU,
		DailyNewReviews:    dailyStats.DailyNewReviews,
		TotalReviews:       dailyStats.TotalReviews,
		TotalCourses:       dailyStats.TotalCourses,
		CoursesWithReviews: dailyStats.CoursesWithReviews,
	}, nil
}

func (r *statisticsRepository) GetDailyStatisticsRange(ctx context.Context, startDate, endDate time.Time) ([]*statistics.DailyStatistics, error) {
	var dailyStats []entity.DailyStatistics
	err := r.db.Where("date BETWEEN ? AND ?", startDate, endDate).
		Order("date ASC").
		Find(&dailyStats).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get daily statistics range: %w", err)
	}

	result := make([]*statistics.DailyStatistics, len(dailyStats))
	for i, stats := range dailyStats {
		result[i] = &statistics.DailyStatistics{
			Date:               stats.Date,
			DAU:                stats.DAU,
			DNU:                stats.DNU,
			MAU:                stats.MAU,
			DailyNewReviews:    stats.DailyNewReviews,
			TotalReviews:       stats.TotalReviews,
			TotalCourses:       stats.TotalCourses,
			CoursesWithReviews: stats.CoursesWithReviews,
		}
	}

	return result, nil
}

func (r *statisticsRepository) SaveDailyStatistics(ctx context.Context, stats *statistics.DailyStatistics) error {
	dailyStats := &entity.DailyStatistics{
		Date:               stats.Date,
		DAU:                stats.DAU,
		DNU:                stats.DNU,
		MAU:                stats.MAU,
		DailyNewReviews:    stats.DailyNewReviews,
		TotalReviews:       stats.TotalReviews,
		TotalCourses:       stats.TotalCourses,
		CoursesWithReviews: stats.CoursesWithReviews,
	}

	// Use upsert to handle duplicate dates
	err := r.db.Where("date = ?", stats.Date).
		Assign(dailyStats).
		FirstOrCreate(dailyStats).Error
	if err != nil {
		return fmt.Errorf("failed to save daily statistics: %w", err)
	}

	return nil
}

func (r *statisticsRepository) GetLatestDailyStatistics(ctx context.Context) (*statistics.DailyStatistics, error) {
	var dailyStats entity.DailyStatistics
	err := r.db.Order("date DESC").First(&dailyStats).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get latest daily statistics: %w", err)
	}

	return &statistics.DailyStatistics{
		Date:               dailyStats.Date,
		DAU:                dailyStats.DAU,
		DNU:                dailyStats.DNU,
		MAU:                dailyStats.MAU,
		DailyNewReviews:    dailyStats.DailyNewReviews,
		TotalReviews:       dailyStats.TotalReviews,
		TotalCourses:       dailyStats.TotalCourses,
		CoursesWithReviews: dailyStats.CoursesWithReviews,
	}, nil
}
