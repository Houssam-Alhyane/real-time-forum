package handlers

import "testing"

func TestParseChatMessagePayload(t *testing.T) {
	payload := []byte(`{"type":"send_message","receiver_id":2,"content":"hello"}`)

	req, err := parseChatMessagePayload(payload)
	if err != nil {
		t.Fatalf("parseChatMessagePayload returned error: %v", err)
	}

	if req.Type != "send_message" {
		t.Fatalf("expected type send_message, got %q", req.Type)
	}
	if req.ReceiverID != 2 {
		t.Fatalf("expected receiver_id 2, got %d", req.ReceiverID)
	}
	if req.Content != "hello" {
		t.Fatalf("expected content hello, got %q", req.Content)
	}
}
