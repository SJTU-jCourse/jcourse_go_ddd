package common

type Pagination struct {
	Page int
	Size int
}

const (
	DefaultPage       = 1
	DefaultPageSize   = 20
	MinimumPageSize   = 1
	MinimumPageNumber = 1
)

func (p *Pagination) Offset() int {
	return (p.Page - 1) * p.Size
}

func NewPagination(page int, size int) Pagination {
	if page < MinimumPageNumber {
		page = DefaultPage
	}
	if size < MinimumPageSize {
		size = DefaultPageSize
	}
	return Pagination{
		Page: page,
		Size: size,
	}
}
