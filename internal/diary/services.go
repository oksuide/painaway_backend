package diary

import (
	"painaway_test/internal/notifications"
	"painaway_test/internal/utils"
	"painaway_test/models"

	"go.uber.org/zap"
)

type Service struct {
	Repo                 Repository
	NotificationsService *notifications.Service
	Logger               *zap.Logger
}

func NewService(repo Repository, notifSrv *notifications.Service, logger *zap.Logger) *Service {
	return &Service{
		Repo:                 repo,
		NotificationsService: notifSrv,
		Logger:               logger,
	}
}

func (s *Service) PatientListLinks(userID uint) ([]utils.PatientLinkDTO, error) {
	subs, err := s.Repo.GetSubscriptionsByPatientID(userID, 0, 20)
	if err != nil {
		return nil, err
	}

	if len(subs) == 0 {
		return []utils.PatientLinkDTO{
			{
				ID:     0,
				Status: "none",
				Doctor: utils.DoctorDTO{},
			},
		}, nil
	}

	var result []utils.PatientLinkDTO
	for _, sub := range subs {
		dto := utils.PatientLinkDTO{
			ID:     sub.ID,
			Status: sub.Status,
			Doctor: utils.DoctorDTO{
				ID:         sub.Doctor.ID,
				Username:   sub.Doctor.Username,
				LastName:   sub.Doctor.LastName,
				FirstName:  sub.Doctor.FirstName,
				FatherName: sub.Doctor.FatherName,
			},
			Prescription: sub.Prescription,
		}
		result = append(result, dto)
	}

	return result, nil
}

func (s *Service) DoctorListLinks(userID uint) ([]utils.DoctorLinkDTO, error) {
	subs, err := s.Repo.GetSubscriptionsByDoctorID(userID, 0, 20)
	if err != nil {
		return nil, err
	}

	if len(subs) == 0 {
		return []utils.DoctorLinkDTO{
			{
				ID:     0,
				Status: "none",
			},
		}, nil
	}

	var result []utils.DoctorLinkDTO
	for _, sub := range subs {
		dto := utils.DoctorLinkDTO{
			ID:           sub.ID,
			Status:       sub.Status,
			Prescription: utils.SetPrescriptionDTO{Prescription: sub.Prescription},
			Diagnosis:    utils.SetDiagnosisDTO{Diagnosis: sub.Diagnosis},
			Patient: utils.PatientDTO{
				ID:          sub.Patient.ID,
				FirstName:   sub.Patient.FirstName,
				LastName:    sub.Patient.LastName,
				FatherName:  sub.Patient.FatherName,
				Sex:         sub.Patient.Sex,
				DateOfBirth: sub.Patient.DateOfBirth.Format("02.01.2006"),
			},
		}
		result = append(result, dto)
	}

	return result, nil
}

func (s *Service) LinkDoc(patientID uint, docUsername string) (*utils.PatientLinkDTO, error) {
	doc, err := s.Repo.GetDoctorByUsername(docUsername)
	if err != nil {
		return nil, err
	}

	sub := &models.Subscription{
		PatientID: patientID,
		DoctorID:  doc.ID,
		Status:    "pending",
	}

	if err := s.Repo.CreateSubscription(sub); err != nil {
		return nil, err
	}
	dto := &utils.PatientLinkDTO{
		ID:     sub.ID,
		Status: sub.Status,
		Doctor: utils.DoctorDTO{
			ID:         doc.ID,
			Username:   doc.Username,
			LastName:   doc.LastName,
			FirstName:  doc.FirstName,
			FatherName: doc.FatherName,
		},
	}
	if err := s.NotificationsService.CreateNotification(
		sub.DoctorID,
		"Новый запрос на прикрепление",
	); err != nil {
		s.Logger.Error("failed to create notification",
			zap.Uint("doctorID", sub.DoctorID),
			zap.Uint("patientID", patientID),
			zap.Error(err),
		)
	}

	return dto, nil
}

func (s *Service) RespondToLinkRequest(doctorID, patientID uint, action string) error {
	link, err := s.Repo.GetLinkByDoctorAndPatient(doctorID, patientID)
	if err != nil {
		return err
	}

	switch action {
	case "accept":
		link.Status = "accepted"
	case "reject":
		link.Status = "rejected"
	default:
		return err
	}

	if err := s.NotificationsService.CreateNotification(
		patientID,
		"Ответ на запрос о прикреплении",
	); err != nil {
		s.Logger.Error("failed to create notification",
			zap.Uint("doctorID", doctorID),
			zap.Uint("patientID", patientID),
			zap.Error(err),
		)
	}
	return s.Repo.UpdateLink(link)
}

func (s *Service) SetPrescription(req utils.SetPrescriptionDTO) error {
	link, err := s.Repo.GetLinkByID(req.Link)
	if err != nil {
		return err
	}
	link.Prescription = req.Prescription

	return s.Repo.UpdateLink(link)
}

func (s *Service) SetDiagnosis(req utils.SetDiagnosisDTO) error {
	link, err := s.Repo.GetLinkByID(req.Link)
	if err != nil {
		return err
	}
	link.Diagnosis = req.Diagnosis

	return s.Repo.UpdateLink(link)
}

func (s *Service) GetUserAllStats(patientID uint) ([]models.Note, error) {
	stats, err := s.Repo.GetAllStatByPatientID(patientID)
	if err != nil {
		return nil, err
	}

	if len(stats) == 0 {
		return []models.Note{}, nil
	}

	return stats, nil
}

func (s *Service) GetBodyParts() []BodyPart {
	return BodyParts
}

func (s *Service) CreateNote(note *models.Note) error {
	return s.Repo.CreateNote(note)
}

func (s *Service) ToNoteDTO(notes []models.Note) []utils.NoteDTO {
	dto := make([]utils.NoteDTO, 0, len(notes))

	for _, n := range notes {
		dto = append(dto, utils.NoteDTO{
			ID:               n.ID,
			DateRecorded:     n.CreatedAt,
			Intensity:        n.Intensity,
			PainType:         n.PainType,
			TookPrescription: n.TookPrescription,
			Description:      n.Description,
			BodyPart:         int(n.BodyPart),
		})
	}

	return dto
}
