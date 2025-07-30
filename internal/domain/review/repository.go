package review

import "context"

type ReviewFilter struct {
	ReviewID      *int
	UserID        *int
	CourseID      *int
	MainTeacherID *int
	Semester      *string
	Rating        *int
}

type ReviewRepository interface {
	Get(ctx context.Context, id int) (*Review, error)
	FindBy(ctx context.Context, filter ReviewFilter) ([]Review, error)
	Save(ctx context.Context, review *Review, revision *ReviewRevision) error
	Delete(ctx context.Context, filter ReviewFilter) error
}

type CourseFilter struct {
	MainTeacherID *int
	Code          *string
	Name          *string
	Credit        []float32
	Categories    []string
	Departments   []string

	HasReviews bool
}

type CourseRepository interface {
	Get(ctx context.Context, id int) (*Course, error)
	FindBy(ctx context.Context, filter CourseFilter) ([]Course, error)
	Save(ctx context.Context, course *Course) error
	Delete(ctx context.Context, filter CourseFilter) error

	FindOfferedCourse(ctx context.Context, courseID int, semester Semester) (*OfferedCourse, error)
}
