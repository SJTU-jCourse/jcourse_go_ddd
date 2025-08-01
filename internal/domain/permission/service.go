package permission

import (
	"jcourse_go/internal/domain/auth"
	"jcourse_go/internal/domain/common"
)

type PermissionService interface {
	CheckPermission(commonCtx *common.CommonContext, ref ResourceRef, action Action) (Result, error)
}

type permissionService struct {
	userRepo auth.UserRepository
}

func (p *permissionService) CheckPermission(commonCtx *common.CommonContext, ref ResourceRef, action Action) (Result, error) {
	switch ref.Type {
	case ResourceTypeReview:
		return p.checkReviewPermission(commonCtx, ref, action)
	case ResourceTypeReviewAction:
		return p.checkReviewActionPermission(commonCtx, ref, action)
	case ResourceTypeUser:
		return p.checkUserPermission(commonCtx, ref, action)
	case ResourceTypePoint:
		return p.checkPointPermission(commonCtx, ref, action)
	case ResourceTypeCourse:
		return p.checkCoursePermission(commonCtx, ref, action)
	default:
		return Result{Allow: false, Reason: "unknown resource type"}, nil
	}
}

func (p *permissionService) checkReviewPermission(commonCtx *common.CommonContext, ref ResourceRef, action Action) (Result, error) {
	switch action {
	case ActionView:
		return Result{Allow: true, Reason: "public access"}, nil
	case ActionCreate:
		if commonCtx.User == nil {
			return Result{Allow: false, Reason: "not authenticated"}, nil
		}
		return Result{Allow: true, Reason: "authenticated user"}, nil
	case ActionUpdate, ActionDelete:
		if commonCtx.User == nil {
			return Result{Allow: false, Reason: "not authenticated"}, nil
		}
		if ref.Owner.ID == 0 {
			return Result{Allow: false, Reason: "resource owner not found"}, nil
		}

		// Check if user is admin or owner
		if commonCtx.User.Role == common.RoleAdmin {
			return Result{Allow: true, Reason: "admin access"}, nil
		}
		if commonCtx.User.UserID == ref.Owner.ID {
			return Result{Allow: true, Reason: "owner access"}, nil
		}

		return Result{Allow: false, Reason: "permission denied"}, nil
	default:
		return Result{Allow: false, Reason: "unknown action"}, nil
	}
}

func (p *permissionService) checkReviewActionPermission(commonCtx *common.CommonContext, ref ResourceRef, action Action) (Result, error) {
	switch action {
	case ActionCreate:
		if commonCtx.User == nil {
			return Result{Allow: false, Reason: "not authenticated"}, nil
		}
		return Result{Allow: true, Reason: "authenticated user"}, nil
	case ActionDelete:
		if commonCtx.User == nil {
			return Result{Allow: false, Reason: "not authenticated"}, nil
		}
		if ref.Owner.ID == 0 {
			return Result{Allow: false, Reason: "resource owner not found"}, nil
		}

		// Check if user is admin or owner
		if commonCtx.User.Role == common.RoleAdmin {
			return Result{Allow: true, Reason: "admin access"}, nil
		}
		if commonCtx.User.UserID == ref.Owner.ID {
			return Result{Allow: true, Reason: "owner access"}, nil
		}

		return Result{Allow: false, Reason: "permission denied"}, nil
	default:
		return Result{Allow: false, Reason: "unknown action"}, nil
	}
}

func (p *permissionService) checkUserPermission(commonCtx *common.CommonContext, ref ResourceRef, action Action) (Result, error) {
	switch action {
	case ActionView:
		return Result{Allow: true, Reason: "public access"}, nil
	case ActionUpdate:
		if commonCtx.User == nil {
			return Result{Allow: false, Reason: "not authenticated"}, nil
		}
		// Check if user is owner or admin
		if commonCtx.User.Role == common.RoleAdmin {
			return Result{Allow: true, Reason: "admin access"}, nil
		}
		if commonCtx.User.UserID == ref.Owner.ID {
			return Result{Allow: true, Reason: "owner access"}, nil
		}

		return Result{Allow: false, Reason: "permission denied"}, nil
	case ActionCreate, ActionDelete:
		return Result{Allow: false, Reason: "action not allowed"}, nil
	default:
		return Result{Allow: false, Reason: "unknown action"}, nil
	}
}

func (p *permissionService) checkPointPermission(commonCtx *common.CommonContext, _ ResourceRef, _ Action) (Result, error) {
	if commonCtx.User == nil {
		return Result{Allow: false, Reason: "not authenticated"}, nil
	}

	// Check if user is admin
	if commonCtx.User.Role == common.RoleAdmin {
		return Result{Allow: true, Reason: "admin access"}, nil
	}

	return Result{Allow: false, Reason: "admin access required"}, nil
}

func (p *permissionService) checkCoursePermission(commonCtx *common.CommonContext, _ ResourceRef, action Action) (Result, error) {
	switch action {
	case ActionView:
		return Result{Allow: true, Reason: "public access"}, nil
	case ActionCreate, ActionUpdate, ActionDelete:
		if commonCtx.User == nil {
			return Result{Allow: false, Reason: "not authenticated"}, nil
		}
		// Check if user is admin
		if commonCtx.User.Role == common.RoleAdmin {
			return Result{Allow: true, Reason: "admin access"}, nil
		}

		return Result{Allow: false, Reason: "admin access required"}, nil
	default:
		return Result{Allow: false, Reason: "unknown action"}, nil
	}
}

func NewPermissionService(userRepo auth.UserRepository) PermissionService {
	return &permissionService{
		userRepo: userRepo,
	}
}

// Helper functions to create ResourceRef objects
func NewReviewResourceRef(reviewID, ownerID int) ResourceRef {
	return ResourceRef{
		ID:   reviewID,
		Type: ResourceTypeReview,
		Owner: ResourceOwner{
			ID: ownerID,
		},
	}
}

func NewReviewActionResourceRef(actionID, ownerID int) ResourceRef {
	return ResourceRef{
		ID:   actionID,
		Type: ResourceTypeReviewAction,
		Owner: ResourceOwner{
			ID: ownerID,
		},
	}
}

func NewUserResourceRef(userID int) ResourceRef {
	return ResourceRef{
		ID:   userID,
		Type: ResourceTypeUser,
		Owner: ResourceOwner{
			ID: userID,
		},
	}
}

func NewPointResourceRef() ResourceRef {
	return ResourceRef{
		ID:   0,
		Type: ResourceTypePoint,
		Owner: ResourceOwner{
			ID: 0,
		},
	}
}

func NewCourseResourceRef() ResourceRef {
	return ResourceRef{
		ID:   0,
		Type: ResourceTypeCourse,
		Owner: ResourceOwner{
			ID: 0,
		},
	}
}
