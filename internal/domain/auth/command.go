package auth

type LoginCommand struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterCommand struct {
	Email    string `json:"email"`
	Code     string `json:"code"`
	Password string `json:"password"`
}

type LogoutCommand struct {
	UserID    int
	SessionID string
}

type SendVerificationCodeCommand struct {
	Email string
}
