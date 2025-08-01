package review

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"jcourse_go/internal/domain/review"
	"jcourse_go/internal/infrastructure/entity"
)

type courseRepository struct {
	db *gorm.DB
}

func NewCourseRepository(db *gorm.DB) review.CourseRepository {
	return &courseRepository{db: db}
}

func (r *courseRepository) Get(ctx context.Context, id int) (*review.Course, error) {
	var courseEntity entity.Course
	result := r.db.WithContext(ctx).
		Preload("MainTeacher").
		First(&courseEntity, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get course: %w", result.Error)
	}
	return r.toDomainCourse(&courseEntity), nil
}

func (r *courseRepository) FindBy(ctx context.Context, filter review.CourseFilter) ([]review.Course, error) {
	var courseEntities []entity.Course
	query := r.db.WithContext(ctx).Preload("MainTeacher")

	if filter.MainTeacherID != nil {
		query = query.Where("main_teacher_id = ?", *filter.MainTeacherID)
	}
	if filter.Code != nil {
		query = query.Where("code = ?", *filter.Code)
	}
	if filter.Name != nil {
		query = query.Where("name LIKE ?", "%"+*filter.Name+"%")
	}
	if len(filter.Credit) > 0 {
		query = query.Where("credit IN ?", filter.Credit)
	}

	if filter.HasReviews {
		query = query.Joins("JOIN reviews ON courses.id = reviews.course_id").
			Group("courses.id")
	}

	result := query.Find(&courseEntities)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to find courses: %w", result.Error)
	}

	courses := make([]review.Course, len(courseEntities))
	for i, courseEntity := range courseEntities {
		courses[i] = *r.toDomainCourse(&courseEntity)
	}

	return courses, nil
}

func (r *courseRepository) Save(ctx context.Context, course *review.Course) error {
	courseORM := r.toORMCourse(course)
	result := r.db.WithContext(ctx).Save(courseORM)
	if result.Error != nil {
		return fmt.Errorf("failed to save course: %w", result.Error)
	}
	return nil
}

func (r *courseRepository) Delete(ctx context.Context, filter review.CourseFilter) error {
	query := r.db.WithContext(ctx)
	if filter.MainTeacherID != nil {
		query = query.Where("main_teacher_id = ?", *filter.MainTeacherID)
	}
	if filter.Code != nil {
		query = query.Where("code = ?", *filter.Code)
	}

	result := query.Delete(&entity.Course{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete courses: %w", result.Error)
	}

	return nil
}

func (r *courseRepository) FindOfferedCourse(ctx context.Context, courseID int, semester review.Semester) (*review.OfferedCourse, error) {
	var offeredCourse review.OfferedCourse
	result := r.db.WithContext(ctx).
		Where("course_id = ? AND semester = ?", courseID, semester).
		First(&offeredCourse)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find offered course: %w", result.Error)
	}
	return &offeredCourse, nil
}

func (r *courseRepository) GetDepartments(ctx context.Context) ([]string, error) {
	var departments []string
	result := r.db.WithContext(ctx).
		Model(&entity.User{}).
		Distinct("department").
		Find(&departments)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get departments: %w", result.Error)
	}
	return departments, nil
}

func (r *courseRepository) GetCategories(ctx context.Context) ([]string, error) {
	var categories []string
	result := r.db.WithContext(ctx).
		Table("offered_courses").
		Distinct("unnest(categories)").
		Find(&categories)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get categories: %w", result.Error)
	}
	return categories, nil
}

func (r *courseRepository) GetUserEnrolledCourses(ctx context.Context, userID int) ([]int, error) {
	var courseIDs []int
	result := r.db.WithContext(ctx).
		Table("user_enrolled_courses").
		Where("user_id = ?", userID).
		Pluck("course_id", &courseIDs)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get user enrolled courses: %w", result.Error)
	}
	return courseIDs, nil
}

func (r *courseRepository) AddUserEnrolledCourse(ctx context.Context, userID int, courseID int) error {
	result := r.db.WithContext(ctx).
		Exec("INSERT INTO user_enrolled_courses (user_id, course_id) VALUES (?, ?) ON CONFLICT (user_id, course_id) DO NOTHING", userID, courseID)
	if result.Error != nil {
		return fmt.Errorf("failed to add user enrolled course: %w", result.Error)
	}
	return nil
}

func (r *courseRepository) WatchCourse(ctx context.Context, userID int, courseID int, watch bool) error {
	if watch {
		result := r.db.WithContext(ctx).
			Exec("INSERT INTO course_watches (user_id, course_id) VALUES (?, ?) ON CONFLICT (user_id, course_id) DO NOTHING", userID, courseID)
		if result.Error != nil {
			return fmt.Errorf("failed to watch course: %w", result.Error)
		}
	} else {
		result := r.db.WithContext(ctx).
			Exec("DELETE FROM course_watches WHERE user_id = ? AND course_id = ?", userID, courseID)
		if result.Error != nil {
			return fmt.Errorf("failed to unwatch course: %w", result.Error)
		}
	}
	return nil
}

// Helper methods to convert between domain and ORM models
func (r *courseRepository) toDomainCourse(courseEntity *entity.Course) *review.Course {
	return &review.Course{
		ID:            courseEntity.ID,
		Name:          courseEntity.Name,
		Code:          courseEntity.Code,
		MainTeacherID: courseEntity.MainTeacherID,
		Credit:        float32(courseEntity.Credits),
	}
}

func (r *courseRepository) toORMCourse(course *review.Course) *entity.Course {
	return &entity.Course{
		ID:            course.ID,
		Name:          course.Name,
		Code:          course.Code,
		MainTeacherID: course.MainTeacherID,
		Department:    "Computer Science", // Default department
		Credits:       float64(course.Credit),
		Description:   "", // Default empty description
	}
}
