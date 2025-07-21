package consumer

import (
	"common/kafka/producer"
	kafkaretry "common/kafka/retry"
	"context"
	"errors"
	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"
	"github.com/vietquan-37/order-service/pkg/config"
	"github.com/vietquan-37/order-service/pkg/enum"
	"github.com/vietquan-37/order-service/pkg/message"
	"github.com/vietquan-37/order-service/pkg/repository"
	"gorm.io/gorm"
	"time"
)

type OrderConsumer struct {
	repo repository.IOrderRepo
	cfg  *config.Config
	p    *producer.Producer
}

func NewOrderConsumer(repo repository.IOrderRepo, cfg *config.Config, p *producer.Producer) *OrderConsumer {
	return &OrderConsumer{
		repo: repo,
		cfg:  cfg,
		p:    p,
	}
}

func (c *OrderConsumer) Process(ctx context.Context, msg kafka.Message) error {
	log.Debug().Msg("Processing message")
	ctx, cancel := context.WithTimeout(ctx, time.Second*4)
	defer cancel()
	event, err := message.ParseStockUpdateMessage(msg.Value)
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse payment succeeded message")
		return kafkaretry.NewNonRetryableError(err, "Cannot parse message")
	}
	order, err := c.repo.GetOrderById(ctx, event.Message.OrderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return kafkaretry.NewNonRetryableError(err, "Order not found")
		}
		return err
	}
	log.Debug().
		Str("event_type", event.EventType).
		Str("event_id", event.EventID).
		Str("source", event.Source).
		Str("version", event.Version).
		Time("OccurredAt", event.OccurredAt)
	switch event.EventType {
	case "update.failed":
		order.Status = enum.CANCELLED

		return c.repo.UpdateOrder(ctx, order)

	case "update.success":
		order.Status = enum.COMPLETED
		order.OrderDate = time.Now()
		err = c.repo.UpdateOrder(ctx, order)
		if err != nil {
			log.Error().Err(err).Msg("Failed to update order")
			return err
		}
		payload, err := message.NewOrderEmailEnvelope("order-service", "1", "order.placed", message.OrderEmailMessage{
			OrderID:   int32(order.ID),
			Amount:    order.Amount,
			Status:    string(order.Status),
			Customer:  event.Message.Customer,
			OrderDate: order.OrderDate,
			Items:     event.Message.Items,
		})
		if err != nil {
			log.Error().Err(err).Msg("Failed to encode order email payload")
			return kafkaretry.NewNonRetryableError(err, "Failed to encode order email payload")
		}
		//send mail
		return c.p.SendMessage(context.Background(), c.cfg.EmailTopic, nil, payload)

	default:
		log.Warn().
			Str("event_type", event.EventType).
			Msg("Unknown event type, skipping")
		return nil
	}

}
func (c *OrderConsumer) MoveToDLQ(ctx context.Context, msg kafka.Message, reason error) {
	dlqMsg := kafkaretry.DLQMessage{
		OriginalMessage: msg.Value,
		Error:           reason.Error(),
		FailedAt:        time.Now(),
	}
	err := c.p.SendMessage(ctx, c.cfg.DLQTopic, msg.Key, dlqMsg)
	if err != nil {
		log.Error().Err(err).
			Str("dlq_topic", c.cfg.DLQTopic).
			Msg("Failed to publish message to DLQ")
	} else {
		log.Info().Str("dlq_topic", c.cfg.DLQTopic).
			Msg("Published message to DLQ")
	}
}
