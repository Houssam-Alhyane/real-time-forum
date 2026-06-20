package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"
	"zone/backend/database"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func Login(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/login" {
		HandleError(w, http.StatusNotFound, "Page not found")
		return
	}

	switch r.Method {
	case http.MethodGet:
		RenderTemplate(w, 200, "index.html", nil)
		return

	case http.MethodPost:
		// Map the field name to match the JavaScript frontend payload ("login")
		identifier := strings.TrimSpace(r.FormValue("login"))
		password := r.FormValue("password")

		// Basic input presence validation
		if identifier == "" || password == "" {
			HandleError(w, http.StatusBadRequest, "Email/username and password are required")
			return
		}

		var userID int
		var hashedPassword string

		// Search for the user by either their email or nickname
		err := database.Database.QueryRow(
			"SELECT id, password FROM users WHERE email = ? OR nickname = ?", identifier, identifier,
		).Scan(&userID, &hashedPassword)
		if err != nil {
			HandleError(w, http.StatusUnauthorized, "Invalid email/username or password")
			return
		}

		// Verify if the provided password matches the hashed password in the database
		if err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
			HandleError(w, http.StatusUnauthorized, "Invalid email/username or password")
			return
		}

		// Delete any existing active sessions for this specific user
		_, err = database.Database.Exec("DELETE FROM sessions WHERE user_id = ?", userID)
		if err != nil {
			HandleError(w, http.StatusInternalServerError, "Server error")
			return
		}

		// Generate a new unique session ID and set the expiration time
		sessionID := uuid.New().String()
		expiration := time.Now().Add(24 * time.Hour)

		// Insert the new session record into the database
		_, err = database.Database.Exec(
			"INSERT INTO sessions (id, expires_at, user_id) VALUES (?, ?, ?)",
			sessionID, expiration, userID,
		)
		if err != nil {
			HandleError(w, http.StatusInternalServerError, "Server error")
			return
		}


		http.SetCookie(w, &http.Cookie{
			Name:     "session_token",
			Value:    sessionID,
			Path:     "/",
			Expires:  expiration,
			SameSite: http.SameSiteLaxMode,
			// Secure: true, // enable once served over HTTPS
		})

		// A non-secret flag, readable by document.cookie, that the
		// frontend uses purely to decide which UI to render.
		http.SetCookie(w, &http.Cookie{
			Name:    "logged_in",
			Value:   "true",
			Path:    "/",
			Expires: expiration,
		})

		// Return a successful JSON response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "login successfully",
		})

	default:
		HandleError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}