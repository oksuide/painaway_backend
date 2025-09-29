package utils

import "time"

type PatientLinkDTO struct {
	ID           uint      `json:"id"`
	Status       string    `json:"status"`
	Doctor       DoctorDTO `json:"doctor"`
	Prescription string    `json:"prescription"`
}

type DoctorLinkDTO struct {
	ID           uint               `json:"id"`
	Status       string             `json:"status"`
	Patient      PatientDTO         `json:"patient"`
	Prescription SetPrescriptionDTO `json:"prescription"`
	Diagnosis    SetDiagnosisDTO    `json:"diagnosis"`
}

type DoctorDTO struct {
	ID         uint   `json:"id"`
	Username   string `json:"username"`
	LastName   string `json:"last_name"`
	FirstName  string `json:"first_name"`
	FatherName string `json:"father_name"`
}

type PatientDTO struct {
	ID          uint   `json:"id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	FatherName  string `json:"father_name"`
	Sex         string `json:"sex"`
	DateOfBirth string `json:"date_of_birth"`
}

type SetPrescriptionDTO struct {
	Link         uint   `json:"link"`
	Prescription string `json:"prescription"`
}

type SetDiagnosisDTO struct {
	Link      uint   `json:"link"`
	Diagnosis string `json:"diagnosis"`
}

type SelectDoctorRequestDTO struct {
	DocUsername string `json:"doc_username" binding:"required"`
}

type DocRespondRequestDTO struct {
	PatientID uint   `json:"patient_id" binding:"required"`
	Action    string `json:"action" binding:"required"` // "accept" | "reject"
}

type NoteDTO struct {
	ID               uint      `json:"id" binding:"required"`
	DateRecorded     time.Time `json:"date_recorded" binding:"required"`
	Intensity        int       `json:"intensity" binding:"required"`
	PainType         string    `json:"pain_type" binding:"required"`
	TookPrescription bool      `json:"tookPrescription" binding:"required"`
	Description      string    `json:"description" binding:"required"`
	BodyPart         int       `json:"body_part" binding:"required"`
}
