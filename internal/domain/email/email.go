package email

import (
	"context"
	"fmt"

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

// EmailServiceImpl implements EmailService interface
type EmailServiceImpl struct {
	sender    Sender
	template Template
}

func NewEmailService() EmailService {
	return nil
}

// NewEmailServiceImpl creates a new email service with sender and template
func NewEmailServiceImpl(sender Sender, template Template) EmailService {
	return &EmailServiceImpl{
		sender:    sender,
		template: template,
	}
}

func (s *EmailServiceImpl) SendVerificationCode(ctx context.Context, emailAddr string, input *auth.VerificationCode) error {
	if s.sender == nil || s.template == nil {
		return fmt.Errorf("email service not properly initialized")
	}

	renderedEmail := s.template.Execute(ctx, input)
	return s.sender.Send(ctx, emailAddr, renderedEmail)
}
