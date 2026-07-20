package middlewares

import (
	"net/http"
	"time"
	"zone/backend/database"
	"zone/backend/handlers"
)

// Auth is a middleware that restricts access to authenticated users.
// It verifies the user's session token and updates their last seen activity.
func Auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := handlers.GetUserIDFromSession(r)
		if err != nil {
			http.SetCookie(w, &http.Cookie{
				Name:     "session_token",
				Value:    "",
				Path:     "/",
				Expires:  time.Unix(0, 0),
				HttpOnly: true,
				SameSite: http.SameSiteLaxMode,
			})
			handlers.HandleError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		// Update activity timestamp
		database.UpdateLastSeen(userID)

		next(w, r)
	}
}
