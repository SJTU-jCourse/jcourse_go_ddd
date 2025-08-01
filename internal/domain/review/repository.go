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
	SaveReviewAction(ctx context.Context, action *ReviewAction) error
	DeleteReviewAction(ctx context.Context, actionID int) error
	GetReviewAction(ctx context.Context, actionID int) (*ReviewAction, error)
	GetReviewRevisions(ctx context.Context, reviewID int) ([]ReviewRevision, error)
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
	GetDepartments(ctx context.Context) ([]string, error)
	GetCategories(ctx context.Context) ([]string, error)
	GetUserEnrolledCourses(ctx context.Context, userID int) ([]int, error)
	AddUserEnrolledCourse(ctx context.Context, userID int, courseID int) error
	WatchCourse(ctx context.Context, userID int, courseID int, watch bool) error
}
