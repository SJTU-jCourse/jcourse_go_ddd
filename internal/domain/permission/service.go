package permission

import (
	"context"
	"strconv"
)

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
		
		// Convert userID string to int for comparison
		userID, err := strconv.Atoi(permCtx.UserID)
		if err != nil {
			return Result{Allow: false, Reason: "invalid user ID"}, nil
		}
		
		// Check if user is admin or owner
		if userID == ref.Owner.ID {
			return Result{Allow: true, Reason: "owner access"}, nil
		}
		
		// TODO: Check admin role - this would require user repository
		// For now, we'll implement this check in the application layer
		return Result{Allow: false, Reason: "permission denied"}, nil
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
