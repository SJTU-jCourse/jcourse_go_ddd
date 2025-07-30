package app

import (
	"jcourse_go/internal/application/auth"
	pointcommand "jcourse_go/internal/application/point/command"
	pointquery "jcourse_go/internal/application/point/query"
	reviewcommand "jcourse_go/internal/application/review/command"
	reviewquery "jcourse_go/internal/application/review/query"
	"jcourse_go/internal/config"
)

type ServiceContainer struct {
	AuthService          auth.AuthService
	CodeService          auth.VerificationCodeService
	CourseQueryService   reviewquery.CourseQueryService
	ReviewCommandService reviewcommand.ReviewCommandService
	ReviewQueryService   reviewquery.ReviewQueryService
	PointCommandService  pointcommand.PointCommandService
	PointQueryService    pointquery.UserPointQueryService
}

func NewServiceContainer(conf config.Config) *ServiceContainer {
	// Note: In a real implementation, you would inject the actual repositories
	// and dependencies here. This is a placeholder structure.
	return &ServiceContainer{
		AuthService:          nil,
		CodeService:          nil,
		CourseQueryService:   nil,
		ReviewCommandService: nil,
		ReviewQueryService:   nil,
		PointCommandService:  nil,
		PointQueryService:    nil,
	}
}
