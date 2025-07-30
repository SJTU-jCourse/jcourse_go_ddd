package viewobject

import "jcourse_go/internal/domain/review"

type CourseInReviewVO struct {
	ID          int               `json:"id"`
	Code        string            `json:"code"`
	Name        string            `json:"name"`
	MainTeacher TeacherListItemVO `json:"main_teacher"`
}

func NewCourseInReviewVO(c *review.Course) CourseInReviewVO {
	return CourseInReviewVO{
		ID:          c.ID,
		Code:        c.Code,
		Name:        c.Name,
		MainTeacher: NewTeacherVO(c.MainTeacher),
	}
}

type ReviewVO struct {
	Course    *CourseInReviewVO
	ID        int
	CourseID  int
	Semester  string
	Grade     string
	Comment   string
	Rating    int
	CreatedAt int64
	UpdatedAt int64
}

type UserInReviewVO struct {
	ID   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

var AnonymousUserVO = UserInReviewVO{
	ID:   -1,
	Name: "匿名用户",
}

func NewReviewVO(r *review.Review, withCourse bool) ReviewVO {
	rvo := ReviewVO{
		ID:        r.ID,
		CourseID:  r.CourseID,
		Semester:  r.Semester.String(),
		Grade:     r.Grade,
		Comment:   r.Comment,
		Rating:    r.Rating.Int(),
		CreatedAt: r.CreatedAt.Unix(),
		UpdatedAt: r.UpdatedAt.Unix(),
	}
	if withCourse {
		course := NewCourseInReviewVO(r.Course)
		rvo.Course = &course
	}
	return rvo
}
