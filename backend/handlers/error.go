package handlers

import (
	"net/http"
	"os"
	"strconv"
	"strings"
)



func HandleError(w http.ResponseWriter, code int, message string) {
	RespondJSON(w, code, map[string]string{
		"error": message,
	})
}

func IsValidSPARoute(path string) bool {
	if path == "/" || path == "/login" || path == "/register" {
		return true
	}
	if strings.HasPrefix(path, "/posts/") {
		idStr := path[len("/posts/"):]
		if idStr != "" {
			_, err := strconv.Atoi(idStr)
			return err == nil
		}
	}
	return false
}

func ServeSPAIndex(w http.ResponseWriter, status int) {
	content, err := os.ReadFile("./frontend/index.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 Internal Server Error"))
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	w.Write(content)
}