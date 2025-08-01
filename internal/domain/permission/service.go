package permission

import (
	"context"

	"jcourse_go/internal/domain/auth"
	"jcourse_go/internal/domain/common"
)

type PermissionService interface {
	CheckPermission(ctx context.Context, ref ResourceRef, action Action, userID int) (Result, error)
}

type permissionService struct {
	userRepo auth.UserRepository
}

func (p *permissionService) CheckPermission(ctx context.Context, ref ResourceRef, action Action, userID int) (Result, error) {
	permCtx := &Ctx{
		ctx:    ctx,
		UserID: userID,
	}

	switch ref.Type {
	case ResourceTypeReview:
		return p.checkReviewPermission(permCtx, ref, action)
	case ResourceTypeUser:
		return p.checkUserPermission(permCtx, ref, action)
	case ResourceTypePoint:
		return p.checkPointPermission(permCtx, ref, action)
	case ResourceTypeCourse:
		return p.checkCoursePermission(permCtx, ref, action)
	default:
		return Result{Allow: false, Reason: "unknown resource type"}, nil
	}
}

func (p *permissionService) checkReviewPermission(permCtx *Ctx, ref ResourceRef, action Action) (Result, error) {
	switch action {
	case ActionView:
		return Result{Allow: true, Reason: "public access"}, nil
	case ActionCreate:
		if permCtx.UserID == 0 {
			return Result{Allow: false, Reason: "not authenticated"}, nil
		}
		return Result{Allow: true, Reason: "authenticated user"}, nil
	case ActionUpdate, ActionDelete:
		if permCtx.UserID == 0 {
			return Result{Allow: false, Reason: "not authenticated"}, nil
		}
		if ref.Owner.ID == 0 {
			return Result{Allow: false, Reason: "resource owner not found"}, nil
		}

		// Check if user is admin or owner
		if permCtx.UserID == ref.Owner.ID {
			return Result{Allow: true, Reason: "owner access"}, nil
		}

		// Check if user is admin
		if p.userRepo != nil {
			user, err := p.userRepo.GetByID(permCtx.ctx, permCtx.UserID)
			if err == nil && user.Role == common.RoleAdmin {
				return Result{Allow: true, Reason: "admin access"}, nil
			}
		}

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
		if permCtx.UserID == 0 {
			return Result{Allow: false, Reason: "not authenticated"}, nil
		}

		// Check if user is owner or admin
		if permCtx.UserID == ref.Owner.ID {
			return Result{Allow: true, Reason: "owner access"}, nil
		}

		// Check if user is admin
		if p.userRepo != nil {
			user, err := p.userRepo.GetByID(permCtx.ctx, permCtx.UserID)
			if err == nil && user.Role == common.RoleAdmin {
				return Result{Allow: true, Reason: "admin access"}, nil
			}
		}

		return Result{Allow: false, Reason: "permission denied"}, nil
	case ActionCreate, ActionDelete:
		return Result{Allow: false, Reason: "action not allowed"}, nil
	default:
		return Result{Allow: false, Reason: "unknown action"}, nil
	}
}

func (p *permissionService) checkPointPermission(permCtx *Ctx, ref ResourceRef, action Action) (Result, error) {
	if permCtx.UserID == 0 {
		return Result{Allow: false, Reason: "admin access required"}, nil
	}

	// Check if user is admin
	if p.userRepo != nil {
		user, err := p.userRepo.GetByID(permCtx.ctx, permCtx.UserID)
		if err == nil && user.Role == common.RoleAdmin {
			return Result{Allow: true, Reason: "admin access"}, nil
		}
	}

	return Result{Allow: false, Reason: "admin access required"}, nil
}

func (p *permissionService) checkCoursePermission(permCtx *Ctx, ref ResourceRef, action Action) (Result, error) {
	switch action {
	case ActionView:
		return Result{Allow: true, Reason: "public access"}, nil
	case ActionCreate, ActionUpdate, ActionDelete:
		if permCtx.UserID == 0 {
			return Result{Allow: false, Reason: "not authenticated"}, nil
		}

		// Check if user is admin
		if p.userRepo != nil {
			user, err := p.userRepo.GetByID(permCtx.ctx, permCtx.UserID)
			if err == nil && user.Role == common.RoleAdmin {
				return Result{Allow: true, Reason: "admin access"}, nil
			}
		}

		return Result{Allow: false, Reason: "admin access required"}, nil
	default:
		return Result{Allow: false, Reason: "unknown action"}, nil
	}
}

func NewPermissionService() PermissionService {
	return &permissionService{}
}

func NewPermissionServiceWithUserRepo(userRepo auth.UserRepository) PermissionService {
	return &permissionService{
		userRepo: userRepo,
	}
}
