package permission

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"jcourse_go/internal/domain/common"
)

func TestPermissionService_CheckReviewPermission(t *testing.T) {
	userRepo := NewMockUserRepository()
	permissionService := NewPermissionService(userRepo)

	tests := []struct {
		name     string
		ref      ResourceRef
		action   Action
		userID   int
		role     common.Role
		expected Result
	}{
		{
			name:     "anonymous user can view review",
			ref:      ResourceRef{Type: ResourceTypeReview, Owner: ResourceOwner{ID: 2}},
			action:   ActionView,
			userID:   0,
			role:     "",
			expected: Result{Allow: true, Reason: "public access"},
		},
		{
			name:     "authenticated user can create review",
			ref:      ResourceRef{Type: ResourceTypeReview, Owner: ResourceOwner{ID: 0}},
			action:   ActionCreate,
			userID:   2,
			role:     common.RoleUser,
			expected: Result{Allow: true, Reason: "authenticated user"},
		},
		{
			name:     "owner can update review",
			ref:      ResourceRef{Type: ResourceTypeReview, Owner: ResourceOwner{ID: 2}},
			action:   ActionUpdate,
			userID:   2,
			role:     common.RoleUser,
			expected: Result{Allow: true, Reason: "owner access"},
		},
		{
			name:     "admin can update review",
			ref:      ResourceRef{Type: ResourceTypeReview, Owner: ResourceOwner{ID: 2}},
			action:   ActionUpdate,
			userID:   1,
			role:     common.RoleAdmin,
			expected: Result{Allow: true, Reason: "admin access"},
		},
		{
			name:     "non-owner cannot update review",
			ref:      ResourceRef{Type: ResourceTypeReview, Owner: ResourceOwner{ID: 2}},
			action:   ActionUpdate,
			userID:   3,
			role:     common.RoleUser,
			expected: Result{Allow: false, Reason: "permission denied"},
		},
		{
			name:     "anonymous user cannot create review",
			ref:      ResourceRef{Type: ResourceTypeReview, Owner: ResourceOwner{ID: 0}},
			action:   ActionCreate,
			userID:   0,
			role:     "",
			expected: Result{Allow: false, Reason: "not authenticated"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			var user *common.User
			if tt.userID != 0 {
				user = &common.User{UserID: tt.userID, Role: tt.role}
			}
			commonCtx := common.NewCommonContext(ctx, user)
			result, err := permissionService.CheckPermission(commonCtx, tt.ref, tt.action)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestPermissionService_CheckUserPermission(t *testing.T) {
	userRepo := NewMockUserRepository()
	permissionService := NewPermissionService(userRepo)

	tests := []struct {
		name     string
		ref      ResourceRef
		action   Action
		userID   int
		role     common.Role
		expected Result
	}{
		{
			name:     "anyone can view user",
			ref:      ResourceRef{Type: ResourceTypeUser, Owner: ResourceOwner{ID: 2}},
			action:   ActionView,
			userID:   0,
			role:     "",
			expected: Result{Allow: true, Reason: "public access"},
		},
		{
			name:     "owner can update user",
			ref:      ResourceRef{Type: ResourceTypeUser, Owner: ResourceOwner{ID: 2}},
			action:   ActionUpdate,
			userID:   2,
			role:     common.RoleUser,
			expected: Result{Allow: true, Reason: "owner access"},
		},
		{
			name:     "admin can update user",
			ref:      ResourceRef{Type: ResourceTypeUser, Owner: ResourceOwner{ID: 2}},
			action:   ActionUpdate,
			userID:   1,
			role:     common.RoleAdmin,
			expected: Result{Allow: true, Reason: "admin access"},
		},
		{
			name:     "non-owner cannot update user",
			ref:      ResourceRef{Type: ResourceTypeUser, Owner: ResourceOwner{ID: 2}},
			action:   ActionUpdate,
			userID:   3,
			role:     common.RoleUser,
			expected: Result{Allow: false, Reason: "permission denied"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			var user *common.User
			if tt.userID != 0 {
				user = &common.User{UserID: tt.userID, Role: tt.role}
			}
			commonCtx := common.NewCommonContext(ctx, user)
			result, err := permissionService.CheckPermission(commonCtx, tt.ref, tt.action)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestPermissionService_CheckPointPermission(t *testing.T) {
	userRepo := NewMockUserRepository()
	permissionService := NewPermissionService(userRepo)

	tests := []struct {
		name     string
		ref      ResourceRef
		action   Action
		userID   int
		role     common.Role
		expected Result
	}{
		{
			name:     "admin can manage points",
			ref:      ResourceRef{Type: ResourceTypePoint, Owner: ResourceOwner{ID: 0}},
			action:   ActionCreate,
			userID:   1,
			role:     common.RoleAdmin,
			expected: Result{Allow: true, Reason: "admin access"},
		},
		{
			name:     "non-admin cannot manage points",
			ref:      ResourceRef{Type: ResourceTypePoint, Owner: ResourceOwner{ID: 0}},
			action:   ActionCreate,
			userID:   2,
			role:     common.RoleUser,
			expected: Result{Allow: false, Reason: "admin access required"},
		},
		{
			name:     "anonymous cannot manage points",
			ref:      ResourceRef{Type: ResourceTypePoint, Owner: ResourceOwner{ID: 0}},
			action:   ActionCreate,
			userID:   0,
			role:     "",
			expected: Result{Allow: false, Reason: "not authenticated"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			var user *common.User
			if tt.userID != 0 {
				user = &common.User{UserID: tt.userID, Role: tt.role}
			}
			commonCtx := common.NewCommonContext(ctx, user)
			result, err := permissionService.CheckPermission(commonCtx, tt.ref, tt.action)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestPermissionService_CheckCoursePermission(t *testing.T) {
	userRepo := NewMockUserRepository()
	permissionService := NewPermissionService(userRepo)

	tests := []struct {
		name     string
		ref      ResourceRef
		action   Action
		userID   int
		role     common.Role
		expected Result
	}{
		{
			name:     "anyone can view course",
			ref:      ResourceRef{Type: ResourceTypeCourse, Owner: ResourceOwner{ID: 0}},
			action:   ActionView,
			userID:   0,
			role:     "",
			expected: Result{Allow: true, Reason: "public access"},
		},
		{
			name:     "admin can create course",
			ref:      ResourceRef{Type: ResourceTypeCourse, Owner: ResourceOwner{ID: 0}},
			action:   ActionCreate,
			userID:   1,
			role:     common.RoleAdmin,
			expected: Result{Allow: true, Reason: "admin access"},
		},
		{
			name:     "non-admin cannot create course",
			ref:      ResourceRef{Type: ResourceTypeCourse, Owner: ResourceOwner{ID: 0}},
			action:   ActionCreate,
			userID:   2,
			role:     common.RoleUser,
			expected: Result{Allow: false, Reason: "admin access required"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			var user *common.User
			if tt.userID != 0 {
				user = &common.User{UserID: tt.userID, Role: tt.role}
			}
			commonCtx := common.NewCommonContext(ctx, user)
			result, err := permissionService.CheckPermission(commonCtx, tt.ref, tt.action)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}
