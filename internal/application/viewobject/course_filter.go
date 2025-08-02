package viewobject

type CourseFilterVO struct {
	Departments []string `json:"departments"`
	Categories  []string `json:"categories"`
}

func NewCourseFilterVO(departments, categories []string) CourseFilterVO {
	return CourseFilterVO{
		Departments: departments,
		Categories:  categories,
	}
}
