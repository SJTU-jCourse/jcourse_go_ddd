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

	// API version 1 group
	v1 := g.Group("/api/v1")

	// Auth routes
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
		courses.GET("/enroll", courseController.GetUserEnrolledCourses)
		courses.POST("/enroll", courseController.AddUserEnrolledCourse)
		courses.GET("/search", courseController.SearchCourses)
		courses.GET("/:id", courseController.GetCourseDetail)
		courses.GET("/:id/review", reviewController.GetCourseReviews)
		courses.POST("/:id/watch", courseController.WatchCourse)
	}

	// Review routes
	reviews := v1.Group("/review")
	{
		reviews.GET("", reviewController.GetLatestReviews)
		reviews.POST("", reviewController.WriteReview)
		reviews.PUT("/:id", reviewController.UpdateReview)
		reviews.DELETE("/:id", reviewController.DeleteReview)
		reviews.POST("/:id/action", reviewController.PostReviewAction)
		reviews.DELETE("/:id/action/:actionID", reviewController.DeleteReviewAction)
		reviews.GET("/:id/revision", reviewController.GetReviewRevisions)
	}

	// User routes
	users := v1.Group("/user")
	{
		users.GET("/info", userController.GetUserInfo)
		users.POST("/info", userController.UpdateUserInfo)
		users.GET("/point", pointController.GetUserPoint)
		users.GET("/review", userController.GetUserReviews)
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
