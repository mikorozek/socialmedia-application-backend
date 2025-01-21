package handlers

import (
	"net/http"
	"socialmedia-backend/internal/shared/services"
	"strconv"

	"github.com/gorilla/websocket"
)

type WebSocketHandler struct {
	wsService *services.WebSocketService
	upgrader  websocket.Upgrader
}

func NewWebSocketHandler(wsService *services.WebSocketService) *WebSocketHandler {
	return &WebSocketHandler{
		wsService: wsService,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}

// HandleWebSocket obsługuje nowe połączenia WebSocket
func (h *WebSocketHandler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Pobierz user_id z query params
	userIDStr := r.URL.Query().Get("user_id")
	if userIDStr == "" {
		http.Error(w, "Missing user_id", http.StatusBadRequest)
		return
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid user_id", http.StatusBadRequest)
		return
	}

	// Upgrade połączenia do WebSocket
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Could not upgrade connection", http.StatusInternalServerError)
		return
	}

	// Zarejestruj połączenie
	h.wsService.RegisterConnection(uint(userID), conn)

	// Czekaj na zamknięcie połączenia
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			h.wsService.RemoveConnection(uint(userID))
			return
		}
	}
}
