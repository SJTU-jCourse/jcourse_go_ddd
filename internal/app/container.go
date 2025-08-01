package app

import (
	announcementquery "jcourse_go/internal/application/announcement/query"
	"jcourse_go/internal/application/auth"
	authcommand "jcourse_go/internal/application/auth/command"
	authquery "jcourse_go/internal/application/auth/query"
	pointcommand "jcourse_go/internal/application/point/command"
	pointquery "jcourse_go/internal/application/point/query"
	reviewcommand "jcourse_go/internal/application/review/command"
	reviewquery "jcourse_go/internal/application/review/query"
	statisticsquery "jcourse_go/internal/application/statistics/query"
	"jcourse_go/internal/application/statistics/service"
	"jcourse_go/internal/config"
	"jcourse_go/internal/domain/email"
	"jcourse_go/internal/domain/event"
	"jcourse_go/internal/domain/permission"
	"jcourse_go/internal/infrastructure/database"
	redisclient "jcourse_go/internal/infrastructure/redis"
	"jcourse_go/internal/infrastructure/repository"
	"jcourse_go/pkg/password"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type ServiceContainer struct {
	DB       *gorm.DB
	Redis    *redis.Client

	AuthCommandService       authcommand.AuthCommandService
	AuthQueryService         authquery.AuthQueryService
	CodeService              auth.VerificationCodeService
	CourseCommandService     reviewcommand.CourseCommandService
	CourseQueryService       reviewquery.CourseQueryService
	ReviewCommandService     reviewcommand.ReviewCommandService
	ReviewQueryService       reviewquery.ReviewQueryService
	PointCommandService      pointcommand.PointCommandService
	PointQueryService        pointquery.UserPointQueryService
	UserCommandService       authcommand.UserCommandService
	UserQueryService         authquery.UserQueryService
	AnnouncementQueryService announcementquery.AnnouncementQueryService
	StatisticsQueryService   statisticsquery.StatisticsQueryService
	DailyStatisticsService   service.DailyStatisticsService
}

func NewServiceContainer(conf config.Config, eventPublisher event.Publisher) (*ServiceContainer, error) {
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
	sessionRepo := repository.NewSessionRepository(redisClient)
	reviewRepo := repository.NewReviewRepository(db)
	courseRepo := repository.NewCourseRepository(db)
	pointRepo := repository.NewUserPointRepository(db)
	announcementRepo := repository.NewAnnouncementRepository(db)
	statisticsRepo := repository.NewStatisticsRepository(db)

	hasher := password.NewHasher()
	permissionService := permission.NewPermissionService(userRepo)

	codeRepo := repository.NewCodeRepository(db)
	emailService := email.NewEmailService()
	codeService := auth.NewVerificationCodeService(emailService, codeRepo)

	container := &ServiceContainer{
		DB:       db,
		Redis:    redisClient,

		AuthCommandService:       authcommand.NewAuthCommandService(userRepo, hasher, sessionRepo, codeService),
		AuthQueryService:         authquery.NewAuthQueryService(userRepo, sessionRepo),
		CodeService:              codeService,
		CourseCommandService:     reviewcommand.NewCourseCommandService(courseRepo),
		CourseQueryService:       reviewquery.NewCourseQueryService(courseRepo, reviewRepo),
		ReviewCommandService:     reviewcommand.NewReviewCommandService(reviewRepo, courseRepo, permissionService, eventPublisher),
		ReviewQueryService:       reviewquery.NewReviewQueryService(reviewRepo, courseRepo),
		PointCommandService:      pointcommand.NewPointCommandService(pointRepo),
		PointQueryService:        pointquery.NewUserPointQueryService(pointRepo),
		UserCommandService:       authcommand.NewUserCommandService(userRepo),
		UserQueryService:         authquery.NewUserQueryService(userRepo),
		AnnouncementQueryService: announcementquery.NewAnnouncementQueryService(announcementRepo),
		StatisticsQueryService:   statisticsquery.NewStatisticsQueryService(statisticsRepo),
		DailyStatisticsService:   service.NewDailyStatisticsService(statisticsRepo),
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
