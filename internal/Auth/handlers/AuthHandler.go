package handlers

import (
	"socialmedia-backend/internal/Auth/usecases"
	"encoding/json"
	"net/http"
)

type AuthHandler struct {
	loginUsecase    *usecases.LoginUsecase
	registerUsecase *usecases.RegisterUsecase
}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{
		loginUsecase:    usecases.NewLoginUsecase(),
		registerUsecase: usecases.NewRegisterUsecase(),
	}
}

// Login handles the /api/auth/login endpoint
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {

	type LoginResponse struct {
		ID        uint   `json:"id"`
		Username  string `json:"username"`
		Email     string `json:"email"`
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := h.loginUsecase.Execute(request.Email, request.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	response := LoginResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Register handles the /api/auth/register endpoint
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var request struct {
		Username  string `json:"username"`
		Email     string `json:"email"`
		Password  string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := h.registerUsecase.Execute(request.Username, request.Email, request.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Registration successful"})
}
