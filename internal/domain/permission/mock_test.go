package permission

import (
	"context"

	"jcourse_go/internal/domain/auth"
	"jcourse_go/internal/domain/common"
)

// MockUserRepository is a mock implementation of auth.UserRepository for testing
type MockUserRepository struct {
	users map[int]*auth.User
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		users: map[int]*auth.User{
			1: {ID: 1, Username: "admin", Email: "admin@test.com", Role: common.RoleAdmin},
			2: {ID: 2, Username: "user1", Email: "user1@test.com", Role: common.RoleUser},
			3: {ID: 3, Username: "user2", Email: "user2@test.com", Role: common.RoleUser},
		},
	}
}

func (m *MockUserRepository) Get(ctx context.Context, email string) (*auth.User, error) {
	for _, user := range m.users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, nil
}

func (m *MockUserRepository) GetByID(ctx context.Context, userID int) (*auth.User, error) {
	user, exists := m.users[userID]
	if !exists {
		return nil, nil
	}
	return user, nil
}

func (m *MockUserRepository) FindBy(ctx context.Context, filter auth.UserFilter) ([]auth.User, error) {
	var users []auth.User
	for _, user := range m.users {
		users = append(users, *user)
	}
	return users, nil
}

func (m *MockUserRepository) Save(ctx context.Context, user *auth.User) (int, error) {
	if user.ID == 0 {
		user.ID = len(m.users) + 1
	}
	m.users[user.ID] = user
	return user.ID, nil
}

func (m *MockUserRepository) Update(ctx context.Context, user *auth.User) error {
	m.users[user.ID] = user
	return nil
}
