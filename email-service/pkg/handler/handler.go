package handler

import (
	"context"
	"github.com/rs/zerolog/log"
	"github.com/vietquan-37/email-service/pkg/email"
	"github.com/vietquan-37/email-service/pkg/message"
)

type MessageHandler struct {
	emailService email.IEmailService
}

func NewMessageHandler(emailService email.IEmailService) *MessageHandler {
	return &MessageHandler{
		emailService: emailService,
	}
}
func (h *MessageHandler) MessageHandler(ctx context.Context, key, value []byte) error {
	log.Debug().Msg("Processing message")
	event, err := message.ParseUserCreatedMessage(value)
	if err != nil {
		log.Error().Err(err).Msg("failed to parse message")
	}
	log.Debug().
		Str("event_type", event.EventType).
		Str("event_id", event.EventID).
		Str("source", event.Source).
		Str("version", event.Version).
		Time("OccurredAt", event.OccurredAt)
	switch event.EventType {
	case "user.created":
		return h.handleUserCreated(event.Message)
	default:
		log.Warn().
			Str("event_type", event.EventType).
			Msg("Unknown event type, skipping")
		return nil

	}
}
func (h *MessageHandler) handleUserCreated(userMsg message.UserCreatePayload) error {

	log.Info().
		Int32("user_id", userMsg.ID).
		Str("email", userMsg.Email).
		Str("full_name", userMsg.FullName).
		Msg("Processing user created event")

	return h.emailService.SendVerificationEmail(userMsg.Email, userMsg.FullName, userMsg.Token)
}
