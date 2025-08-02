package command

type LoginCommand struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RegisterCommand struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
	Code     string `json:"code" binding:"required"`
}

type LogoutCommand struct {
	SessionID string `json:"session_id" binding:"required"`
}

type SendVerificationCodeCommand struct {
	Email string `json:"email" binding:"required,email"`
}
