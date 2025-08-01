package dto

import "time"

// Validation DTOs

type EmailValidation struct {
	Email string `json:"email" binding:"required,email" example:"user@example.com"`
}

type PasswordValidation struct {
	Password string `json:"password" binding:"required,min=6" example:"password123"`
}

type CodeValidation struct {
	Code string `json:"code" binding:"required,len=6" example:"123456"`
}

type NicknameValidation struct {
	Nickname string `json:"nickname" binding:"required,min=1,max=50" example:"张三"`
}

type IDValidation struct {
	ID int `json:"id" binding:"required,min=1" example:"1"`
}

type RatingValidation struct {
	Rating int `json:"rating" binding:"required,min=1,max=5" example:"5"`
}

type CommentValidation struct {
	Comment string `json:"comment" binding:"required,min=1,max=1000" example:"这门课程很棒！"`
}

type SemesterValidation struct {
	Semester string `json:"semester" binding:"required" example:"2024-秋"`
}

type GradeValidation struct {
	Grade string `json:"grade" binding:"required" example:"A"`
}

type AmountValidation struct {
	Amount int `json:"amount" binding:"required,min=1" example:"100"`
}

type ReasonValidation struct {
	Reason string `json:"reason" binding:"required,min=1,max=200" example:"活动奖励"`
}

// Search DTOs

type SearchRequest struct {
	Query     string `json:"query" binding:"required,min=1" example:"数据结构"`
	Page      int    `json:"page" binding:"min=1" example:"1"`
	PageSize  int    `json:"page_size" binding:"min=1,max=100" example:"20"`
	SortBy    string `json:"sort_by" binding:"omitempty,oneof=name rating created_at" example:"rating"`
	SortOrder string `json:"sort_order" binding:"omitempty,oneof=asc desc" example:"desc"`
}

type SearchResult struct {
	Items      []interface{} `json:"items"`
	Total      int           `json:"total"`
	Page       int           `json:"page"`
	PageSize   int           `json:"page_size"`
	TotalPages int           `json:"total_pages"`
}

// Filter DTOs

type FilterRequest struct {
	Filters  map[string]interface{} `json:"filters"`
	Page     int                    `json:"page" binding:"min=1" example:"1"`
	PageSize int                    `json:"page_size" binding:"min=1,max=100" example:"20"`
}

type FilterOption struct {
	Value string `json:"value"`
	Label string `json:"label"`
	Count int    `json:"count"`
}

type FilterResponse struct {
	Options []FilterOption `json:"options"`
}

// Upload DTOs

type FileUploadRequest struct {
	File     interface{} `json:"file" binding:"required"`
	FileType string      `json:"file_type" binding:"required,oneof=image avatar document"`
}

type FileUploadResponse struct {
	FileID   string `json:"file_id"`
	FileName string `json:"file_name"`
	FileSize int64  `json:"file_size"`
	FileType string `json:"file_type"`
	URL      string `json:"url"`
}

// Notification DTOs

type NotificationRequest struct {
	Type    string `json:"type" binding:"required,oneof=system user review point"`
	Title   string `json:"title" binding:"required,min=1,max=100"`
	Message string `json:"message" binding:"required,min=1,max=500"`
	UserID  *int   `json:"user_id,omitempty"`
}

type NotificationResponse struct {
	ID        int       `json:"id"`
	Type      string    `json:"type"`
	Title     string    `json:"title"`
	Message   string    `json:"message"`
	UserID    *int      `json:"user_id,omitempty"`
	IsRead    bool      `json:"is_read"`
	CreatedAt time.Time `json:"created_at"`
}

// Health Check DTOs

type HealthCheckResponse struct {
	Status    string            `json:"status"`
	Timestamp time.Time         `json:"timestamp"`
	Services  map[string]string `json:"services"`
	Version   string            `json:"version"`
}

// Configuration DTOs

type ConfigResponse struct {
	Features     map[string]bool `json:"features"`
	Limits       map[string]int  `json:"limits"`
	RateLimits   map[string]int  `json:"rate_limits"`
	AllowedHosts []string        `json:"allowed_hosts"`
}

// Error DTOs

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Code    string `json:"code"`
}

type APIError struct {
	Code      int         `json:"code"`
	Message   string      `json:"message"`
	Details   interface{} `json:"details,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
	RequestID string      `json:"request_id"`
}

// Meta DTOs

type MetaResponse struct {
	Version     string    `json:"version"`
	Environment string    `json:"environment"`
	Timestamp   time.Time `json:"timestamp"`
	Uptime      int64     `json:"uptime"`
}

// WebSocket DTOs

type WebSocketMessage struct {
	Type      string      `json:"type"`
	Data      interface{} `json:"data"`
	Timestamp time.Time   `json:"timestamp"`
	UserID    *int        `json:"user_id,omitempty"`
}

type WebSocketRequest struct {
	Action string                 `json:"action"`
	Data   map[string]interface{} `json:"data"`
}

// Batch DTOs

type BatchRequest struct {
	Items []interface{} `json:"items" binding:"required,min=1,max=100"`
}

type BatchResponse struct {
	Success []interface{} `json:"success"`
	Failed  []BatchError  `json:"failed"`
	Total   int           `json:"total"`
}

type BatchError struct {
	Index int         `json:"index"`
	Error APIError    `json:"error"`
	Item  interface{} `json:"item"`
}

// Export DTOs

type ExportRequest struct {
	Format  string                 `json:"format" binding:"required,oneof=csv json xml"`
	Filters map[string]interface{} `json:"filters"`
	Fields  []string               `json:"fields"`
}

type ExportResponse struct {
	FileID      string    `json:"file_id"`
	FileName    string    `json:"file_name"`
	FileSize    int64     `json:"file_size"`
	DownloadURL string    `json:"download_url"`
	ExpiresAt   time.Time `json:"expires_at"`
}
