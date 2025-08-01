package migrations

import (
	"gorm.io/gorm"

	"jcourse_go/internal/infrastructure/entity"
)

type Migration struct {
	ID          uint   `gorm:"primaryKey"`
	Name        string `gorm:"unique;not null"`
	AppliedAt   int64  `gorm:"not null"`
	Description string `gorm:"type:text"`
}

func Migrate(db *gorm.DB) error {
	if err := db.AutoMigrate(&Migration{}); err != nil {
		return err
	}

	migrations := []struct {
		name        string
		description string
		migrate     func(*gorm.DB) error
	}{
		{
			name:        "001_initial_schema",
			description: "Create initial database schema for all entities",
			migrate:     migrateInitialSchema,
		},
	}

	for _, migration := range migrations {
		var count int64
		if err := db.Model(&Migration{}).Where("name = ?", migration.name).Count(&count).Error; err != nil {
			return err
		}

		if count > 0 {
			continue
		}

		if err := migration.migrate(db); err != nil {
			return err
		}

		if err := db.Create(&Migration{
			Name:        migration.name,
			AppliedAt:   0,
			Description: migration.description,
		}).Error; err != nil {
			return err
		}
	}

	return nil
}

func migrateInitialSchema(db *gorm.DB) error {
	if err := db.AutoMigrate(
		&entity.User{},
		&entity.VerificationCode{},
		&entity.Course{},
		&entity.Review{},
		&entity.ReviewRevision{},
		&entity.ReviewAction{},
		&entity.UserPointRecord{},
		&entity.UserEnrolledCourse{},
		&entity.CourseWatch{},
	); err != nil {
		return err
	}

	return nil
}
