package database

// Database helper additions for presence and messaging.
// Includes GetUsersByLastSeen which returns users ordered by recent activity (last_seen).

import (
	"database/sql"
	"log"
	"zone/backend/types"
)

// UpdateLastSeen updates the last_seen timestamp for the given user.
// Shared by handlers (login) and middleware (auth).
func UpdateLastSeen(userID int) {
	_, err := Database.Exec(
		"UPDATE users SET last_seen = CURRENT_TIMESTAMP WHERE id = ?",
		userID,
	)
	if err != nil {
		log.Println("UpdateLastSeen:", err)
	}
}


func GetUsersByLastSeen(currentUserID int) ([]types.UserStatus, error) {
	rows, err := Database.Query(
		`SELECT
			u.id,
			u.nickname,
			u.last_seen,
			COALESCE((
				SELECT m.content FROM messages m
				WHERE (m.sender_id = ? AND m.receiver_id = u.id)
				   OR (m.sender_id = u.id AND m.receiver_id = ?)
				ORDER BY m.created_at DESC
				LIMIT 1
			), '') AS last_message,
			COALESCE((
				SELECT m.created_at FROM messages m
				WHERE (m.sender_id = ? AND m.receiver_id = u.id)
				   OR (m.sender_id = u.id AND m.receiver_id = ?)
				ORDER BY m.created_at DESC
				LIMIT 1
			), '') AS last_message_time
		FROM users u
		WHERE u.id != ?
		ORDER BY
			CASE WHEN last_message_time = '' THEN 1 ELSE 0 END,
			last_message_time DESC,
			u.nickname ASC`,
		currentUserID, currentUserID, currentUserID, currentUserID, currentUserID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []types.UserStatus
	for rows.Next() {
		var user types.UserStatus
		if err := rows.Scan(
			&user.UserID,
			&user.Nickname,
			&user.LastSeen,
			&user.LastMessage,
			&user.LastMessageTime,
		); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

// InsertMessage stores a private message in the database and updates user activity.
func InsertMessage(senderID, receiverID int, content string) (int, error) {
	result, err := Database.Exec(
		"INSERT INTO messages (sender_id, receiver_id, content) VALUES (?, ?, ?)",
		senderID, receiverID, content,
	)
	if err != nil {
		return 0, err
	}

	messageID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	// Update the sender's last seen timestamp.
	UpdateLastSeen(senderID)

	return int(messageID), nil
}

// GetMessageByID returns one private message with sender nickname.
func GetMessageByID(messageID int) (types.Message, error) {
	row := Database.QueryRow(
		`SELECT m.id, m.sender_id, m.receiver_id, m.content, m.created_at, u.nickname
		FROM messages m
		JOIN users u ON u.id = m.sender_id
		WHERE m.id = ?`,
		messageID,
	)

	var message types.Message
	if err := row.Scan(&message.ID, &message.SenderID, &message.ReceiverID, &message.Content, &message.CreatedAt, &message.SenderNickname); err != nil {
		return types.Message{}, err
	}

	return message, nil
}

// GetMessages retrieves private messages between two users with pagination.
// beforeID specifies the oldest message ID already loaded; pass 0 to get the latest messages.
func GetMessages(userID1, userID2, limit, beforeID int) ([]types.Message, bool, error) {
	var query string
	var args []interface{}

	if beforeID <= 0 {
		query = `SELECT m.id, m.sender_id, m.receiver_id, m.content, m.created_at, u.nickname
		FROM messages m
		JOIN users u ON u.id = m.sender_id
		WHERE (m.sender_id = ? AND m.receiver_id = ?) OR (m.sender_id = ? AND m.receiver_id = ?)
		ORDER BY m.created_at DESC
		LIMIT ?`
		args = []interface{}{userID1, userID2, userID2, userID1, limit + 1}
	} else {
		query = `SELECT m.id, m.sender_id, m.receiver_id, m.content, m.created_at, u.nickname
		FROM messages m
		JOIN users u ON u.id = m.sender_id
		WHERE ((m.sender_id = ? AND m.receiver_id = ?) OR (m.sender_id = ? AND m.receiver_id = ?))
		AND m.id < ?
		ORDER BY m.created_at DESC
		LIMIT ?`
		args = []interface{}{userID1, userID2, userID2, userID1, beforeID, limit + 1}
	}

	rows, err := Database.Query(query, args...)
	if err != nil {
		return nil, false, err
	}
	defer rows.Close()

	var messages []types.Message
	for rows.Next() {
		var message types.Message
		if err := rows.Scan(&message.ID, &message.SenderID, &message.ReceiverID, &message.Content, &message.CreatedAt, &message.SenderNickname); err != nil {
			return nil, false, err
		}
		messages = append(messages, message)
	}

	if err := rows.Err(); err != nil {
		return nil, false, err
	}

	hasMore := len(messages) > limit
	if hasMore {
		messages = messages[:limit]
	}

	return messages, hasMore, nil
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
func GetComments(postID, userID, limit, offset int) ([]types.Comment, error) {
	rows, err := Database.Query(
		`SELECT c.id, u.nickname, c.content, c.created_at,
			COALESCE(cr_counts.like_count, 0),
			COALESCE(cr_counts.dislike_count, 0),
			COALESCE(cr_user.is_like, -1)
		FROM comments c
		JOIN users u ON c.user_id = u.id
		LEFT JOIN (
			SELECT comment_id,
				SUM(CASE WHEN is_like = 1 THEN 1 ELSE 0 END) AS like_count,
				SUM(CASE WHEN is_like = 0 THEN 1 ELSE 0 END) AS dislike_count
			FROM comment_reactions
			GROUP BY comment_id
		) cr_counts ON cr_counts.comment_id = c.id
		LEFT JOIN comment_reactions cr_user ON cr_user.comment_id = c.id AND cr_user.user_id = ?
		WHERE c.post_id = ?
		ORDER BY c.created_at DESC
		LIMIT ? OFFSET ?`, userID, postID, limit, offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []types.Comment
	for rows.Next() {
		var comment types.Comment
		var userReaction int
		if err := rows.Scan(&comment.ID, &comment.Nickname, &comment.Content, &comment.CreatedAt, &comment.LikeCount, &comment.DislikeCount, &userReaction); err != nil {
			return nil, err
		}
		if userReaction == 1 {
			comment.UserReaction = "like"
		} else if userReaction == 0 {
			comment.UserReaction = "dislike"
		} else {
			comment.UserReaction = ""
		}
		comments = append(comments, comment)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}

// DoesCommentExist checks if a comment with the given ID exists.
func DoesCommentExist(commentID int) (bool, error) {
	var exists bool
	err := Database.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM comments WHERE id = ?)", commentID,
	).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

// GetCommentReactionSummary returns like/dislike counts and the user's reaction for a comment.
func GetCommentReactionSummary(commentID, userID int) (types.CommentReactResponse, error) {
	var resp types.CommentReactResponse

	err := Database.QueryRow(`
		SELECT
			COALESCE(SUM(CASE WHEN is_like = 1 THEN 1 ELSE 0 END), 0),
			COALESCE(SUM(CASE WHEN is_like = 0 THEN 1 ELSE 0 END), 0)
		FROM comment_reactions WHERE comment_id = ?`,
		commentID,
	).Scan(&resp.LikeCount, &resp.DislikeCount)
	if err != nil {
		return resp, err
	}

	var userIsLike int
	err = Database.QueryRow(
		"SELECT is_like FROM comment_reactions WHERE comment_id = ? AND user_id = ?",
		commentID, userID,
	).Scan(&userIsLike)

	if err == nil {
		if userIsLike == 1 {
			resp.UserReaction = "like"
		} else {
			resp.UserReaction = "dislike"
		}
	} else if err != sql.ErrNoRows {
		return resp, err
	}

	return resp, nil
}