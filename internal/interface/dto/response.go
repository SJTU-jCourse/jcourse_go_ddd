package dto

import "time"

// Course Response DTOs

type CourseDetailResponse struct {
	Course     CourseDetail        `json:"course"`
	Teachers   []TeacherInfo       `json:"teachers"`
	Offerings  []OfferedCourseInfo `json:"offerings"`
	Rating     RatingInfo          `json:"rating"`
	Reviews    []ReviewSummary     `json:"reviews"`
	IsEnrolled bool                `json:"is_enrolled"`
	IsWatched  bool                `json:"is_watched"`
	UserReview *UserReviewInfo     `json:"user_review,omitempty"`
}

type CourseDetail struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Code        string `json:"code"`
	Department  string `json:"department"`
	Credits     int    `json:"credits"`
	Description string `json:"description"`
	Category    string `json:"category"`
}

type TeacherInfo struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Title      string `json:"title"`
	Department string `json:"department"`
}

type OfferedCourseInfo struct {
	ID        int     `json:"id"`
	CourseID  int     `json:"course_id"`
	TeacherID int     `json:"teacher_id"`
	Semester  string  `json:"semester"`
	Capacity  int     `json:"capacity"`
	Enrolled  int     `json:"enrolled"`
	Rating    float64 `json:"rating"`
}

type RatingInfo struct {
	TotalRating  float64 `json:"total_rating"`
	RatingCount  int     `json:"rating_count"`
	Rating1Count int     `json:"rating_1_count"`
	Rating2Count int     `json:"rating_2_count"`
	Rating3Count int     `json:"rating_3_count"`
	Rating4Count int     `json:"rating_4_count"`
	Rating5Count int     `json:"rating_5_count"`
}

type ReviewSummary struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Nickname  string    `json:"nickname"`
	Rating    int       `json:"rating"`
	Comment   string    `json:"comment"`
	Semester  string    `json:"semester"`
	Grade     string    `json:"grade"`
	CreatedAt time.Time `json:"created_at"`
	LikeCount int       `json:"like_count"`
}

type UserReviewInfo struct {
	ID        int       `json:"id"`
	Rating    int       `json:"rating"`
	Comment   string    `json:"comment"`
	Semester  string    `json:"semester"`
	Grade     string    `json:"grade"`
	CreatedAt time.Time `json:"created_at"`
}

// Review Response DTOs

type ReviewDetailResponse struct {
	Review     ReviewDetail      `json:"review"`
	Course     CourseInReview    `json:"course"`
	User       UserInReview      `json:"user"`
	Revisions  []ReviewRevision  `json:"revisions"`
	Actions    []ReviewAction    `json:"actions"`
	UserAction *ReviewUserAction `json:"user_action,omitempty"`
}

type ReviewDetail struct {
	ID        int       `json:"id"`
	CourseID  int       `json:"course_id"`
	UserID    int       `json:"user_id"`
	Comment   string    `json:"comment"`
	Rating    int       `json:"rating"`
	Semester  string    `json:"semester"`
	Grade     string    `json:"grade"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	LikeCount int       `json:"like_count"`
}

type CourseInReview struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
}

type UserInReview struct {
	ID       int    `json:"id"`
	Nickname string `json:"nickname"`
	Role     string `json:"role"`
}

type ReviewRevision struct {
	ID        int       `json:"id"`
	ReviewID  int       `json:"review_id"`
	Comment   string    `json:"comment"`
	Rating    int       `json:"rating"`
	Semester  string    `json:"semester"`
	Grade     string    `json:"grade"`
	CreatedAt time.Time `json:"created_at"`
}

type ReviewAction struct {
	ID         int           `json:"id"`
	ReviewID   int           `json:"review_id"`
	UserID     int           `json:"user_id"`
	ActionType string        `json:"action_type"`
	CreatedAt  time.Time     `json:"created_at"`
	User       *UserInReview `json:"user,omitempty"`
}

type ReviewUserAction struct {
	ActionType string `json:"action_type"`
}

// User Response DTOs

type UserInfoResponse struct {
	User       UserInfo       `json:"user"`
	Statistics UserStatistics `json:"statistics"`
}

type UserInfo struct {
	ID         int       `json:"id"`
	Email      string    `json:"email"`
	Nickname   string    `json:"nickname"`
	Role       string    `json:"role"`
	CreatedAt  time.Time `json:"created_at"`
	LastSeenAt time.Time `json:"last_seen_at"`
}

type UserStatistics struct {
	ReviewCount  int `json:"review_count"`
	LikeCount    int `json:"like_count"`
	PointBalance int `json:"point_balance"`
}

// Point Response DTOs

type PointHistoryResponse struct {
	Balance      int                `json:"balance"`
	Transactions []PointTransaction `json:"transactions"`
	Pagination   PaginationResponse `json:"pagination"`
}

type PointTransaction struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Amount    int       `json:"amount"`
	Type      string    `json:"type"`
	Reason    string    `json:"reason"`
	CreatedAt time.Time `json:"created_at"`
}

// Filter Response DTOs

type CourseFilterResponse struct {
	Departments []string `json:"departments"`
	Categories  []string `json:"categories"`
}

// Announcement Response DTOs

type AnnouncementResponse struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Statistics Response DTOs

type SystemStatisticsResponse struct {
	TotalUsers       int `json:"total_users"`
	TotalCourses     int `json:"total_courses"`
	TotalReviews     int `json:"total_reviews"`
	TotalPoints      int `json:"total_points"`
	ActiveUsersCount int `json:"active_users_count"`
	RecentReviews    int `json:"recent_reviews"`
}
