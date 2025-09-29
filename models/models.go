package models

import "time"

//TODO: Разбить на сущности

type User struct {
	ID          uint      `gorm:"primaryKey"`
	Username    string    `gorm:"unique;not null"`
	Email       string    `gorm:"unique;not null"`
	Password    string    `gorm:"not null"`
	FirstName   string    `gorm:"not null"`
	LastName    string    `gorm:"not null"`
	FatherName  string    `gorm:"not null"`
	Sex         string    `gorm:"not null"`
	DateOfBirth time.Time `gorm:"not null"`
	Groups      string    `gorm:"default:Patient; not null"`
	Created_at  time.Time `gorm:"not null"`
	Updated_at  time.Time `gorm:"not null"`
}

type Subscription struct {
	ID           uint   `gorm:"primaryKey"`
	DoctorID     uint   `gorm:"not null"`
	PatientID    uint   `gorm:"not null"`
	Status       string `gorm:"not null; default:pending"` // pending / accepted / rejected
	Prescription string
	Diagnosis    string
	CreatedAt    time.Time `gorm:"not null"`
	UpdatedAt    time.Time `gorm:"not null"`

	Doctor  User `gorm:"foreignKey:DoctorID"`
	Patient User `gorm:"foreignKey:PatientID"`
}

type Note struct {
	ID               uint      `json:"id" gorm:"primaryKey"`
	CreatedAt        time.Time `json:"created_at" gorm:"autoCreateTime"`
	Intensity        int       `json:"intensity" gorm:"not null"`
	PainType         string    `json:"pain_type" gorm:"not null"`
	TookPrescription bool      `json:"tookPrescription" gorm:"not null"`
	Description      string    `json:"description"`
	BodyPart         uint      `json:"body_part" gorm:"not null"`
	PatientID        uint      `json:"patient_id" gorm:"not null"`
}

type Notification struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id" gorm:"not null"`
	Message   string    `json:"message"`
	IsRead    bool      `json:"is_read" gorm:"not null; default:False"`
	CreatedAt time.Time `json:"created_at" gorm:"not null"`
}
