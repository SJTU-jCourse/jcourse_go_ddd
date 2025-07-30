package viewobject

import "jcourse_go/internal/domain/review"

type TeacherListItemVO struct {
	ID   int
	Code string
	Name string
}

func NewTeacherVO(t *review.Teacher) TeacherListItemVO {
	return TeacherListItemVO{
		ID:   t.ID,
		Code: t.Code,
		Name: t.Name,
	}
}

type OfferedCourseVO struct {
	Semester     string              `json:"semester"`
	Department   string              `json:"department"`
	Grades       []string            `json:"grades"`
	Categories   []string            `json:"categories"`
	TeacherGroup []TeacherListItemVO `json:"teacher_group"`
}
type RatingInfoVO struct {
	Count int                   `json:"count"`
	Avg   float32               `json:"avg"`
	Dist  map[review.Rating]int `json:"dist"`
}

type CourseListItemVO struct {
	ID          int               `json:"id"`
	Code        string            `json:"code"`
	Name        string            `json:"name"`
	Credit      float32           `json:"credit"`
	MainTeacher TeacherListItemVO `json:"main_teacher"`
	Rating      RatingInfoVO      `json:"rating"`

	// 最新开课记录
	Categories []string `json:"categories"`
	Department string   `json:"department"`
}

type CourseDetailVO struct {
	ID          int               `json:"id"`
	Code        string            `json:"code"`
	Name        string            `json:"name"`
	Credit      float32           `json:"credit"`
	MainTeacher TeacherListItemVO `json:"main_teacher"`
	Rating      RatingInfoVO      `json:"rating"`

	OfferedCourses []OfferedCourseVO `json:"offered_courses,omitempty"`

	CoursesUnderSameTeacher []CourseListItemVO `json:"courses_under_same_teacher"`
	CoursesByOtherTeachers  []CourseListItemVO `json:"courses_by_other_teachers"`
}

func NewCourseListItemVO(c *review.Course) CourseListItemVO {
	return CourseListItemVO{}
}

func NewCourseDetailVO(c *review.Course) CourseDetailVO {
	return CourseDetailVO{}
}
