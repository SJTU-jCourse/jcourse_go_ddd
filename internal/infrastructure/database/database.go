package database

import (
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"jcourse_go/internal/config"
)

const (
	MaxIdleConns     = 10
	MaxOpenConns     = 100
	ConnMaxLifetime  = time.Hour
)

func NewDatabase(cfg config.DBConfig) (*gorm.DB, error) {
	dsn := cfg.DSN

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	sqlDB.SetMaxIdleConns(MaxIdleConns)
	sqlDB.SetMaxOpenConns(MaxOpenConns)
	sqlDB.SetConnMaxLifetime(ConnMaxLifetime)

	return db, nil
}
