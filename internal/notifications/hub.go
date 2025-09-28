package notifications

import (
	"sync"

	"painaway_test/models"

	"github.com/gorilla/websocket"
)

type Hub struct {
	mu      sync.RWMutex
	clients map[uint]*websocket.Conn // userID → conn
}

func NewHub() *Hub {
	return &Hub{
		clients: make(map[uint]*websocket.Conn),
	}
}

func (h *Hub) Register(userID uint, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.clients[userID] = conn
}

func (h *Hub) Unregister(userID uint) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.clients, userID)
}

func (h *Hub) Send(userID uint, notification *models.Notification) error {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if conn, ok := h.clients[userID]; ok {
		return conn.WriteJSON(notification)
	}
	return nil
}
