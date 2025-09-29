package models

import "time"

//TODO: Разбить на сущности

type User struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Username    string    `gorm:"unique;not null" json:"username"`
	Email       string    `gorm:"unique;not null" json:"email"`
	Password    string    `gorm:"not null" json:"-"`
	FirstName   string    `gorm:"not null" json:"first_name"`
	LastName    string    `gorm:"not null" json:"last_name"`
	FatherName  string    `gorm:"not null" json:"father_name"`
	Sex         string    `gorm:"not null" json:"sex"`
	DateOfBirth time.Time `gorm:"not null" json:"date_of_birth"`
	Groups      string    `gorm:"default:Patient; not null" json:"groups"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type Subscription struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	DoctorID     uint      `gorm:"not null" json:"doctor_id"`
	PatientID    uint      `gorm:"not null" json:"patient_id"`
	Status       string    `gorm:"not null;default:pending" json:"status"` // pending / accepted / rejected
	Prescription string    `json:"prescription,omitempty"`
	Diagnosis    string    `json:"diagnosis,omitempty"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	Doctor  User `gorm:"foreignKey:DoctorID" json:"doctor"`
	Patient User `gorm:"foreignKey:PatientID" json:"patient"`
}

type Note struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	CreatedAt        time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	Intensity        int       `gorm:"not null" json:"intensity"`
	PainType         string    `gorm:"not null" json:"pain_type"`
	TookPrescription bool      `gorm:"not null" json:"took_prescription"`
	Description      string    `json:"description,omitempty"`
	BodyPart         uint      `gorm:"not null" json:"body_part"`
	PatientID        uint      `gorm:"not null" json:"patient_id"`
}

type Notification struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null" json:"user_id"`
	Message   string    `json:"message"`
	IsRead    bool      `gorm:"not null;default:false" json:"is_read"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
