package permission

import "context"

type PermissionService interface {
	CheckPermission(ctx context.Context, ref ResourceRef, action Action, userID string) (Result, error)
}

type permissionService struct{}

func (p *permissionService) CheckPermission(ctx context.Context, ref ResourceRef, action Action, userID string) (Result, error) {
	permCtx := &Ctx{
		ctx:    ctx,
		UserID: userID,
	}

	switch ref.Type {
	case ResourceTypeReview:
		return p.checkReviewPermission(permCtx, ref, action)
	case ResourceTypeUser:
		return p.checkUserPermission(permCtx, ref, action)
	default:
		return Result{Allow: false, Reason: "unknown resource type"}, nil
	}
}

func (p *permissionService) checkReviewPermission(permCtx *Ctx, ref ResourceRef, action Action) (Result, error) {
	switch action {
	case ActionView:
		return Result{Allow: true, Reason: "public access"}, nil
	case ActionCreate:
		return Result{Allow: true, Reason: "authenticated user"}, nil
	case ActionUpdate, ActionDelete:
		if permCtx.UserID == "" {
			return Result{Allow: false, Reason: "not authenticated"}, nil
		}
		if ref.Owner.ID == 0 {
			return Result{Allow: false, Reason: "resource owner not found"}, nil
		}
		// TODO: check if user is owner or has admin role
		// For now, allow owner to edit/delete their own reviews
		return Result{Allow: true, Reason: "owner access"}, nil
	default:
		return Result{Allow: false, Reason: "unknown action"}, nil
	}
}

func (p *permissionService) checkUserPermission(permCtx *Ctx, ref ResourceRef, action Action) (Result, error) {
	switch action {
	case ActionView:
		return Result{Allow: true, Reason: "public access"}, nil
	case ActionUpdate:
		if permCtx.UserID == "" {
			return Result{Allow: false, Reason: "not authenticated"}, nil
		}
		// TODO: check if user is owner or has admin role
		// For now, allow users to update their own profile
		return Result{Allow: true, Reason: "owner access"}, nil
	case ActionCreate, ActionDelete:
		return Result{Allow: false, Reason: "action not allowed"}, nil
	default:
		return Result{Allow: false, Reason: "unknown action"}, nil
	}
}

func NewPermissionService() PermissionService {
	return &permissionService{}
}
