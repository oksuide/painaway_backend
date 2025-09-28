package diary

import (
	"painaway_test/models"

	"gorm.io/gorm"
)

type Repo struct {
	DB *gorm.DB
}

type Repository interface {
	GetSubscriptionsByPatientID(patientID uint) ([]models.Subscription, error)
	FindDoctorByUsername(username string) (*models.User, error)
	CreateSubscription(sub *models.Subscription) error
	GetLinkByDoctorAndPatient(doctorID, patientID uint) (*models.Subscription, error)
	UpdateLink(link *models.Subscription) error
	GetAllBodyStatByPatientID(patientID uint) ([]models.Note, error)
	CreateNote(note *models.Note) error
	GetGroupByUserID(userID uint) (string, error)
	GetSubscriptionsByDoctorID(doctorID uint) ([]models.Subscription, error)
	GetLinkByID(linkID uint) (*models.Subscription, error)
}

func NewRepository(db *gorm.DB) Repository {
	return &Repo{DB: db}
}

func (r *Repo) GetSubscriptionsByPatientID(patientID uint) ([]models.Subscription, error) {
	var subs []models.Subscription
	if err := r.DB.Preload("Doctor").Where("patient_id = ? AND status = ?", patientID, "accepted").Find(&subs).Error; err != nil {
		return nil, err
	}
	return subs, nil
}

func (r *Repo) GetSubscriptionsByDoctorID(doctorID uint) ([]models.Subscription, error) {
	var subs []models.Subscription
	if err := r.DB.Preload("Patient").Where("doctor_id = ?", doctorID).Find(&subs).Error; err != nil {
		return nil, err
	}
	return subs, nil
}

func (r *Repo) GetGroupByUserID(userID uint) (string, error) {
	var user models.User
	if err := r.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		return "", err
	}
	return user.Groups, nil
}

func (r *Repo) FindDoctorByUsername(username string) (*models.User, error) {
	var doctor models.User
	if err := r.DB.Where("username = ? AND groups = ?", username, "Doctor").First(&doctor).Error; err != nil {
		return nil, err
	}
	return &doctor, nil
}

func (r *Repo) CreateSubscription(sub *models.Subscription) error {
	return r.DB.Create(sub).Error
}

func (r *Repo) GetLinkByDoctorAndPatient(doctorID, patientID uint) (*models.Subscription, error) {
	var link models.Subscription
	if err := r.DB.Where("doctor_id = ? AND patient_id = ?", doctorID, patientID).First(&link).Error; err != nil {
		return nil, err
	}
	return &link, nil
}

func (r *Repo) GetLinkByID(linkID uint) (*models.Subscription, error) {
	var link models.Subscription
	if err := r.DB.Where("id = ?", linkID).First(&link).Error; err != nil {
		return nil, err
	}
	return &link, nil
}

func (r *Repo) UpdateLink(link *models.Subscription) error {
	return r.DB.Save(link).Error
}

func (r *Repo) GetAllBodyStatByPatientID(patientID uint) ([]models.Note, error) {
	var stats []models.Note
	if err := r.DB.Where("patient_id = ?", patientID).Find(&stats).Error; err != nil {
		return nil, err
	}
	return stats, nil
}

func (r *Repo) CreateNote(note *models.Note) error {
	return r.DB.Create(note).Error
}
