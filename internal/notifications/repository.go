package notifications

import (
	"painaway_test/models"

	"gorm.io/gorm"
)

type Repo struct {
	DB *gorm.DB
}

type Repository interface {
	CreateNotification(notification *models.Notification) error
	GetNotifications(userID uint) ([]models.Notification, error)
	DeleteNotification(id uint, userID uint) error
	MarkNotificationRead(id uint, userID uint) error
}

func NewRepository(db *gorm.DB) Repository {
	return &Repo{DB: db}
}

func (r *Repo) CreateNotification(notification *models.Notification) error {
	return r.DB.Create(notification).Error
}

func (r *Repo) GetNotifications(userID uint) ([]models.Notification, error) {
	var notifications []models.Notification
	err := r.DB.Where("user_id = ?", userID).Order("created_at DESC").Find(&notifications).Error
	return notifications, err
}

func (r *Repo) MarkNotificationRead(id uint, userID uint) error {
	return r.DB.Model(&models.Notification{}).
		Where("id = ? AND user_id = ?", id, userID).
		Update("is_read", true).Error
}

func (r *Repo) DeleteNotification(id uint, userID uint) error {
	return r.DB.Where("id = ? AND user_id = ?", id, userID).
		Delete(&models.Notification{}).Error
}
