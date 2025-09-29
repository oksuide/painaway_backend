package auth

import (
	"errors"
	"fmt"
	"painaway_test/internal/users"
	"painaway_test/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Service struct {
	UserRepo users.Repository
}

func NewService(userRepo users.Repository) *Service {
	return &Service{UserRepo: userRepo}
}

func (s *Service) Register(user *models.User) error {
	exists, err := s.UserRepo.IsUserExistWithEmail(user.Email)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("email already registered")
	}

	exists, err = s.UserRepo.IsUserExistWithUsername(user.Username)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("username already taken")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashed)

	return s.UserRepo.CreateUser(user)
}

func (s *Service) Login(username, password string) (*models.User, error) {
	user, err := s.UserRepo.GetUserByUsername(username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("invalid credentials")
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}
	return user, nil
}
