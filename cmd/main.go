package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"socialmedia-backend/internal/Auth/handlers"
	"socialmedia-backend/internal/shared/db"
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

	http.HandleFunc("/api/health", enableCORS(healthCheck))

	http.HandleFunc("/api/verify/login", enableCORS(authHandler.Login))
	http.HandleFunc("/api/verify/register", enableCORS(authHandler.Register))

	fmt.Printf("Server starting on http://localhost:%s\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
