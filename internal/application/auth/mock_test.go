package auth

import (
	"context"

	"jcourse_go/internal/domain/auth"
)

// MockEmailService is a mock implementation of email.EmailService for testing
type MockEmailService struct {
	SendError error
}

func (m *MockEmailService) SendVerificationCode(ctx context.Context, emailAddr string, input *auth.VerificationCode) error {
	return m.SendError
}

// MockCodeRepository is a mock implementation of auth.CodeRepository for testing
type MockCodeRepository struct {
	GetCode     *auth.VerificationCode
	GetError    error
	SaveError   error
	DeleteError error
}

func (m *MockCodeRepository) Get(ctx context.Context, email string) (*auth.VerificationCode, error) {
	return m.GetCode, m.GetError
}

func (m *MockCodeRepository) Save(ctx context.Context, code *auth.VerificationCode) error {
	return m.SaveError
}

func (m *MockCodeRepository) Delete(ctx context.Context, code *auth.VerificationCode) error {
	return m.DeleteError
}
