package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"zone/backend/database"
	"zone/backend/types"
)

func ReactToComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		HandleError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	userID, err := GetUserIDFromSession(r)
	if err != nil {
		HandleError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req types.CommentReactRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		HandleError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	if req.Type != "like" && req.Type != "dislike" {
		HandleError(w, http.StatusBadRequest, "type must be 'like' or 'dislike'")
		return
	}

	newIsLike := 0
	if req.Type == "like" {
		newIsLike = 1
	}

	exists, err := database.DoesCommentExist(req.CommentID)
	if err != nil {
		HandleError(w, http.StatusInternalServerError, "Database error")
		return
	}
	if !exists {
		HandleError(w, http.StatusBadRequest, "Invalid comment")
		return
	}

	var existingIsLike int
	err = database.Database.QueryRow(
		"SELECT is_like FROM comment_reactions WHERE user_id = ? AND comment_id = ?",
		userID, req.CommentID,
	).Scan(&existingIsLike)

	switch {
	case err == sql.ErrNoRows:
		_, err = database.Database.Exec(
			"INSERT INTO comment_reactions (user_id, comment_id, is_like) VALUES (?, ?, ?)",
			userID, req.CommentID, newIsLike,
		)
	case err == nil && existingIsLike == newIsLike:
		_, err = database.Database.Exec(
			"DELETE FROM comment_reactions WHERE user_id = ? AND comment_id = ?",
			userID, req.CommentID,
		)
	case err == nil:
		_, err = database.Database.Exec(
			"UPDATE comment_reactions SET is_like = ? WHERE user_id = ? AND comment_id = ?",
			newIsLike, userID, req.CommentID,
		)
	}
	if err != nil {
		log.Printf("ReactToComment: %v", err)
		HandleError(w, http.StatusInternalServerError, "Could not save reaction")
		return
	}

	resp, err := database.GetCommentReactionSummary(req.CommentID, userID)
	if err != nil {
		log.Printf("ReactToComment summary: %v", err)
		HandleError(w, http.StatusInternalServerError, "Could not load reaction counts")
		return
	}

	RespondJSON(w, http.StatusOK, resp)
}
