package types

type CommentPayload struct {
	PostID  int    `json:"post_id"`
	Content string `json:"content"`
	UserID  int    `json:"user_id"`
}

type CommentsResponse struct {
	Comments []Comment `json:"comments"`
}

type Comment struct {
	ID           int    `json:"id"`
	Nickname     string `json:"nickname"`
	Content      string `json:"content"`
	CreatedAt    string `json:"created_at"`
	LikeCount    int    `json:"like_count"`
	DislikeCount int    `json:"dislike_count"`
	UserReaction string `json:"user_reaction"`
}

type CommentReactRequest struct {
	CommentID int    `json:"comment_id"`
	Type      string `json:"type"`
}

type CommentReactResponse struct {
	LikeCount    int    `json:"like_count"`
	DislikeCount int    `json:"dislike_count"`
	UserReaction string `json:"user_reaction"`
}

// MessagePayload represents the data a client sends to send a private message.
type MessagePayload struct {
	SenderID   int    `json:"sender_id"`
	ReceiverID int    `json:"receiver_id"`
	Content    string `json:"content"`
}

// WebSocketPayload represents the envelope used for real-time socket events.
type WebSocketPayload struct {
	Type     string   `json:"type"`
	Message  *Message `json:"message,omitempty"`
	UserID   int      `json:"user_id,omitempty"`
	Nickname string   `json:"nickname,omitempty"`
}

// Message represents the actual private message data stored and exchanged.
type Message struct {
	ID             int    `json:"id"`
	SenderID       int    `json:"sender_id"`
	ReceiverID     int    `json:"receiver_id"`
	Content        string `json:"content"`
	CreatedAt      string `json:"created_at"`
	SenderNickname string `json:"sender_nickname"`
}

// In backend/types/types.go, update UserStatus to include the two new
// nullable fields returned by GetUsersByLastSeen. Do not replace your whole
// types.go — just add LastMessage / LastMessageTime to the existing struct.

type UserStatus struct {
	UserID          int
	Nickname        string
	LastSeen        string
	LastMessage     string // "" if no message history with this user
	LastMessageTime string // "" if no message history with this user
}
	
// Make sure "database/sql" is imported in types.go if it isn't already.
// WebSocketPayload.Type values used by the server:
// - "user_online": a user came online (first tab connected)
// - "user_offline": a user went fully offline (last tab disconnected)
