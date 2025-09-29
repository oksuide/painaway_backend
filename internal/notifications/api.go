package notifications

import (
	"net/http"
	"painaway_test/internal/response"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type Handler struct {
	Service *Service
	Hub     *Hub
	Logger  *zap.Logger
}

func RegisterRoutes(rg *gin.RouterGroup, service *Service, hub *Hub, logger *zap.Logger) {
	h := &Handler{Service: service, Logger: logger, Hub: hub}
	rg.GET("/diary/notifications/", h.GetNotifications)
	rg.PATCH("/diary/notifications/", h.MarkNotificationRead)
	rg.DELETE("/diary/notifications/", h.DeleteNotification)
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func (h *Handler) WsNotifications(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.NewErrorResponse(c, http.StatusUnauthorized, "unauthorized", h.Logger)
		return
	}
	uid := userID.(uint)

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		h.Logger.Error("failed to upgrade websocket connection", zap.Uint("userID", uid), zap.Error(err))
		response.NewErrorResponse(c, http.StatusInternalServerError, "failed to establish websocket connection", h.Logger)
		return
	}
	defer conn.Close()

	h.Hub.Register(uid, conn)
	h.Logger.Info("user connected to notifications hub", zap.Uint("userID", uid))
	defer func() {
		h.Hub.Unregister(uid)
		h.Logger.Info("user disconnected from notifications hub", zap.Uint("userID", uid))
	}()

	for {
		if _, _, err := conn.NextReader(); err != nil {
			break
		}
	}
}

func (h *Handler) GetNotifications(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.NewErrorResponse(c, http.StatusUnauthorized, "unauthorized", h.Logger)
		return
	}

	notifications, err := h.Service.GetNotifications(userID.(uint))
	if err != nil {
		h.Logger.Error("failed to get notifications", zap.Uint("userID", userID.(uint)), zap.Error(err))
		response.NewErrorResponse(c, http.StatusInternalServerError, "failed to fetch notifications", h.Logger)
		return
	}

	h.Logger.Info("notifications retrieved", zap.Uint("userID", userID.(uint)), zap.Int("count", len(notifications)))
	c.JSON(http.StatusOK, notifications)
}

func (h *Handler) MarkNotificationRead(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.NewErrorResponse(c, http.StatusUnauthorized, "unauthorized", h.Logger)
		return
	}

	var req struct {
		NotificationID uint `json:"notification_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, "invalid request body", h.Logger)
		return
	}

	if err := h.Service.MarkNotificationRead(req.NotificationID, userID.(uint)); err != nil {
		h.Logger.Error("failed to mark notification as read", zap.Uint("userID", userID.(uint)), zap.Uint("notificationID", req.NotificationID), zap.Error(err))
		response.NewErrorResponse(c, http.StatusInternalServerError, "failed to mark notification as read", h.Logger)
		return
	}
	h.Logger.Info("notification marked as read", zap.Uint("userID", userID.(uint)), zap.Uint("notificationID", req.NotificationID))
	c.Status(http.StatusNoContent)
}

func (h *Handler) DeleteNotification(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.NewErrorResponse(c, http.StatusUnauthorized, "unauthorized", h.Logger)
		return
	}
	var req struct {
		NotificationID uint `json:"notification_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, "invalid request body", h.Logger)
		return
	}

	if err := h.Service.DeleteNotification(req.NotificationID, userID.(uint)); err != nil {
		h.Logger.Error("failed to delete notification", zap.Uint("userID", userID.(uint)), zap.Uint("notificationID", req.NotificationID), zap.Error(err))
		response.NewErrorResponse(c, http.StatusInternalServerError, "failed to delete notification", h.Logger)
		return
	}

	h.Logger.Info("notification deleted", zap.Uint("userID", userID.(uint)), zap.Uint("notificationID", req.NotificationID))
	c.Status(http.StatusNoContent)
}
