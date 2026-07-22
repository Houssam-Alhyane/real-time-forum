package handlers


import (
	"encoding/json"
	"errors"
	"log"
)


type typingRequest struct {
	Type       string `json:"type"`
	ReceiverID int    `json:"receiver_id"`
}


type typingBroadcast struct {
	Type           string `json:"type"` // always "typing"
	SenderID       int    `json:"sender_id"`
	SenderNickname string `json:"sender_nickname"`
	Status         string `json:"status"` // "start" or "stop"
}

func handleTypingEvent(userID int, client *wsClient, envelopeType string, payload []byte) error {
	var req typingRequest
	if err := json.Unmarshal(payload, &req); err != nil {
		return sendSocketError(client, "invalid typing payload")
	}

	if err := validateTypingRequest(req); err != nil {
		return sendSocketError(client, err.Error())
	}

	status := "start"
	if envelopeType == "typing_stop" {
		status = "stop"
	}

	nickname, err := getNicknameByUserID(userID)
	if err != nil {
		log.Println("getNicknameByUserID (typing):", err)
	}

	broadcast := typingBroadcast{
		Type:           "typing",
		SenderID:       userID,
		SenderNickname: nickname,
		Status:         status,
	}

	broadcastTyping(broadcast, req.ReceiverID)

	return nil
}

func validateTypingRequest(req typingRequest) error {
	if req.ReceiverID <= 0 {
		return errors.New("invalid receiver_id")
	}
	return nil
}

func broadcastTyping(payload typingBroadcast, receiverID int) {
	clientsMu.Lock()
	conns := append([]*wsClient(nil), clients[receiverID]...)
	clientsMu.Unlock()

	for _, c := range conns {
		if err := c.sendJSON(payload); err != nil {
			log.Println("broadcastTyping:", err)
		}
	}
}