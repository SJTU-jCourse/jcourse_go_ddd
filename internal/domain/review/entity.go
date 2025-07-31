package review

import (
	"time"

	"jcourse_go/internal/domain/auth"
)

type Teacher struct {
	ID   int
	Code string // 工号
	Name string

	Department string
	Title      string

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

// Course = Teacher X BaseCourse
type Course struct {
	ID     int
	Code   string // 课程号
	Name   string
	Credit float32 // 学分

	MainTeacher   *Teacher
	MainTeacherID int

	OfferedCourses []OfferedCourse

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

func (c *Course) BaseCourse() BaseCourse {
	return BaseCourse{
		Code:   c.Code,
		Name:   c.Name,
		Credit: c.Credit,
	}
}

// OfferedCourse = Course X Semester
type OfferedCourse struct {
	CourseID int
	Semester Semester

	TeacherIDs   []int
	TeacherGroup []Teacher

	Categories []Category
	Grades     []string

	Language string
}

type Review struct {
	ID int

	CourseID int
	Course   *Course

	UserID int
	User   *auth.User

	Comment  string
	Rating   Rating
	Semester Semester
	Grade    string // 成绩

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

func (r *Review) Update(c *ReviewContent) {
	r.Comment = c.Comment
	r.Rating = NewRating(c.Rating)
	r.Semester = NewSemester(c.Semester)
	r.Grade = c.Grade
	r.UpdatedAt = time.Now()
}

type ReviewRevision struct {
	ID       int
	ReviewID int
	UserID   int

	Comment  string
	Semester Semester
	Grade    string
	Rating   Rating

	CreatedAt time.Time
	DeletedAt *time.Time
}

type ReviewAction struct {
	ID         int
	ReviewID   int
	UserID     int
	ActionType string

	CreatedAt time.Time
	DeletedAt *time.Time
}
