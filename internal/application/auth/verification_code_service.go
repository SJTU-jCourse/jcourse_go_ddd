package auth

import (
	"context"
	"crypto/rand"
	"math/big"
	"time"

	"jcourse_go/internal/domain/auth"
	"jcourse_go/internal/domain/email"
	"jcourse_go/pkg/apperror"
)

type VerificationCodeService interface {
	SendCode(ctx context.Context, email string) error
	Verify(ctx context.Context, inputCode string, email string) error
}

func NewVerificationCodeService(email email.EmailService, codeRepo auth.CodeRepository) VerificationCodeService {
	return &verificationCodeService{
		email:       email,
		codeRepo:    codeRepo,
		codeLength:  6,
		codeCharset: "0123456789",
		ttl:         10 * time.Minute,
		interval:    1 * time.Minute,
	}
}

type verificationCodeService struct {
	email       email.EmailService
	codeRepo    auth.CodeRepository
	ttl         time.Duration
	interval    time.Duration
	codeLength  int
	codeCharset string
}

func (v *verificationCodeService) Verify(ctx context.Context, inputCode string, email string) error {
	code, err := v.codeRepo.Get(ctx, email)
	if err != nil {
		return err
	}
	if code.IsExpired(time.Now()) {
		return apperror.ErrExpired
	}
	if code.Code != inputCode {
		return apperror.ErrWrongInput
	}
	err = v.codeRepo.Delete(ctx, code)
	return err
}

func (v *verificationCodeService) createCode(ctx context.Context, email string) *auth.VerificationCode {
	now := time.Now()
	code := &auth.VerificationCode{
		Code:      GenerateRandomCode(v.codeCharset, v.codeLength),
		Email:     email,
		ExpiresAt: now.Add(v.ttl),
		CreatedAt: now,
	}
	return code
}

func (v *verificationCodeService) canSendCode(ctx context.Context, email string) error {
	code, err := v.codeRepo.Get(ctx, email)
	if err != nil {
		return err
	}
	// If no existing code, allow sending
	if code == nil {
		return nil
	}
	now := time.Now()
	if now.Sub(code.CreatedAt) < v.interval {
		return apperror.ErrRateLimit
	}
	return nil
}

func (v *verificationCodeService) SendCode(ctx context.Context, email string) error {
	if err := v.canSendCode(ctx, email); err != nil {
		return err
	}
	code := v.createCode(ctx, email)
	if err := v.codeRepo.Save(ctx, code); err != nil {
		return err
	}
	return v.email.SendVerificationCode(ctx, code.Email, code)
}

func GenerateRandomCode(charSet string, length int) string {
	number := make([]byte, length)
	maxIdx := int64(len(charSet))
	for i := 0; i < length; i++ {
		n, _ := rand.Int(rand.Reader, big.NewInt(maxIdx))
		number[i] = charSet[n.Int64()]
	}
	return string(number)
}
