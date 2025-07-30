package review

type BaseCourse struct {
	Code   string
	Name   string
	Credit float32
}

type Semester string

func (s *Semester) String() string {
	return string(*s)
}

func NewSemester(val string) Semester {
	return Semester(val)
}

type Rating int

func (r *Rating) Int() int {
	return int(*r)
}

func NewRating(val int) Rating {
	if val < 1 {
		return 1
	}
	if val > 5 {
		return 5
	}
	return Rating(val)
}

type Category string

type ReviewContent struct {
	Comment  string
	Rating   int
	Semester string
	Grade    string
}
