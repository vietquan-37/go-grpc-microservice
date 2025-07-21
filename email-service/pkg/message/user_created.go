package message

import (
	"common/kafka/event"
	"encoding/json"
)

type UserCreatePayload struct {
	ID       int32  `json:"id"`
	Email    string `json:"email"`
	Token    string `json:"token"`
	FullName string `json:"full_name"`
}
type UserCreatedEvent struct {
	event.Envelope
	Message UserCreatePayload
}

func ParseMessage(data []byte) (*event.Envelope, error) {
	var envelope event.Envelope
	if err := json.Unmarshal(data, &envelope); err != nil {
		return nil, err
	}

	return &envelope, nil
}
