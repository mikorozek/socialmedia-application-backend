package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"socialmedia-backend/internal/handlers"
	"socialmedia-backend/internal/shared/db"
	"socialmedia-backend/internal/shared/services"
	"time"
)

func enableCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		frontendHost := os.Getenv("FRONTEND_HOST")
		if frontendHost == "" {
			frontendHost = "localhost"
		}

		frontendPort := os.Getenv("FRONTEND_PORT")
		if frontendPort == "" {
			frontendPort = "3000"
		}

		allowedOrigin := fmt.Sprintf("http://%s:%s", frontendHost, frontendPort)
		w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	dbConn := db.GetDB()
	sqlDB, err := dbConn.DB()

	status := "healthy"
	dbStatus := "up"

	if err != nil || sqlDB.Ping() != nil {
		status = "unhealthy"
		dbStatus = "down"
	}

	response := map[string]interface{}{
		"status":    status,
		"timestamp": time.Now().UTC(),
		"services": map[string]string{
			"database": dbStatus,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	port := os.Getenv("BACKEND_PORT")
	if port == "" {
		port = "8080"
	}

	db.InitDB()

	authHandler := handlers.NewAuthHandler()
	wsService := services.NewWebSocketService()
	wsHandler := handlers.NewWebSocketHandler(wsService)
	conversationHandler := handlers.NewConversationHandler(wsService)

	// Health check endpoint
	http.HandleFunc("/api/health", enableCORS(healthCheck))
	http.HandleFunc("/ws", enableCORS(wsHandler.HandleWebSocket))

	// Auth endpoints
	http.HandleFunc("/api/verify/login", enableCORS(authHandler.Login))
	http.HandleFunc("/api/verify/register", enableCORS(authHandler.Register))

	// Conversation endpoints
	http.HandleFunc("/api/conversations/create", enableCORS(conversationHandler.CreateConversation))
	http.HandleFunc("/api/conversations/messages", enableCORS(conversationHandler.SendMessage))     // POST do wysyłania
	http.HandleFunc("/api/conversations/messages/get", enableCORS(conversationHandler.GetMessages)) // GET do pobierania
	http.HandleFunc("/api/conversations/messages/edit", enableCORS(conversationHandler.EditMessage))
	//http.HandleFunc("/api/conversations/messages/delete", enableCORS(conversationHandler.DeleteMessage))
	http.HandleFunc("/api/conversations/recent", enableCORS(conversationHandler.GetRecentConversations))
	http.HandleFunc("/api/conversations/unread", enableCORS(conversationHandler.GetUnreadConversations))
	http.HandleFunc("/api/conversations/mark-read", enableCORS(conversationHandler.MarkAsRead))

	fmt.Printf("Server starting on http://localhost:%s\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
