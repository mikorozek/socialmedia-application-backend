package handlers

import (
	"encoding/json"
	"net/http"
	"socialmedia-backend/internal/usecases"
)

type UserSearchHandler struct {
	searchUsecase *usecases.UserSearchUsecase
}

func NewUserSearchHandler() *UserSearchHandler {
	return &UserSearchHandler{
		searchUsecase: usecases.NewUserSearchUsecase(),
	}
}

// GET /api/users/search?query=abc
func (h *UserSearchHandler) SearchUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	query := r.URL.Query().Get("query")
	if query == "" {
		http.Error(w, "Query parameter is required", http.StatusBadRequest)
		return
	}

	results, err := h.searchUsecase.SearchUsers(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}
