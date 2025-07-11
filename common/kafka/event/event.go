package event

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"time"
)

type Envelope struct {
	EventID    string          `json:"event_id"` // UUID (idempotency)
	EventType  string          `json:"event_type"`
	OccurredAt time.Time       `json:"occurred_at"` // timestamp
	Source     string          `json:"source"`      // service name
	Version    string          `json:"version"`     // version of message schema
	Payload    json.RawMessage `json:"payload"`     // actual data (generic)
}

func GenerateUniqueId(eventType string) string {
	return fmt.Sprintf("%s_%s", eventType, uuid.NewString())
}
