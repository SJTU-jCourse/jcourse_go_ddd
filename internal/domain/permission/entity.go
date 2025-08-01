package permission

import "context"

type Strategy interface {
	Check(permCtx *Ctx, ref ResourceRef) (Result, error)
}

type Ctx struct {
	ctx    context.Context
	UserID int
}

type Result struct {
	Allow  bool
	Reason string
}

type ResourceRef struct {
	ID    int
	Type  ResourceType
	Owner ResourceOwner
}

type ResourceOwner struct {
	ID int
}

type ResourceType int

const (
	ResourceTypeReview ResourceType = iota
	ResourceTypeUser   ResourceType = iota
	ResourceTypePoint  ResourceType = iota
	ResourceTypeCourse ResourceType = iota
)

type Action int

const (
	ActionView Action = iota
	ActionCreate
	ActionUpdate
	ActionDelete
)
