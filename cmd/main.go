package main

import (
	"fmt"
)

// func enableCORS(next http.HandlerFunc) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		w.Header().Set("Access-Control-Allow-Origin", "*")
// 		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
// 		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

// 		if r.Method == "OPTIONS" {
// 			w.WriteHeader(http.StatusOK)
// 			return
// 		}
// 		next(w, r)
// 	}
// }

func main() {
	// // Get port from environment variable or use default
	// port := os.Getenv("PORT")
	// if port == "" {
	// 	port = "8080"
	// }

	// // Initialize database connection
	// if err := db.InitDB(); err != nil {
	// 	log.Fatal("Failed to initialize database:", err)
	// }

	// // Initialize handlers
	// // Handler := handlers.NewAuthHandler()

	// // http.HandleFunc("/api/auth/login", enableCORS(authHandler.Login))
	// // http.HandleFunc("/api/auth/register", enableCORS(authHandler.Register))
	// // http.HandleFunc("/api/auth/verify", enableCORS(authHandler.VerifyToken))

	// // Start server
	fmt.Printf("Server starting on http://localhost\n")

}
