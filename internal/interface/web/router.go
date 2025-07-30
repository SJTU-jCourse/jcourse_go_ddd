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
		courses.GET("/search", courseController.SearchCourses)
		courses.GET("/:id", courseController.GetCourseDetail)
	}

	// Review routes
	reviews := v1.Group("/review")
	{
		reviews.GET("/latest", reviewController.GetLatestReviews)
		reviews.POST("/", reviewController.WriteReview)
		reviews.PUT("/:id", reviewController.UpdateReview)
		reviews.DELETE("/:id", reviewController.DeleteReview)
	}

	// Course-specific review routes
	courseReviews := v1.Group("/course/:courseID/reviews")
	{
		courseReviews.GET("/", reviewController.GetCourseReviews)
	}

	// User points routes
	points := v1.Group("/user/:userId/point")
	{
		points.GET("/", pointController.GetUserPoint)
	}
}
