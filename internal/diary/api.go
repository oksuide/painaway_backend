package diary

import (
	"fmt"
	"net/http"
	"painaway_test/internal/utils"
	"painaway_test/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Handler struct {
	Service *Service
	Logger  *zap.Logger
}

func RegisterRoutes(rg *gin.RouterGroup, service *Service, logger *zap.Logger) {
	h := &Handler{Service: service, Logger: logger}
	rg.GET("/diary/list_links", h.ListLinks)
	rg.POST("/diary/link_doc/", h.LinkDoc)
	rg.POST("/diary/doc_respond", h.DocRespond)
	rg.GET("/diary/stats/", h.GetUserBodyStats)
	rg.GET("/diary/bodyparts/", h.GetBodyParts)
	rg.POST("/diary/stats/", h.CreateNote)
	rg.POST("/diary/diagnosis", h.SetDiagnosis)
	rg.POST("/diary/prescription", h.SetPrescription)
}

func (h *Handler) ListLinks(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		h.Logger.Warn("unauthorized access to ListLinks")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userGroups, _ := c.Get("groups")

	switch userGroups {
	case "Doctor":
		links, err := h.Service.DoctorListLinks(userID.(uint))
		if err != nil {
			h.Logger.Error("failed to list links", zap.Uint("userID", userID.(uint)), zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		h.Logger.Info("links retrieved", zap.Uint("userID", userID.(uint)), zap.Int("count", len(links)))
		c.JSON(http.StatusOK, links)

	case "Patient":
		links, err := h.Service.PatientListLinks(userID.(uint))
		if err != nil {
			h.Logger.Error("failed to list links", zap.Uint("userID", userID.(uint)), zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		h.Logger.Info("links retrieved", zap.Uint("userID", userID.(uint)), zap.Int("count", len(links)))
		c.JSON(http.StatusOK, links)

	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid groups"})
		return
	}

}

func (h *Handler) LinkDoc(c *gin.Context) {
	var req utils.SelectDoctorRequestDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		h.Logger.Warn("invalid request body in LinkDoc", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	patientID, exists := c.Get("userID")
	if !exists {
		h.Logger.Warn("unauthorized access to LinkDoc")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	link, err := h.Service.LinkDoc(patientID.(uint), req.DocUsername)
	if err != nil {
		h.Logger.Error("failed to link doctor", zap.Uint("patientID", patientID.(uint)), zap.String("docUsername", req.DocUsername), zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	h.Logger.Info("doctor linked successfully", zap.Uint("patientID", patientID.(uint)), zap.String("docUsername", req.DocUsername))
	c.JSON(http.StatusOK, link)
}

func (h *Handler) DocRespond(c *gin.Context) {
	var req utils.DocRespondRequestDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		h.Logger.Warn("invalid request body in DocRespond", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	doctorID, exists := c.Get("userID")
	if !exists {
		h.Logger.Warn("unauthorized access to DocRespond")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if err := h.Service.RespondToLinkRequest(doctorID.(uint), req.PatientID, req.Action); err != nil {
		h.Logger.Error("failed to respond to link request", zap.Uint("doctorID", doctorID.(uint)), zap.Uint("patientID", req.PatientID), zap.String("action", req.Action), zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	h.Logger.Info("link request responded", zap.Uint("doctorID", doctorID.(uint)), zap.Uint("patientID", req.PatientID), zap.String("action", req.Action))
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *Handler) SetPrescription(c *gin.Context) {
	var req utils.SetPrescriptionDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		h.Logger.Warn("invalid request body in SetPrescription", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.Service.SetPrescription(req); err != nil {
		h.Logger.Error("failed to respond to set prescription", zap.Uint("linkID", uint(req.Link)), zap.String("prescription", req.Prescription), zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	h.Logger.Info("prescription set successfully", zap.Uint("linkID", uint(req.Link)), zap.String("prescription", req.Prescription))
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *Handler) SetDiagnosis(c *gin.Context) {
	var req utils.SetDiagnosisDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		h.Logger.Warn("invalid request body in SetDiagnosis", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.Service.SetDiagnosis(req); err != nil {
		h.Logger.Error("failed to respond to set diagnosis", zap.Uint("linkID", uint(req.Link)), zap.String("diagnosis", req.Diagnosis), zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	h.Logger.Info("diagnosis set successfully", zap.Uint("linkID", uint(req.Link)), zap.String("diagnosis", req.Diagnosis))
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *Handler) GetUserBodyStats(c *gin.Context) {
	userID, err := h.resolveUserID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	stats, err := h.Service.GetUserAllBodyStats(userID)
	if err != nil {
		h.Logger.Warn("failed to get body stats", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get body stats"})
		return
	}

	c.JSON(http.StatusOK, stats)
}

func (h *Handler) GetBodyParts(c *gin.Context) {
	bodyParts := h.Service.GetBodyParts()
	c.JSON(http.StatusOK, bodyParts)
}

func (h *Handler) CreateNote(c *gin.Context) {
	patientID, exists := c.Get("userID")
	if !exists {
		h.Logger.Warn("unauthorized access to CreateNote")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req models.Note
	if err := c.ShouldBindJSON(&req); err != nil {
		h.Logger.Warn("invalid request body in CreateNote", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.PatientID = patientID.(uint)

	if err := h.Service.CreateNote(&req); err != nil {
		h.Logger.Error("failed to create note", zap.Uint("patientID", patientID.(uint)), zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.Logger.Info("note created", zap.Uint("patientID", patientID.(uint)), zap.Uint("noteID", req.ID))
	c.JSON(http.StatusOK, gin.H{"message": "Note created"})
}

func (h *Handler) resolveUserID(c *gin.Context) (uint, error) {
	idStr := c.Query("patient_id")

	if idStr != "" {
		idUint64, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			return 0, fmt.Errorf("invalid id query parameter: %w", err)
		}
		userID := uint(idUint64)
		return userID, nil
	}

	userIDValue, exists := c.Get("userID")
	if !exists {
		return 0, fmt.Errorf("unauthorized")
	}

	uid, ok := userIDValue.(uint)
	if !ok {
		return 0, fmt.Errorf("invalid userID type in context")
	}

	return uid, nil
}
