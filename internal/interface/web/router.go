package web

import (
	"github.com/gin-gonic/gin"

	"jcourse_go/internal/app"
)

func RegisterRouter(g *gin.Engine, s *app.ServiceContainer) {
	authController := NewAuthController(s.AuthService, s.CodeService)
	courseController := NewCourseController(s.CourseQueryService)
	reviewController := NewReviewController(s.ReviewCommandService, s.ReviewQueryService)
	pointController := NewUserPointController(s.PointCommandService, s.PointQueryService)
	userController := NewUserController(s.UserQueryService, s.ReviewQueryService)
	announcementController := NewAnnouncementController(s.AnnouncementQueryService)
	statisticsController := NewStatisticsController(s.StatisticsQueryService)

	// Apply authentication middleware to all routes
	g.Use(AuthMiddleware(s.AuthService))

	// API version 1 group
	v1 := g.Group("/api/v1")

	// Auth routes (no auth required)
	auth := v1.Group("/auth")
	{
		auth.POST("/login", authController.Login)
		auth.POST("/register", authController.Register)
		auth.POST("/logout", authController.Logout)
		auth.POST("/send-code", authController.SendVerificationCode)
	}

	// Course routes
	courses := v1.Group("/course")
	{
		courses.GET("/filter", courseController.GetCourseFilter)
		courses.GET("/enroll", RequireAuth(), courseController.GetUserEnrolledCourses)
		courses.POST("/enroll", RequireAuth(), courseController.AddUserEnrolledCourse)
		courses.GET("/search", courseController.SearchCourses)
		courses.GET("/:id", courseController.GetCourseDetail)
		courses.GET("/:id/review", reviewController.GetCourseReviews)
		courses.POST("/:id/watch", RequireAuth(), courseController.WatchCourse)
	}

	// Review routes
	reviews := v1.Group("/review")
	{
		reviews.GET("", reviewController.GetLatestReviews)
		reviews.POST("", RequireAuth(), reviewController.WriteReview)
		reviews.PUT("/:id", RequireAuth(), reviewController.UpdateReview)
		reviews.DELETE("/:id", RequireAuth(), reviewController.DeleteReview)
		reviews.POST("/:id/action", RequireAuth(), reviewController.PostReviewAction)
		reviews.DELETE("/:id/action/:actionID", RequireAuth(), reviewController.DeleteReviewAction)
		reviews.GET("/:id/revision", reviewController.GetReviewRevisions)
	}

	// User routes
	users := v1.Group("/user")
	{
		users.GET("/info", RequireAuth(), userController.GetUserInfo)
		users.POST("/info", RequireAuth(), userController.UpdateUserInfo)
		users.GET("/point", RequireAuth(), pointController.GetUserPoint)
		users.GET("/review", RequireAuth(), userController.GetUserReviews)
	}

	// Admin routes
	admin := v1.Group("/admin")
	admin.Use(RequireAdmin())
	{
		admin.POST("/point", pointController.CreatePoint)
		admin.POST("/point/transaction", pointController.Transaction)
	}

	announcements := v1.Group("/announcement")
	{
		announcements.GET("", announcementController.GetAnnouncements)
	}

	statistics := v1.Group("/statistics")
	{
		statistics.GET("", statisticsController.GetSystemStatistics)
	}
}
