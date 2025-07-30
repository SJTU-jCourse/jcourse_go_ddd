package viewobject

import "jcourse_go/internal/domain/point"

type UserPointVO struct {
	TotalPoint int                 `json:"total_point"`
	Records    []UserPointRecordVO `json:"records"`
}

type UserPointRecordVO struct {
	Point       int    `json:"point"`
	Description string `json:"description"`
	CreatedAt   int64  `json:"created_at"`
}

func NewUserPointRecordVO(r *point.UserPointRecord) UserPointRecordVO {
	return UserPointRecordVO{
		Point:       r.Point,
		Description: r.Description,
		CreatedAt:   r.CreatedAt.Unix(),
	}
}

func NewUserPointVO(p *point.UserPoint) UserPointVO {
	vo := UserPointVO{
		TotalPoint: p.TotalPoint,
		Records:    make([]UserPointRecordVO, 0),
	}
	for _, r := range p.Records {
		vo.Records = append(vo.Records, NewUserPointRecordVO(&r))
	}
	return vo
}
