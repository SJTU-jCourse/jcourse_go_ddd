package auth

import (
	"context"
	"errors"
	"testing"
	"time"

	"jcourse_go/internal/domain/auth"
	"jcourse_go/pkg/apperror"

	"github.com/stretchr/testify/assert"
)

// Mock implementations for the dependencies
type mockEmailService struct {
	sendError error
}

func (m *mockEmailService) SendVerificationCode(ctx context.Context, emailAddr string, input *auth.VerificationCode) error {
	return m.sendError
}

type mockCodeRepository struct {
	getCode     *auth.VerificationCode
	getError    error
	saveError   error
	deleteError error
}

func (m *mockCodeRepository) Get(ctx context.Context, email string) (*auth.VerificationCode, error) {
	return m.getCode, m.getError
}

func (m *mockCodeRepository) Save(ctx context.Context, code *auth.VerificationCode) error {
	return m.saveError
}

func (m *mockCodeRepository) Delete(ctx context.Context, code *auth.VerificationCode) error {
	return m.deleteError
}

// Tests
func TestVerificationCodeService_SendCode(t *testing.T) {
	tests := []struct {
		name          string
		setupMocks    func() *verificationCodeService
		expectedError error
	}{
		{
			name: "success",
			setupMocks: func() *verificationCodeService {
				mockRepo := &mockCodeRepository{}
				mockEmail := &mockEmailService{}
				return &verificationCodeService{
					email:       mockEmail,
					codeRepo:    mockRepo,
					codeLength:  6,
					codeCharset: "0123456789",
					ttl:         10 * time.Minute,
					interval:    1 * time.Minute,
				}
			},
			expectedError: nil,
		},
		{
			name: "cannot send rate limit error",
			setupMocks: func() *verificationCodeService {
				mockRepo := &mockCodeRepository{
					getCode: &auth.VerificationCode{CreatedAt: time.Now()},
				}
				mockEmail := &mockEmailService{sendError: nil}
				return &verificationCodeService{
					email:       mockEmail,
					codeRepo:    mockRepo,
					codeLength:  6,
					codeCharset: "0123456789",
					ttl:         10 * time.Minute,
					interval:    1 * time.Minute,
				}
			},
			expectedError: apperror.ErrRateLimit,
		},
		{
			name: "email sending error",
			setupMocks: func() *verificationCodeService {
				mockRepo := &mockCodeRepository{
					getCode: nil,
				}
				mockEmail := &mockEmailService{sendError: errors.New("email error")}
				return &verificationCodeService{
					email:       mockEmail,
					codeRepo:    mockRepo,
					codeLength:  6,
					codeCharset: "0123456789",
					ttl:         10 * time.Minute,
					interval:    1 * time.Minute,
				}
			},
			expectedError: errors.New("email error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := tt.setupMocks()
			err := service.SendCode(context.Background(), "test@example.com")
			assert.Equal(t, tt.expectedError, err)
		})
	}
}

func TestVerificationCodeService_Verify(t *testing.T) {
	tests := []struct {
		name          string
		setupMocks    func() *verificationCodeService
		inputCode     string
		expectedError error
	}{
		{
			name: "success",
			setupMocks: func() *verificationCodeService {
				mockRepo := &mockCodeRepository{
					getCode: &auth.VerificationCode{Code: "123456", ExpiresAt: time.Now().Add(5 * time.Minute)},
				}
				return &verificationCodeService{
					email:       nil,
					codeRepo:    mockRepo,
					codeLength:  6,
					codeCharset: "0123456789",
					ttl:         10 * time.Minute,
					interval:    1 * time.Minute,
				}
			},
			inputCode:     "123456",
			expectedError: nil,
		},
		{
			name: "code doesn't match",
			setupMocks: func() *verificationCodeService {
				mockRepo := &mockCodeRepository{
					getCode: &auth.VerificationCode{Code: "123456", ExpiresAt: time.Now().Add(5 * time.Minute)},
				}
				return &verificationCodeService{
					email:       nil,
					codeRepo:    mockRepo,
					codeLength:  6,
					codeCharset: "0123456789",
					ttl:         10 * time.Minute,
					interval:    1 * time.Minute,
				}
			},
			inputCode:     "654321",
			expectedError: apperror.ErrWrongInput,
		},
		{
			name: "code expired",
			setupMocks: func() *verificationCodeService {
				mockRepo := &mockCodeRepository{
					getCode: &auth.VerificationCode{Code: "123456", ExpiresAt: time.Now().Add(-1 * time.Minute)},
				}
				return &verificationCodeService{
					email:       nil,
					codeRepo:    mockRepo,
					codeLength:  6,
					codeCharset: "0123456789",
					ttl:         10 * time.Minute,
					interval:    1 * time.Minute,
				}
			},
			inputCode:     "123456",
			expectedError: apperror.ErrExpired,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := tt.setupMocks()
			err := service.Verify(context.Background(), tt.inputCode, "test@example.com")
			assert.Equal(t, tt.expectedError, err)
		})
	}
}
