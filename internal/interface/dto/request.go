package dto

import "time"

// Auth DTOs

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email" example:"user@example.com"`
	Password string `json:"password" binding:"required,min=6" example:"password123"`
}

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email" example:"user@example.com"`
	Code     string `json:"code" binding:"required" example:"123456"`
	Password string `json:"password" binding:"required,min=6" example:"password123"`
}

type SendVerificationCodeRequest struct {
	Email string `json:"email" binding:"required,email" example:"user@example.com"`
}

type AuthResponse struct {
	SessionID string `json:"session_id" example:"abc123"`
}

// Course DTOs

type AddEnrolledCourseRequest struct {
	CourseID int `json:"course_id" binding:"required" example:"1"`
}

type WatchCourseRequest struct {
	Watch bool `json:"watch" binding:"required" example:"true"`
}

type CourseSearchQuery struct {
	Name       string `form:"name" binding:"omitempty" example:"数据结构"`
	Code       string `form:"code" binding:"omitempty" example:"CS101"`
	Department string `form:"department" binding:"omitempty" example:"计算机系"`
}

// Review DTOs

type WriteReviewRequest struct {
	CourseID int    `json:"course_id" binding:"required" example:"1"`
	Comment  string `json:"comment" binding:"required,min=1,max=1000" example:"这门课程很棒！"`
	Rating   int    `json:"rating" binding:"required,min=1,max=5" example:"5"`
	Semester string `json:"semester" binding:"required" example:"2024-秋"`
	Grade    string `json:"grade" binding:"required" example:"A"`
}

type UpdateReviewRequest struct {
	Comment  string `json:"comment" binding:"required,min=1,max=1000" example:"这门课程很棒！"`
	Rating   int    `json:"rating" binding:"required,min=1,max=5" example:"5"`
	Semester string `json:"semester" binding:"required" example:"2024-秋"`
	Grade    string `json:"grade" binding:"required" example:"A"`
}

type PostReviewActionRequest struct {
	ActionType string `json:"action_type" binding:"required" example:"like"`
}

// User DTOs

type UpdateUserInfoRequest struct {
	Nickname string `json:"nickname" binding:"required,min=1,max=50" example:"张三"`
}

// Point DTOs (Admin)

type CreatePointRequest struct {
	UserID int    `json:"user_id" binding:"required" example:"1"`
	Amount int    `json:"amount" binding:"required" example:"100"`
	Reason string `json:"reason" binding:"required" example:"活动奖励"`
}

type PointTransactionRequest struct {
	FromUserID int    `json:"from_user_id" binding:"required" example:"1"`
	ToUserID   int    `json:"to_user_id" binding:"required" example:"2"`
	Amount     int    `json:"amount" binding:"required" example:"50"`
	Reason     string `json:"reason" binding:"required" example:"转账"`
}

// Response DTOs

type CourseListResponse struct {
	Courses []CourseListItem `json:"courses"`
	Total   int              `json:"total"`
}

type ReviewListResponse struct {
	Reviews []ReviewItem `json:"reviews"`
	Total   int          `json:"total"`
}

type UserPointResponse struct {
	Balance      int                `json:"balance"`
	Transactions []PointTransaction `json:"transactions"`
}

type CourseListItem struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Code        string  `json:"code"`
	Department  string  `json:"department"`
	Credits     int     `json:"credits"`
	Rating      float64 `json:"rating"`
	ReviewCount int     `json:"review_count"`
}

type ReviewItem struct {
	ID        int       `json:"id"`
	CourseID  int       `json:"course_id"`
	UserID    int       `json:"user_id"`
	Comment   string    `json:"comment"`
	Rating    int       `json:"rating"`
	Semester  string    `json:"semester"`
	Grade     string    `json:"grade"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PointTransaction struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Amount    int       `json:"amount"`
	Type      string    `json:"type"`
	Reason    string    `json:"reason"`
	CreatedAt time.Time `json:"created_at"`
}

// Pagination DTOs

type PaginationQuery struct {
	Page     int `form:"page" binding:"min=1" example:"1"`
	PageSize int `form:"page_size" binding:"min=1,max=100" example:"20"`
}

type PaginationResponse struct {
	Page       int `json:"page"`
	PageSize   int `json:"page_size"`
	TotalItems int `json:"total_items"`
	TotalPages int `json:"total_pages"`
}

// Error Response DTOs

type ValidationErrorResponse struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Code    int                       `json:"code"`
	Message string                    `json:"message"`
	Details []ValidationErrorResponse `json:"details,omitempty"`
}
