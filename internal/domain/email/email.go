package email

import (
	"context"
	"fmt"

	"jcourse_go/internal/domain/auth"
)

type Template interface {
	Execute(ctx context.Context, input any) RenderedEmail
}

type RenderedEmail struct {
	Title string
	Body  string
}

type Sender interface {
	Send(ctx context.Context, emailAddr string, email RenderedEmail) error
}

type EmailService interface {
	SendVerificationCode(ctx context.Context, emailAddr string, input *auth.VerificationCode) error
}

type mockEmailService struct{}

func NewEmailService() EmailService {
	return &mockEmailService{}
}

func (s *mockEmailService) SendVerificationCode(ctx context.Context, emailAddr string, input *auth.VerificationCode) error {
	fmt.Printf("Mock email sent to %s with code %s\n", emailAddr, input.Code)
	return nil
}
