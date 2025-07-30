package permission

import "context"

type PermissionService interface {
	CheckPermission(ctx context.Context, ref ResourceRef, action Action, userID string) (Result, error)
}

type permissionService struct{}

func (p *permissionService) CheckPermission(ctx context.Context, ref ResourceRef, action Action, userID string) (Result, error) {
	// TODO implement me
	panic("implement me")
}

func NewPermissionService() PermissionService {
	return &permissionService{}
}
