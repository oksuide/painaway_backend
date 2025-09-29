package diary

import (
	"painaway_test/models"

	"gorm.io/gorm"
)

type Repo struct {
	DB *gorm.DB
}

type Repository interface {
	CreateNote(note *models.Note) error
	CreateSubscription(sub *models.Subscription) error
	GetSubscriptionsByPatientID(patientID uint, offset, limit int) ([]models.Subscription, error)
	GetSubscriptionsByDoctorID(doctorID uint, offset, limit int) ([]models.Subscription, error)
	GetLinkByDoctorAndPatient(doctorID, patientID uint) (*models.Subscription, error)
	GetAllStatByPatientID(patientID uint) ([]models.Note, error)
	GetDoctorByUsername(username string) (*models.User, error)
	GetGroupByUserID(userID uint) (string, error)
	GetLinkByID(linkID uint) (*models.Subscription, error)
	UpdateLink(link *models.Subscription) error
}

func NewRepository(db *gorm.DB) Repository {
	return &Repo{DB: db}
}

func (r *Repo) GetSubscriptionsByPatientID(patientID uint, offset, limit int) ([]models.Subscription, error) {
	var subs []models.Subscription
	if err := r.DB.Preload("Doctor").
		Where("patient_id = ? AND status = ?", patientID, "accepted").
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&subs).Error; err != nil {
		return nil, err
	}
	return subs, nil
}
func (r *Repo) GetSubscriptionsByDoctorID(doctorID uint, offset, limit int) ([]models.Subscription, error) {
	var subs []models.Subscription
	if err := r.DB.Preload("Patient").
		Where("doctor_id = ?", doctorID).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&subs).Error; err != nil {
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

func (r *Repo) GetDoctorByUsername(username string) (*models.User, error) {
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

func (r *Repo) GetAllStatByPatientID(patientID uint) ([]models.Note, error) {
	var stats []models.Note
	if err := r.DB.Where("patient_id = ?", patientID).Find(&stats).Error; err != nil {
		return nil, err
	}
	return stats, nil
}

func (r *Repo) CreateNote(note *models.Note) error {
	return r.DB.Create(note).Error
}
