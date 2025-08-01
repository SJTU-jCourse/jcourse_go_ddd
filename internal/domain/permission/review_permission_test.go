package permission

import (
	"context"
	"testing"

	"jcourse_go/internal/domain/auth"
	"jcourse_go/internal/domain/common"
	"jcourse_go/internal/domain/review"
)

// MockUserRepository for testing
type mockUserRepository struct{}

func (m *mockUserRepository) Get(ctx context.Context, email string) (*auth.User, error) {
	return nil, nil
}

func (m *mockUserRepository) GetByID(ctx context.Context, userID int) (*auth.User, error) {
	return &auth.User{ID: userID, Role: common.RoleUser}, nil
}

func (m *mockUserRepository) FindBy(ctx context.Context, filter auth.UserFilter) ([]auth.User, error) {
	return nil, nil
}

func (m *mockUserRepository) Save(ctx context.Context, user *auth.User) (int, error) {
	return user.ID, nil
}

func (m *mockUserRepository) Update(ctx context.Context, user *auth.User) error {
	return nil
}

func TestReviewPermissionChecker_CanUpdateReview(t *testing.T) {
	mockRepo := &mockUserRepository{}
	service := NewReviewPermissionChecker(mockRepo)

	tests := []struct {
		name     string
		review   *Review
		user     *common.User
		expected bool
		reason   string
	}{
		{
			name: "admin can update any review",
			review: &Review{
				ID:     1,
				UserID: 2,
			},
			user: &common.User{
				UserID: 1,
				Role:   common.RoleAdmin,
			},
			expected: true,
			reason:   "admin access",
		},
		{
			name: "owner can update own review",
			review: &Review{
				ID:     1,
				UserID: 1,
			},
			user: &common.User{
				UserID: 1,
				Role:   common.RoleUser,
			},
			expected: true,
			reason:   "owner access",
		},
		{
			name: "user cannot update others review",
			review: &Review{
				ID:     1,
				UserID: 2,
			},
			user: &common.User{
				UserID: 1,
				Role:   common.RoleUser,
			},
			expected: false,
			reason:   "permission denied",
		},
		{
			name:     "unauthenticated user cannot update",
			review:   &Review{ID: 1, UserID: 2},
			user:     nil,
			expected: false,
			reason:   "not authenticated",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			canUpdate, reason := service.CanUpdateReview(ctx, tt.review, tt.user)
			
			if canUpdate != tt.expected {
				t.Errorf("CanUpdateReview() = %v, want %v", canUpdate, tt.expected)
			}
			if reason != tt.reason {
				t.Errorf("CanUpdateReview() reason = %v, want %v", reason, tt.reason)
			}
		})
	}
}

func TestReviewPermissionChecker_CanDeleteReview(t *testing.T) {
	mockRepo := &mockUserRepository{}
	service := NewReviewPermissionChecker(mockRepo)

	tests := []struct {
		name     string
		review   *Review
		user     *common.User
		expected bool
		reason   string
	}{
		{
			name: "admin can delete any review",
			review: &Review{
				ID:     1,
				UserID: 2,
			},
			user: &common.User{
				UserID: 1,
				Role:   common.RoleAdmin,
			},
			expected: true,
			reason:   "admin access",
		},
		{
			name: "owner can delete own review",
			review: &Review{
				ID:     1,
				UserID: 1,
			},
			user: &common.User{
				UserID: 1,
				Role:   common.RoleUser,
			},
			expected: true,
			reason:   "owner access",
		},
		{
			name: "user cannot delete others review",
			review: &Review{
				ID:     1,
				UserID: 2,
			},
			user: &common.User{
				UserID: 1,
				Role:   common.RoleUser,
			},
			expected: false,
			reason:   "permission denied",
		},
		{
			name:     "unauthenticated user cannot delete",
			review:   &Review{ID: 1, UserID: 2},
			user:     nil,
			expected: false,
			reason:   "not authenticated",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			canDelete, reason := service.CanDeleteReview(ctx, tt.review, tt.user)
			
			if canDelete != tt.expected {
				t.Errorf("CanDeleteReview() = %v, want %v", canDelete, tt.expected)
			}
			if reason != tt.reason {
				t.Errorf("CanDeleteReview() reason = %v, want %v", reason, tt.reason)
			}
		})
	}
}

func TestToPermissionReview(t *testing.T) {
	domainReview := &review.Review{
		ID:     123,
		UserID: 456,
	}
	
	permissionReview := ToPermissionReview(domainReview)
	
	if permissionReview.ID != domainReview.ID {
		t.Errorf("ToPermissionReview() ID = %v, want %v", permissionReview.ID, domainReview.ID)
	}
	if permissionReview.UserID != domainReview.UserID {
		t.Errorf("ToPermissionReview() UserID = %v, want %v", permissionReview.UserID, domainReview.UserID)
	}
}