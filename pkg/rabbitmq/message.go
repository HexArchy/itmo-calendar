package rabbitmq

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// Message is a standard message for RabbitMQ transport.
type Message struct {
	MessageID string                 `json:"message_id"`
	CreatedAt time.Time              `json:"created_at"`
	Headers   map[string]interface{} `json:"headers,omitempty"`
	Body      json.RawMessage        `json:"body"`
}

// NewMessage creates a new Message with autogenerated headers.
func NewMessage(body interface{}, headers map[string]interface{}) (*Message, error) {
	raw, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	msg := &Message{
		MessageID: uuid.NewString(),
		CreatedAt: time.Now().UTC(),
		Headers:   headers,
		Body:      raw,
	}

	return msg, nil
}
