package handler

import (
	commonclient "common/client"
	"fmt"

	"context"
	"encoding/json"
	"github.com/rs/zerolog/log"
	"github.com/vietquan-37/email-service/pkg/email"
	"github.com/vietquan-37/email-service/pkg/message"
)

type MessageHandler struct {
	emailService email.IEmailService
	authService  *commonclient.AuthClient
}

func NewMessageHandler(emailService email.IEmailService) *MessageHandler {
	return &MessageHandler{
		emailService: emailService,
	}
}
func (h *MessageHandler) ProcessMessage(ctx context.Context, key, value []byte) error {
	log.Debug().Msg("Processing message")
	event, err := message.ParseMessage(value)

	if err != nil {
		log.Error().Err(err).Msg("failed to parse message")
		return err
	}
	log.Debug().
		Str("key", string(key)).
		Str("event_type", event.EventType).
		Str("event_id", event.EventID).
		Str("source", event.Source).
		Str("version", event.Version).
		Time("OccurredAt", event.OccurredAt)
	switch event.EventType {
	case "user.created":
		return h.handleUserCreated(event.Payload)
	case "order.placed":
		return h.handleOrderPlaced(event.Payload)
	default:
		log.Warn().
			Str("event_type", event.EventType).
			Msg("Unknown event type, skipping")
		return nil

	}
}
func (h *MessageHandler) handleUserCreated(value []byte) error {
	var userMsg message.UserCreatePayload
	if err := json.Unmarshal(value, &userMsg); err != nil {
		return fmt.Errorf("failed to unmarshal user message: %w", err)
	}
	log.Info().
		Int32("user_id", userMsg.ID).
		Str("email", userMsg.Email).
		Str("full_name", userMsg.FullName).
		Msg("Processing user created event")

	return h.emailService.SendVerificationEmail(userMsg.Email, userMsg.FullName, userMsg.Token, email.TypeVerification)
}
func (h *MessageHandler) handleOrderPlaced(order []byte) error {
	var orderMsg message.OrderEmailMessage
	if err := json.Unmarshal(order, &orderMsg); err != nil {
		return fmt.Errorf("failed to unmarshal order message: %w", err)
	}
	customer := orderMsg.Customer
	log.Info().
		Str("email", customer.CustomerEmail).
		Str("full_name", customer.CustomerName).
		Int32("order_id", orderMsg.OrderID).
		Msg("Processing order placed event")
	orderData := email.OrderData{
		OrderID:     orderMsg.OrderID,
		TotalAmount: orderMsg.Amount,
		OrderDate:   orderMsg.OrderDate.Format("2006-01-02 15:04:05"),
		Status:      orderMsg.Status,
		Items:       convertToItems(orderMsg.Items),
	}
	return h.emailService.SendOrderConfirmationEmail(customer.CustomerEmail, customer.CustomerName, orderData, email.TypeOrderConfirmation)
}
func convertToItems(items []message.ItemPurchased) []email.OrderItem {
	var result []email.OrderItem
	for _, item := range items {
		result = append(result, email.OrderItem{
			Name:     item.ProductName,
			Quantity: item.Quantity,
			Price:    item.Price,
		})
	}
	return result
}
