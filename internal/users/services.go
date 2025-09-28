package users

import (
	"painaway_test/models"
)

type Service struct {
	Repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{Repo: repo}
}

func (s *Service) GetProfile(userID uint) (*models.User, error) {
	return s.Repo.GetUserByID(userID)
}
