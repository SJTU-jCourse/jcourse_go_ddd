package password

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
)

type Hasher interface {
	Hash(password string) string
	Validate(password, hash string) error
}

type bcryptHasher struct{}

func NewHasher() Hasher {
	return &bcryptHasher{}
}

func (h *bcryptHasher) Hash(password string) string {
	hasher := sha256.New()
	hasher.Write([]byte(password))
	return hex.EncodeToString(hasher.Sum(nil))
}

func (h *bcryptHasher) Validate(password, hash string) error {
	hashedPassword := h.Hash(password)
	if hashedPassword != hash {
		return errors.New("invalid password")
	}
	return nil
}
