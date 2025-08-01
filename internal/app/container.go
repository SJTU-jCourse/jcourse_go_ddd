package app

import (
	announcementquery "jcourse_go/internal/application/announcement/query"
	app_auth "jcourse_go/internal/application/auth"
	authquery "jcourse_go/internal/application/auth/query"
	pointcommand "jcourse_go/internal/application/point/command"
	pointquery "jcourse_go/internal/application/point/query"
	reviewcommand "jcourse_go/internal/application/review/command"
	reviewquery "jcourse_go/internal/application/review/query"
	statisticsquery "jcourse_go/internal/application/statistics/query"
	"jcourse_go/internal/config"
	"jcourse_go/internal/domain/email"
	"jcourse_go/internal/infrastructure/database"
	redisclient "jcourse_go/internal/infrastructure/redis"
	"jcourse_go/internal/infrastructure/repository"
	"jcourse_go/pkg/password"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type ServiceContainer struct {
	DB    *gorm.DB
	Redis *redis.Client

	AuthService              app_auth.AuthService
	CodeService              app_auth.VerificationCodeService
	CourseQueryService       reviewquery.CourseQueryService
	ReviewCommandService     reviewcommand.ReviewCommandService
	ReviewQueryService       reviewquery.ReviewQueryService
	PointCommandService      pointcommand.PointCommandService
	PointQueryService        pointquery.UserPointQueryService
	UserQueryService         authquery.UserQueryService
	AnnouncementQueryService announcementquery.AnnouncementQueryService
	StatisticsQueryService   statisticsquery.StatisticsQueryService
}

func NewServiceContainer(conf config.Config) (*ServiceContainer, error) {
	db, err := database.NewDatabase(conf.DB)
	if err != nil {
		return nil, err
	}

	redisClient, err := redisclient.NewRedisClient(conf.Redis)
	if err != nil {
		return nil, err
	}

	// Auto-migrate is disabled - use manual migration command instead
	// if err := migrations.Migrate(db); err != nil {
	// 	return nil, err
	// }

	userRepo := repository.NewUserRepository(db)
	codeRepo := repository.NewCodeRepository(db)
	sessionRepo := repository.NewSessionRepository(redisClient)
	reviewRepo := repository.NewReviewRepository(db)
	courseRepo := repository.NewCourseRepository(db)
	pointRepo := repository.NewUserPointRepository(db)

	hasher := password.NewHasher()
	emailService := email.NewEmailService()

	container := &ServiceContainer{
		DB:    db,
		Redis: redisClient,

		AuthService:              app_auth.NewAuthService(userRepo, hasher, sessionRepo, nil),
		CodeService:              app_auth.NewVerificationCodeService(emailService, codeRepo),
		CourseQueryService:       reviewquery.NewCourseQueryService(courseRepo, reviewRepo),
		ReviewCommandService:     reviewcommand.NewReviewCommandService(reviewRepo, courseRepo),
		ReviewQueryService:       reviewquery.NewReviewQueryService(reviewRepo, courseRepo),
		PointCommandService:      pointcommand.NewPointCommandService(pointRepo),
		PointQueryService:        pointquery.NewUserPointQueryService(pointRepo),
		UserQueryService:         authquery.NewUserQueryService(userRepo),
		AnnouncementQueryService: announcementquery.NewAnnouncementQueryService(nil),
		StatisticsQueryService:   statisticsquery.NewStatisticsQueryService(nil),
	}

	return container, nil
}

func (c *ServiceContainer) Close() error {
	if c.DB != nil {
		sqlDB, _ := c.DB.DB()
		if sqlDB != nil {
			sqlDB.Close()
		}
	}
	if c.Redis != nil {
		c.Redis.Close()
	}
	return nil
}
