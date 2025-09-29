package users

import (
	"errors"
	"painaway_test/models"

	"gorm.io/gorm"
)

type Repo struct {
	DB *gorm.DB
}

type Repository interface {
	CreateUser(user *models.User) error
	GetUserByUsername(username string) (*models.User, error)
	GetUserByID(ID uint) (*models.User, error)
	IsUserExistWithEmail(email string) (bool, error)
	IsUserExistWithUsername(username string) (bool, error)
}

func NewRepository(db *gorm.DB) Repository {
	return &Repo{DB: db}
}

func (r *Repo) CreateUser(user *models.User) error {
	return r.DB.Create(user).Error
}

func (r *Repo) IsUserExistWithEmail(email string) (bool, error) {
	var user models.User
	err := r.DB.Select("id").Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (r *Repo) IsUserExistWithUsername(username string) (bool, error) {
	var user models.User
	err := r.DB.Select("id").Where("username = ?", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (r *Repo) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	if err := r.DB.
		Where("username = ?", username).
		First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repo) GetUserByID(ID uint) (*models.User, error) {
	var user models.User
	if err := r.DB.
		Where("id = ?", ID).
		First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
