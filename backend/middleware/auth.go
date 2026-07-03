package middlewares

import (
	"database/sql"
	"net/http"
	"time"

	"zone/backend/database"
	"zone/backend/handlers"
)

func Auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		cookie, err := r.Cookie("session_token")
		if err != nil {
			handlers.HandleError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		var (
			userID  int
			expires time.Time
		)

		err = database.Database.QueryRow(`
			SELECT user_id, expires_at
			FROM sessions
			WHERE id = ?
		`, cookie.Value).Scan(&userID, &expires)

		if err == sql.ErrNoRows {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		if expires.Before(time.Now()) {
			http.Error(w, "Session expired", http.StatusUnauthorized)
			return
		}

		// Update activity
		database.UpdateLastSeen(userID)

		next(w, r)
	}
}