package handlers

import (
	"encoding/json"
	"net/http"
	"socialmedia-backend/internal/shared/services"
	"socialmedia-backend/internal/usecases"
	"strconv"
	"time"
)

type ConversationHandler struct {
	conversationUsecase *usecases.ConversationUsecase
	wsService           *services.WebSocketService
}

func NewConversationHandler(wsService *services.WebSocketService) *ConversationHandler {
	return &ConversationHandler{
		conversationUsecase: usecases.NewConversationUsecase(wsService),
		wsService:           wsService,
	}
}

// POST /api/conversations/create
func (h *ConversationHandler) CreateConversation(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request struct {
		UserIDs []uint `json:"user_ids"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	conversation, err := h.conversationUsecase.CreateConversation(request.UserIDs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(conversation)
}

// POST /api/conversations/messages
func (h *ConversationHandler) SendMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request struct {
		ConversationID uint   `json:"conversation_id"`
		SenderID       uint   `json:"sender_id"`
		Content        string `json:"content"`
		PhotoURL       string `json:"photo_url,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	message, err := h.conversationUsecase.SendMessage(
		request.ConversationID,
		request.SenderID,
		request.Content,
		request.PhotoURL,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(message)
}

// GET /api/conversations/messages?conversation_id=1&user_id=1&limit=50&offset=0
func (h *ConversationHandler) GetMessages(w http.ResponseWriter, r *http.Request) {
	type MessageResponse struct {
		ID             uint      `json:"id"`
		ConversationID uint      `json:"conversation_id"`
		UserID         uint      `json:"user_id"`
		Content        string    `json:"content"`
		MessageDate    time.Time `json:"message_date"`
		PhotoURL       string    `json:"photo_url"`
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	conversationID, err := strconv.ParseUint(r.URL.Query().Get("conversation_id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid conversation ID", http.StatusBadRequest)
		return
	}

	userID, err := strconv.ParseUint(r.URL.Query().Get("user_id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	if limit == 0 {
		limit = 50 // Default limit
	}

	messages, err := h.conversationUsecase.GetConversationMessages(uint(conversationID), uint(userID), limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := make([]MessageResponse, 0)

	for _, msg := range messages {
		response = append(response, MessageResponse{
			ID:             msg.ID,
			ConversationID: msg.ConversationID,
			UserID:         msg.UserID,
			Content:        msg.Content,
			MessageDate:    msg.MessageDate,
			PhotoURL:       msg.PhotoURL,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// POST /api/conversations/messages/edit
func (h *ConversationHandler) EditMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request struct {
		MessageID uint   `json:"message_id"`
		UserID    uint   `json:"user_id"`
		Content   string `json:"content"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.conversationUsecase.EditMessage(request.MessageID, request.UserID, request.Content); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// POST /api/conversations/mark-read
func (h *ConversationHandler) MarkAsRead(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request struct {
		ConversationID uint `json:"conversation_id"`
		UserID         uint `json:"user_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.conversationUsecase.MarkConversationAsRead(request.ConversationID, request.UserID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// GET /api/conversations/unread?user_id=1
func (h *ConversationHandler) GetUnreadConversations(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, err := strconv.ParseUint(r.URL.Query().Get("user_id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	conversations, err := h.conversationUsecase.GetUnreadConversations(uint(userID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(conversations)
}

// GET /api/conversations/recent?user_id=1&limit=50
func (h *ConversationHandler) GetRecentConversations(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, err := strconv.ParseUint(r.URL.Query().Get("user_id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit == 0 {
		limit = 50
	}

	conversations, err := h.conversationUsecase.GetRecentConversations(uint(userID), limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Struktura dla odpowiedzi API
	type UserResponse struct {
		ID       uint   `json:"id"`
		Username string `json:"username"`
		Email    string `json:"email"`
	}

	type MessageResponse struct {
		ID          uint      `json:"id"`
		Content     string    `json:"content"`
		MessageDate time.Time `json:"message_date"`
		UserID      uint      `json:"user_id"`
	}

	type ConversationResponse struct {
		ID          uint             `json:"id"`
		Users       []UserResponse   `json:"users"`
		LastMessage *MessageResponse `json:"last_message"`
	}

	// Przekształcamy dane do żądanego formatu
	response := make([]ConversationResponse, 0)
	for _, conv := range conversations {
		users := make([]UserResponse, 0)
		for _, user := range conv.Users {
			users = append(users, UserResponse{
				ID:       user.ID,
				Username: user.Username,
				Email:    user.Email,
			})
		}

		var lastMessage *MessageResponse
		if len(conv.Messages) > 0 && conv.Messages[0].ID != 0 {
			lastMessage = &MessageResponse{
				ID:          conv.Messages[0].ID,
				Content:     conv.Messages[0].Content,
				MessageDate: conv.Messages[0].MessageDate,
				UserID:      conv.Messages[0].UserID,
			}
		}

		response = append(response, ConversationResponse{
			ID:          conv.ID,
			Users:       users,
			LastMessage: lastMessage,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
