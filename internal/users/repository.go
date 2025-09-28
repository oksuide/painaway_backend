package users

import (
	"painaway_test/models"

	"gorm.io/gorm"
)

type Repo struct {
	DB *gorm.DB
}

type Repository interface {
	CreateUser(user *models.User) error
	EmailExists(email string) (bool, error)
	UserExists(username string) (bool, error)
	GetUserByUsername(username string) (*models.User, error)
	GetUserByID(ID uint) (*models.User, error)
}

func NewRepository(db *gorm.DB) Repository {
	return &Repo{DB: db}
}

func (r *Repo) CreateUser(user *models.User) error {
	return r.DB.Create(user).Error
}

func (r *Repo) EmailExists(email string) (bool, error) {
	var count int64
	err := r.DB.Model(&models.User{}).
		Where("email = ?", email).
		Count(&count).
		Error

	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *Repo) UserExists(username string) (bool, error) {
	var count int64
	if err := r.DB.
		Model(&models.User{}).
		Where("username = ?", username).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
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
