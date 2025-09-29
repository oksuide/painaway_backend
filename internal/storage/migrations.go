package storage

import (
	"painaway_test/models"

	"gorm.io/gorm"
)

// TODO: Полноценные миграции вместо авто

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.User{},
		&models.Note{},
		&models.Subscription{},
		&models.Notification{},
	)
}
