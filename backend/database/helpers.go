package database

import (
	"log"
	"zone/backend/types"
)

// UpdateLastSeen updates the last_seen timestamp for the given user.
// Shared by handlers (login) and middleware (auth).
func UpdateLastSeen(userID int) {
	_, err := Database.Exec(
		"UPDATE users SET last_seen = CURRENT_TIMESTAMP WHERE id = ?", userID,
	)
	if err != nil {
		log.Println("UpdateLastSeen:", err)
	}
}

// DoesPostExist checks if a post with the given ID exists in the database.
func DoesPostExist(postID int) (bool, error) {
	var exists bool
	err := Database.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM posts WHERE id = ?)", postID,
	).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

// CreateComment inserts a new comment into the database
func CreateComment(postID, userID int, content string) error {
	_, err := Database.Exec(
		"INSERT INTO comments (post_id, user_id, content) VALUES (?, ?, ?)",
		postID, userID, content,
	)
	if err != nil {
		return err
	}
	return nil
}

// GetComments retrieves comments for a specific post with pagination along with the commenter's nickname and the creation timestamp.
func GetComments(postID, limit, offset int) ([]types.Comment, error) {
	// Query the database
	rows, err := Database.Query(
		`SELECT u.nickname, c.content, c.created_at
		FROM comments c
		JOIN users u ON c.user_id = u.id
		WHERE c.post_id = ?
		ORDER BY c.created_at DESC
		LIMIT ? OFFSET ?`, postID, limit, offset,
	)
	// Handle any errors that occur during the query
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Iterate through the result set and populate the comments slice
	var comments []types.Comment
	for rows.Next() {
		var comment types.Comment
		if err := rows.Scan(&comment.Nickname, &comment.Content, &comment.CreatedAt); err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	// Check for any errors that occurred during iteration
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}
