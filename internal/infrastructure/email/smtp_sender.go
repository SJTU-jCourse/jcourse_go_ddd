package email

import (
	"context"
	"fmt"
	"time"

	"gopkg.in/gomail.v2"

	"jcourse_go/internal/config"
	"jcourse_go/internal/domain/email"
)

type SMTPSender struct {
	config config.SMTPConfig
	dialer *gomail.Dialer
}

func NewSMTPSender(config config.SMTPConfig) *SMTPSender {
	return &SMTPSender{
		config: config,
		dialer: gomail.NewDialer(
			config.Host,
			config.Port,
			config.Username,
			config.Password,
		),
	}
}

func (s *SMTPSender) Send(ctx context.Context, emailAddr string, email email.RenderedEmail) error {
	if s.config.Host == "" || s.config.Username == "" || s.config.Password == "" {
		return fmt.Errorf("SMTP configuration is incomplete")
	}

	// Create message
	m := gomail.NewMessage()
	m.SetHeader("From", s.config.Sender)
	m.SetHeader("To", emailAddr)
	m.SetHeader("Subject", email.Title)
	m.SetBody("text/html", email.Body)

	// Send with timeout
	sendCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// Use channel to handle the blocking send operation
	resultChan := make(chan error, 1)
	
	go func() {
		resultChan <- s.dialer.DialAndSend(m)
	}()

	select {
	case err := <-resultChan:
		if err != nil {
			return fmt.Errorf("failed to send email: %w", err)
		}
		return nil
	case <-sendCtx.Done():
		return fmt.Errorf("email send timed out")
	}
}