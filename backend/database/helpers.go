package database

import "log"

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
