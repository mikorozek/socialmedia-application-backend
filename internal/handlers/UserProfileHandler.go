package handlers

import (
	"encoding/json"
	"net/http"
	"socialmedia-backend/internal/usecases"
	"strconv"
)

type UserProfileHandler struct {
	profileUsecase *usecases.UserProfileUsecase
}

func NewUserProfileHandler() *UserProfileHandler {
	return &UserProfileHandler{
		profileUsecase: usecases.NewUserProfileUsecase(),
	}
}

// GET /api/users/profile?user_id=1
func (h *UserProfileHandler) GetUserProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, err := strconv.ParseUint(r.URL.Query().Get("user_id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	profile, err := h.profileUsecase.GetUserProfile(uint(userID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(profile)
}

type EditProfileRequest struct {
	UserID      uint   `json:"user_id"`
	Username    string `json:"username,omitempty"`
	Password    string `json:"password,omitempty"`
	Description string `json:"description,omitempty"`
}

// POST /api/users/profile/edit
func (h *UserProfileHandler) EditUserProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var rawRequest map[string]interface{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&rawRequest); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var request EditProfileRequest
	requestData, _ := json.Marshal(rawRequest)
	if err := json.Unmarshal(requestData, &request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := h.profileUsecase.UpdateUserProfile(
		request.UserID,
		request.Username,
		request.Password,
		request.Description,
		rawRequest,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}
