package handlers

import (
	"database/sql"
	"net/http"
	"zone/backend/database"
)

//take userid in database with cookie

func GetUserIDFromSession(r *http.Request) (int, error) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		return 0, sql.ErrNoRows
	}

	var userID int
	err = database.Database.QueryRow(
		"SELECT user_id FROM sessions WHERE id = ? AND expires_at > DATETIME('now')",
		cookie.Value,
	).Scan(&userID)
	return userID, err
}