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
	mainTeacher := TeacherListItemVO{}
	if c.MainTeacher != nil {
		mainTeacher = NewTeacherVO(c.MainTeacher)
	}

	return CourseListItemVO{
		ID:          c.ID,
		Code:        c.Code,
		Name:        c.Name,
		Credit:      c.Credit,
		MainTeacher: mainTeacher,
		Rating:      NewRatingInfoVO(), // Default empty rating
		Categories:  []string{},        // Will be populated from latest offered course
		Department:  "",                // Will be populated from latest offered course
	}
}

func NewCourseDetailVO(c *review.Course) CourseDetailVO {
	mainTeacher := TeacherListItemVO{}
	if c.MainTeacher != nil {
		mainTeacher = NewTeacherVO(c.MainTeacher)
	}

	offeredCourses := []OfferedCourseVO{}
	for _, oc := range c.OfferedCourses {
		offeredCourses = append(offeredCourses, NewOfferedCourseVO(oc))
	}

	return CourseDetailVO{
		ID:                      c.ID,
		Code:                    c.Code,
		Name:                    c.Name,
		Credit:                  c.Credit,
		MainTeacher:             mainTeacher,
		Rating:                  NewRatingInfoVO(), // Default empty rating
		OfferedCourses:          offeredCourses,
		CoursesUnderSameTeacher: []CourseListItemVO{}, // Will be populated separately
		CoursesByOtherTeachers:  []CourseListItemVO{}, // Will be populated separately
	}
}

func NewRatingInfoVO() RatingInfoVO {
	return RatingInfoVO{
		Count: 0,
		Avg:   0.0,
		Dist:  map[review.Rating]int{},
	}
}

func NewOfferedCourseVO(oc review.OfferedCourse) OfferedCourseVO {
	teacherGroup := []TeacherListItemVO{}
	for _, t := range oc.TeacherGroup {
		teacherGroup = append(teacherGroup, NewTeacherVO(&t))
	}

	categories := []string{}
	for _, cat := range oc.Categories {
		categories = append(categories, string(cat))
	}

	return OfferedCourseVO{
		Semester:     oc.Semester.String(),
		Department:   "", // Will be populated from teacher info
		Grades:       oc.Grades,
		Categories:   categories,
		TeacherGroup: teacherGroup,
	}
}
