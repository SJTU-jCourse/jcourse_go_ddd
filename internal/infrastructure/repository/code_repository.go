package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"jcourse_go/internal/domain/auth"
	"jcourse_go/internal/infrastructure/entity"
)

type codeRepository struct {
	db *gorm.DB
}

func NewCodeRepository(db *gorm.DB) auth.CodeRepository {
	return &codeRepository{db: db}
}

func (r *codeRepository) Get(ctx context.Context, email string) (*auth.VerificationCode, error) {
	var codeEntity entity.VerificationCode
	result := r.db.WithContext(ctx).Where("email = ?", email).First(&codeEntity)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get verification code: %w", result.Error)
	}
	return r.toDomainCode(&codeEntity), nil
}

func (r *codeRepository) Save(ctx context.Context, code *auth.VerificationCode) error {
	codeEntity := r.toORMCode(code)
	result := r.db.WithContext(ctx).Where("email = ?", code.Email).First(&entity.VerificationCode{})
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			result = r.db.WithContext(ctx).Create(codeEntity)
			if result.Error != nil {
				return fmt.Errorf("failed to create verification code: %w", result.Error)
			}
			return nil
		}
		return fmt.Errorf("failed to check existing verification code: %w", result.Error)
	}

	updateData := map[string]interface{}{
		"code":       code.Code,
		"expires_at": code.ExpiresAt,
		"created_at": time.Now(),
	}
	result = r.db.WithContext(ctx).Model(&entity.VerificationCode{}).Where("email = ?", code.Email).Updates(updateData)
	if result.Error != nil {
		return fmt.Errorf("failed to update verification code: %w", result.Error)
	}
	return nil
}

func (r *codeRepository) Delete(ctx context.Context, code *auth.VerificationCode) error {
	result := r.db.WithContext(ctx).Where("email = ?", code.Email).Delete(&entity.VerificationCode{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete verification code: %w", result.Error)
	}
	return nil
}

// Helper methods to convert between domain and ORM models
func (r *codeRepository) toDomainCode(codeEntity *entity.VerificationCode) *auth.VerificationCode {
	return &auth.VerificationCode{
		Email:     codeEntity.Email,
		Code:      codeEntity.Code,
		ExpiresAt: codeEntity.ExpiresAt,
		CreatedAt: codeEntity.CreatedAt,
	}
}

func (r *codeRepository) toORMCode(code *auth.VerificationCode) *entity.VerificationCode {
	return &entity.VerificationCode{
		Email:     code.Email,
		Code:      code.Code,
		ExpiresAt: code.ExpiresAt,
		CreatedAt: code.CreatedAt,
	}
}
