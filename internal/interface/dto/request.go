package dto

// Course Request DTOs

type AddEnrolledCourseRequest struct {
	CourseID int `json:"course_id" binding:"required" example:"1"`
}

type WatchCourseRequest struct {
	Watch bool `json:"watch" binding:"required" example:"true"`
}

// Review Request DTOs

type PostReviewActionRequest struct {
	ActionType string `json:"action_type" binding:"required" example:"like"`
}

// User Request DTOs

type UpdateUserInfoRequest struct {
	Nickname string `json:"nickname" binding:"required,min=1,max=50" example:"张三"`
}

// Point Request DTOs (Admin)

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
