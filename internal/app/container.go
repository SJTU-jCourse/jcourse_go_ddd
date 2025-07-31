package app

import (
	announcementquery "jcourse_go/internal/application/announcement/query"
	"jcourse_go/internal/application/auth"
	authquery "jcourse_go/internal/application/auth/query"
	pointcommand "jcourse_go/internal/application/point/command"
	pointquery "jcourse_go/internal/application/point/query"
	reviewcommand "jcourse_go/internal/application/review/command"
	reviewquery "jcourse_go/internal/application/review/query"
	statisticsquery "jcourse_go/internal/application/statistics/query"
	"jcourse_go/internal/config"
)

type ServiceContainer struct {
	AuthService              auth.AuthService
	CodeService              auth.VerificationCodeService
	CourseQueryService       reviewquery.CourseQueryService
	ReviewCommandService     reviewcommand.ReviewCommandService
	ReviewQueryService       reviewquery.ReviewQueryService
	PointCommandService      pointcommand.PointCommandService
	PointQueryService        pointquery.UserPointQueryService
	UserQueryService         authquery.UserQueryService
	AnnouncementQueryService announcementquery.AnnouncementQueryService
	StatisticsQueryService   statisticsquery.StatisticsQueryService
}

func NewServiceContainer(conf config.Config) *ServiceContainer {
	// Note: In a real implementation, you would inject the actual repositories
	// and dependencies here. This is a placeholder structure.
	return &ServiceContainer{
		AuthService:              nil,
		CodeService:              nil,
		CourseQueryService:       nil,
		ReviewCommandService:     nil,
		ReviewQueryService:       nil,
		PointCommandService:      nil,
		PointQueryService:        nil,
		UserQueryService:         nil,
		AnnouncementQueryService: nil,
		StatisticsQueryService:   nil,
	}
}
