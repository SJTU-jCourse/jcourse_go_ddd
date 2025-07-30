package common

type Pagination struct {
	Page int
	Size int
}

func (p *Pagination) Offset() int {
	return (p.Page - 1) * p.Size
}

func NewPagination(page int, size int) Pagination {
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 20
	}
	return Pagination{
		Page: page,
		Size: size,
	}
}
