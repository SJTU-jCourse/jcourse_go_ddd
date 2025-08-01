package app

import (
	"strconv"

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
	eventbusimpl "jcourse_go/internal/infrastructure/eventbus"
	redisclient "jcourse_go/internal/infrastructure/redis"
	"jcourse_go/internal/infrastructure/repository"
	eventhandler "jcourse_go/internal/interface/handler"
	"jcourse_go/pkg/password"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type ServiceContainer struct {
	DB       *gorm.DB
	Redis    *redis.Client
	EventBus event.EventBusPublisher

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

	var eventBus event.EventBusPublisher
	var eventPublisher event.Publisher = nil

	if conf.Event.Enabled {
		redisAddr := conf.Redis.Addr
		if conf.Redis.Port != 0 {
			redisAddr = redisAddr + ":" + strconv.Itoa(conf.Redis.Port)
		}

		var err error
		eventBus, err = eventbusimpl.NewAsynqEventBus(redisAddr)
		if err != nil {
			return nil, err
		}

		// Register event handlers
		reviewHandler := eventhandler.NewReviewEventHandler()
		pointHandler := eventhandler.NewPointEventHandler()
		statsHandler := eventhandler.NewStatisticsEventHandler()

		if err := eventBus.Register(event.TypeReviewCreated, reviewHandler); err != nil {
			return nil, err
		}
		if err := eventBus.Register(event.TypeReviewModified, reviewHandler); err != nil {
			return nil, err
		}
		if err := eventBus.Register(event.TypeReviewCreated, pointHandler); err != nil {
			return nil, err
		}
		if err := eventBus.Register(event.TypeReviewModified, statsHandler); err != nil {
			return nil, err
		}

		eventPublisher = eventBus
	}

	container := &ServiceContainer{
		DB:       db,
		Redis:    redisClient,
		EventBus: eventBus,

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
	if c.EventBus != nil {
		c.EventBus.Shutdown()
	}
	return nil
}
