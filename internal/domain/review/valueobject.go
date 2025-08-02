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
	if val < MinRating {
		return MinRating
	}
	if val > MaxRating {
		return MaxRating
	}
	return Rating(val)
}

type Category string

const (
	MinRating = 1
	MaxRating = 5
)

type ReviewContent struct {
	Comment  string
	Rating   int
	Semester string
	Grade    string
}
