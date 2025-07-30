package password

type Hasher interface {
	Hash(password string) string
	Validate(password, hash string) error
}
