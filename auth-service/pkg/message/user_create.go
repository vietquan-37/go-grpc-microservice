package message

import (
	"common/kafka/event"
	"encoding/json"
	"time"
)

type UserCreateMessage struct {
	ID       int32  `json:"id"`
	Email    string `json:"email"`
	Token    string `json:"token"`
	FullName string `json:"full_name"`
}

func NewUserCreatedEnvelope(source, version string, payload UserCreateMessage) (*event.Envelope, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	return &event.Envelope{
		EventID:    event.GenerateUniqueId("user.created"),
		EventType:  "user.created",
		OccurredAt: time.Now(),
		Source:     source,
		Version:    version,
		Payload:    data,
	}, nil
}
