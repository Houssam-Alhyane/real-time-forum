package handlers

import (
	"net/http"
)



func HandleError(w http.ResponseWriter, code int, message string) {
	RespondJSON(w, code, map[string]string{
		"error": message,
	})
}