package permission

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"jcourse_go/internal/domain/auth"
	"jcourse_go/internal/domain/common"
)

// MockUserRepository for testing
type testMockUserRepository struct {
	users map[int]*auth.User
}

func (m *testMockUserRepository) Get(ctx context.Context, email string) (*auth.User, error) {
	for _, user := range m.users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, nil
}

func (m *testMockUserRepository) GetByID(ctx context.Context, userID int) (*auth.User, error) {
	user, exists := m.users[userID]
	if !exists {
		return nil, nil
	}
	return user, nil
}

func (m *testMockUserRepository) FindBy(ctx context.Context, filter auth.UserFilter) ([]auth.User, error) {
	var users []auth.User
	for _, user := range m.users {
		users = append(users, *user)
	}
	return users, nil
}

func (m *testMockUserRepository) Save(ctx context.Context, user *auth.User) (int, error) {
	if user.ID == 0 {
		user.ID = len(m.users) + 1
	}
	m.users[user.ID] = user
	return user.ID, nil
}

func (m *testMockUserRepository) Update(ctx context.Context, user *auth.User) error {
	m.users[user.ID] = user
	return nil
}

func createTestMockUserRepository() *testMockUserRepository {
	return &testMockUserRepository{
		users: map[int]*auth.User{
			1: {ID: 1, Username: "admin", Email: "admin@test.com", Role: common.RoleAdmin},
			2: {ID: 2, Username: "user1", Email: "user1@test.com", Role: common.RoleUser},
			3: {ID: 3, Username: "user2", Email: "user2@test.com", Role: common.RoleUser},
		},
	}
}

func TestPermissionService_CheckReviewPermission(t *testing.T) {
	userRepo := createTestMockUserRepository()
	permissionService := NewPermissionServiceWithUserRepo(userRepo)

	tests := []struct {
		name     string
		ref      ResourceRef
		action   Action
		userID   int
		expected Result
	}{
		{
			name:     "anonymous user can view review",
			ref:      ResourceRef{Type: ResourceTypeReview, Owner: ResourceOwner{ID: 2}},
			action:   ActionView,
			userID:   0,
			expected: Result{Allow: true, Reason: "public access"},
		},
		{
			name:     "authenticated user can create review",
			ref:      ResourceRef{Type: ResourceTypeReview, Owner: ResourceOwner{ID: 0}},
			action:   ActionCreate,
			userID:   2,
			expected: Result{Allow: true, Reason: "authenticated user"},
		},
		{
			name:     "owner can update review",
			ref:      ResourceRef{Type: ResourceTypeReview, Owner: ResourceOwner{ID: 2}},
			action:   ActionUpdate,
			userID:   2,
			expected: Result{Allow: true, Reason: "owner access"},
		},
		{
			name:     "admin can update review",
			ref:      ResourceRef{Type: ResourceTypeReview, Owner: ResourceOwner{ID: 2}},
			action:   ActionUpdate,
			userID:   1,
			expected: Result{Allow: true, Reason: "admin access"},
		},
		{
			name:     "non-owner cannot update review",
			ref:      ResourceRef{Type: ResourceTypeReview, Owner: ResourceOwner{ID: 2}},
			action:   ActionUpdate,
			userID:   3,
			expected: Result{Allow: false, Reason: "permission denied"},
		},
		{
			name:     "anonymous user cannot create review",
			ref:      ResourceRef{Type: ResourceTypeReview, Owner: ResourceOwner{ID: 0}},
			action:   ActionCreate,
			userID:   0,
			expected: Result{Allow: false, Reason: "not authenticated"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := permissionService.CheckPermission(context.Background(), tt.ref, tt.action, tt.userID)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestPermissionService_CheckUserPermission(t *testing.T) {
	userRepo := createTestMockUserRepository()
	permissionService := NewPermissionServiceWithUserRepo(userRepo)

	tests := []struct {
		name     string
		ref      ResourceRef
		action   Action
		userID   int
		expected Result
	}{
		{
			name:     "anyone can view user",
			ref:      ResourceRef{Type: ResourceTypeUser, Owner: ResourceOwner{ID: 2}},
			action:   ActionView,
			userID:   0,
			expected: Result{Allow: true, Reason: "public access"},
		},
		{
			name:     "owner can update user",
			ref:      ResourceRef{Type: ResourceTypeUser, Owner: ResourceOwner{ID: 2}},
			action:   ActionUpdate,
			userID:   2,
			expected: Result{Allow: true, Reason: "owner access"},
		},
		{
			name:     "admin can update user",
			ref:      ResourceRef{Type: ResourceTypeUser, Owner: ResourceOwner{ID: 2}},
			action:   ActionUpdate,
			userID:   1,
			expected: Result{Allow: true, Reason: "admin access"},
		},
		{
			name:     "non-owner cannot update user",
			ref:      ResourceRef{Type: ResourceTypeUser, Owner: ResourceOwner{ID: 2}},
			action:   ActionUpdate,
			userID:   3,
			expected: Result{Allow: false, Reason: "permission denied"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := permissionService.CheckPermission(context.Background(), tt.ref, tt.action, tt.userID)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestPermissionService_CheckPointPermission(t *testing.T) {
	userRepo := createTestMockUserRepository()
	permissionService := NewPermissionServiceWithUserRepo(userRepo)

	tests := []struct {
		name     string
		ref      ResourceRef
		action   Action
		userID   int
		expected Result
	}{
		{
			name:     "admin can manage points",
			ref:      ResourceRef{Type: ResourceTypePoint, Owner: ResourceOwner{ID: 0}},
			action:   ActionCreate,
			userID:   1,
			expected: Result{Allow: true, Reason: "admin access"},
		},
		{
			name:     "non-admin cannot manage points",
			ref:      ResourceRef{Type: ResourceTypePoint, Owner: ResourceOwner{ID: 0}},
			action:   ActionCreate,
			userID:   2,
			expected: Result{Allow: false, Reason: "admin access required"},
		},
		{
			name:     "anonymous cannot manage points",
			ref:      ResourceRef{Type: ResourceTypePoint, Owner: ResourceOwner{ID: 0}},
			action:   ActionCreate,
			userID:   0,
			expected: Result{Allow: false, Reason: "admin access required"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := permissionService.CheckPermission(context.Background(), tt.ref, tt.action, tt.userID)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestPermissionService_CheckCoursePermission(t *testing.T) {
	userRepo := createTestMockUserRepository()
	permissionService := NewPermissionServiceWithUserRepo(userRepo)

	tests := []struct {
		name     string
		ref      ResourceRef
		action   Action
		userID   int
		expected Result
	}{
		{
			name:     "anyone can view course",
			ref:      ResourceRef{Type: ResourceTypeCourse, Owner: ResourceOwner{ID: 0}},
			action:   ActionView,
			userID:   0,
			expected: Result{Allow: true, Reason: "public access"},
		},
		{
			name:     "admin can create course",
			ref:      ResourceRef{Type: ResourceTypeCourse, Owner: ResourceOwner{ID: 0}},
			action:   ActionCreate,
			userID:   1,
			expected: Result{Allow: true, Reason: "admin access"},
		},
		{
			name:     "non-admin cannot create course",
			ref:      ResourceRef{Type: ResourceTypeCourse, Owner: ResourceOwner{ID: 0}},
			action:   ActionCreate,
			userID:   2,
			expected: Result{Allow: false, Reason: "admin access required"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := permissionService.CheckPermission(context.Background(), tt.ref, tt.action, tt.userID)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}
