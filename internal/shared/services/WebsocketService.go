package services

import (
	"encoding/json"
	"sync"

	"github.com/gorilla/websocket"
)

type MessageNotification struct {
	ConversationID uint   `json:"conversation_id"`
	SenderID       uint   `json:"sender_id"`
	SenderUsername string `json:"sender_username"`
	Content        string `json:"content"`
}

type WebSocketService struct {
	connections map[uint]*websocket.Conn
	mutex       sync.RWMutex
}

func NewWebSocketService() *WebSocketService {
	return &WebSocketService{
		connections: make(map[uint]*websocket.Conn),
	}
}

func (s *WebSocketService) RegisterConnection(userID uint, conn *websocket.Conn) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.connections[userID] = conn
}

func (s *WebSocketService) RemoveConnection(userID uint) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	delete(s.connections, userID)
}

func (s *WebSocketService) NotifyUsers(notification MessageNotification, recipientIDs []uint) {
	message, err := json.Marshal(notification)
	if err != nil {
		return
	}

	s.mutex.RLock()
	defer s.mutex.RUnlock()

	for _, userID := range recipientIDs {
		if conn, exists := s.connections[userID]; exists {
			conn.WriteMessage(websocket.TextMessage, message)
		}
	}
}
