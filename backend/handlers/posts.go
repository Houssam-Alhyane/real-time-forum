package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"zone/backend/database"
)

// PostResponse represents the structure of data sent to the frontend
type PostResponse struct {
	ID           int    `json:"id"`
	Title        string `json:"title"`
	Content      string `json:"content"`
	UserID       int    `json:"user_id"`
	CategoryName string `json:"category_name"`
}

// CategoryResponse represents the basic category structure for frontend forms
type CategoryResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Helper to send JSON error responses to matching SPA expectations
func sendJSONError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

// GetPostsAPI fetches all posts combined with their category names from the database
func GetPostsAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		sendJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	query := `
		SELECT p.id, p.title, p.content, p.user_id, c.name 
		FROM posts p 
		JOIN categories c ON p.category_id = c.id
		ORDER BY p.id DESC`

	rows, err := database.Database.Query(query)
	if err != nil {
		log.Printf("Error selecting posts: %v", err)
		sendJSONError(w, http.StatusInternalServerError, "Database error loading posts")
		return
	}
	defer rows.Close()

	var posts []PostResponse = []PostResponse{} // Initialize as empty slice instead of nil
	for rows.Next() {
		var p PostResponse
		if err := rows.Scan(&p.ID, &p.Title, &p.Content, &p.UserID, &p.CategoryName); err != nil {
			log.Printf("Error scanning posts: %v", err)
			continue
		}
		posts = append(posts, p)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

// GetCategoriesAPI returns the list of available categories for the dropdown menu
func GetCategoriesAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		sendJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	rows, err := database.Database.Query("SELECT id, name FROM categories ORDER BY id ASC")
	if err != nil {
		log.Printf("Error selecting categories: %v", err)
		sendJSONError(w, http.StatusInternalServerError, "Database error loading categories")
		return
	}
	defer rows.Close()

	var categories []CategoryResponse = []CategoryResponse{}
	for rows.Next() {
		var c CategoryResponse
		if err := rows.Scan(&c.ID, &c.Name); err != nil {
			continue
		}
		categories = append(categories, c)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}
// CreatePostAPI handles incoming POST requests to submit a new user topic
func CreatePostAPI(w http.ResponseWriter, r *http.Request) {
  if r.Method != http.MethodPost {
    sendJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
    return
  }

  cookie, err := r.Cookie("session_token")
  if err != nil {
    sendJSONError(w, http.StatusUnauthorized, "You must be logged in to post")
    return
  }

  // Retrieve user_id from active session
  var userID int
  query := "SELECT user_id FROM sessions WHERE id = ? AND expires_at > DATETIME('now')"
  err = database.Database.QueryRow(query, cookie.Value).Scan(&userID)
  
  if err != nil {
    if err == sql.ErrNoRows {
      sendJSONError(w, http.StatusUnauthorized, "Invalid or expired session. Please re-login.")
    } else {
      log.Printf("Session query database crash: %v", err)
      sendJSONError(w, http.StatusInternalServerError, "Internal auth database validation failure")
    }
    return
  }

  title := strings.TrimSpace(r.FormValue("title"))
  content := strings.TrimSpace(r.FormValue("content"))
  categoryIDStr := r.FormValue("category_id")

  if title == "" || content == "" || categoryIDStr == "" {
    sendJSONError(w, http.StatusBadRequest, "All fields are required")
    return
  }

  categoryID, err := strconv.Atoi(categoryIDStr)
  if err != nil {
    sendJSONError(w, http.StatusBadRequest, "Invalid category selection")
    return
  }

  _, err = database.Database.Exec(
    "INSERT INTO posts (title, content, user_id, category_id) VALUES (?, ?, ?, ?)",
    title, content, userID, categoryID,
  )
  if err != nil {
    log.Printf("Error inserting post: %v", err)
    sendJSONError(w, http.StatusInternalServerError, "Failed to create post inside database")
    return
  }

  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(http.StatusCreated)
  json.NewEncoder(w).Encode(map[string]string{"message": "Post created successfully!"})
}