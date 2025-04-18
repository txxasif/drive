package handler

import (
	"encoding/json"
	"net/http"
)

// HealthHandler is a simple health check handler.
// It returns a JSON response with a success message.
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{
		"status":  "ok",
		"message": "Server is running successfully!",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
