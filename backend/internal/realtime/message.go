package realtime

import "encoding/json"

// WebSocketMessage is the generic structure for communication.
type WebSocketMessage struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload,omitempty"`
}

// ErrorPayload is a standard structure for sending errors to the client.
type ErrorPayload struct {
	Message string `json:"message"`
}

func createWebSocketMessage(msgType string, payload any) ([]byte, error) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	msg := WebSocketMessage{Type: msgType, Payload: payloadBytes}
	return json.Marshal(msg)
}
