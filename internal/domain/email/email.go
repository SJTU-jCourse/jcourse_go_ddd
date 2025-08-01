package email

import (
	"context"

	"jcourse_go/internal/domain/auth"
)

type Template interface {
	Execute(ctx context.Context, input *auth.VerificationCode) RenderedEmail
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

// NewEmailService creates a new email service instance
// For now, returns a mock implementation that logs to console
func NewEmailService() EmailService {
	return nil
}
