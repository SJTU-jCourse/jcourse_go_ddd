package auth

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"jcourse_go/internal/domain/auth"
	"jcourse_go/internal/domain/common"
	"jcourse_go/internal/infrastructure/entity"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) auth.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Get(ctx context.Context, email string) (*auth.User, error) {
	var userEntity entity.User
	result := r.db.WithContext(ctx).Where("email = ?", email).First(&userEntity)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user by email: %w", result.Error)
	}
	return r.toDomainUser(&userEntity), nil
}

func (r *userRepository) GetByID(ctx context.Context, userID int) (*auth.User, error) {
	var userEntity entity.User
	result := r.db.WithContext(ctx).First(&userEntity, userID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user by id: %w", result.Error)
	}
	return r.toDomainUser(&userEntity), nil
}

func (r *userRepository) FindBy(ctx context.Context, filter auth.UserFilter) ([]auth.User, error) {
	var userEntitys []entity.User
	query := r.db.WithContext(ctx)

	if len(filter.UserIDs) > 0 {
		query = query.Where("id IN ?", filter.UserIDs)
	}

	result := query.Find(&userEntitys)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to find users: %w", result.Error)
	}

	users := make([]auth.User, len(userEntitys))
	for i, userEntity := range userEntitys {
		users[i] = *r.toDomainUser(&userEntity)
	}

	return users, nil
}

func (r *userRepository) Save(ctx context.Context, user *auth.User) (int, error) {
	userEntity := r.toORMUser(user)
	result := r.db.WithContext(ctx).Create(userEntity)
	if result.Error != nil {
		return 0, fmt.Errorf("failed to save user: %w", result.Error)
	}
	return userEntity.ID, nil
}

func (r *userRepository) Update(ctx context.Context, user *auth.User) error {
	userEntity := r.toORMUser(user)
	result := r.db.WithContext(ctx).Save(userEntity)
	if result.Error != nil {
		return fmt.Errorf("failed to update user: %w", result.Error)
	}
	return nil
}

// Helper methods to convert between domain and ORM models
func (r *userRepository) toDomainUser(userEntity *entity.User) *auth.User {
	return &auth.User{
		ID:       userEntity.ID,
		Username: userEntity.Username,
		Password: userEntity.PasswordHash,
		Email:    userEntity.Email,
		Role:     common.Role(userEntity.Role),
	}
}

func (r *userRepository) toORMUser(user *auth.User) *entity.User {
	return &entity.User{
		ID:           user.ID,
		Username:     user.Username,
		Email:        user.Email,
		PasswordHash: user.Password,
		Role:         string(user.Role),
		IsVerified:   true, // Default to true since domain doesn't have this field
	}
}
