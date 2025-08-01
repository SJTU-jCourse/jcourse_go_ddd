package email

import (
	"context"
	"fmt"

	"jcourse_go/internal/domain/auth"
)

// MockEmailService is a mock implementation of EmailService for testing
type MockEmailService struct{}

func NewMockEmailService() *MockEmailService {
	return &MockEmailService{}
}

func (s *MockEmailService) SendVerificationCode(ctx context.Context, emailAddr string, input *auth.VerificationCode) error {
	fmt.Printf("Mock email sent to %s with code %s\n", emailAddr, input.Code)
	return nil
}
