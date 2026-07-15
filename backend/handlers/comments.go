package handlers

import (
	"net/http"
	"strconv"
	"zone/backend/database"
)

type CommentPayload struct {
	PostID  int    `json:"post_id"`
	Content string `json:"content"`
	UserID  int    `json:"user_id"`
}

// GetCommentsAPI fetches comments with pagination for a specific post.
func GetCommentsAPI(w http.ResponseWriter, r *http.Request) {
}

// CreateCommentAPI handles the creation of a new comment for a specific post.
func CreateCommentAPI(w http.ResponseWriter, r *http.Request) {
	// Restrict the request method to POST
	if r.Method != http.MethodPost {
		HandleError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// check authentication
	payload := CommentPayload{}
	userID, err := GetUserIDFromSession(r)
	if err != nil || userID == 0 {
		HandleError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	payload.UserID = userID
	// Parse Payload: Read post_id and content from the request
	payload.PostID, err = strconv.Atoi(r.FormValue("post_id"))
	// validate post_id and content
	if err != nil || payload.PostID <= 0 {
		HandleError(w, http.StatusBadRequest, "Invalid post_id")
		return
	}
	payload.Content = r.FormValue("content")
	if payload.Content == "" || len(payload.Content) > 500 {
		HandleError(w, http.StatusBadRequest, "Content cannot be empty or too long")
		return
	}

	//Foreign Key Check: Query the database to verify the post_id actually exists
	doesPostExist, err := database.DoesPostExist(payload.PostID)
	if err != nil {
		HandleError(w, http.StatusInternalServerError, "Database error")
		return
	}
	if !doesPostExist {
		HandleError(w, http.StatusBadRequest, "Post does not exist")
		return
	}

	// Database Insert: Execute the insert query
	err = database.CreateComment(payload.PostID, payload.UserID, payload.Content)
	// failed to insert comment into the database
	if err != nil {
		HandleError(w, http.StatusInternalServerError, "Database error")
		return
	}

	// Success Response: Return a success message
	RespondJSON(w, http.StatusCreated, map[string]string{"message": "Comment created successfully"})
}
