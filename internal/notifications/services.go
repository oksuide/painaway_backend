package notifications

import (
	"painaway_test/models"
)

type Service struct {
	Repo Repository
	Hub  *Hub
}

func NewService(repo Repository, hub *Hub) *Service {
	return &Service{Repo: repo, Hub: hub}
}

func (s *Service) CreateNotification(userID uint, message string) error {
	notification := models.Notification{
		UserID:  userID,
		Message: message,
		IsRead:  false,
	}

	if err := s.Repo.CreateNotification(&notification); err != nil {
		return err
	}

	// пушим сразу в сокет
	_ = s.Hub.Send(userID, &notification)
	return nil
}

func (s *Service) GetNotifications(userID uint) ([]models.Notification, error) {
	return s.Repo.GetNotifications(userID)
}

func (s *Service) MarkNotificationRead(notificationID, userID uint) error {
	return s.Repo.MarkNotificationRead(notificationID, userID)
}

func (s *Service) DeleteNotification(notificationID, userID uint) error {
	return s.Repo.DeleteNotification(notificationID, userID)
}
