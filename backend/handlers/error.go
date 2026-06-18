package handlers

import (
	"encoding/json"
	"net/http"
)

func HandleError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	json.NewEncoder(w).Encode(map[string]string{
		"error": message,
	})
}