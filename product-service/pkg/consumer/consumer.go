package consumer

import (
	"common/kafka/producer"
	kafkaretry "common/kafka/retry"
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"
	"github.com/vietquan-37/product-service/pkg/config"
	"github.com/vietquan-37/product-service/pkg/message"
	"github.com/vietquan-37/product-service/pkg/repository"
	"time"
)

type ProductConsumer struct {
	repo repository.IProductRepo
	cfg  *config.Config
	p    *producer.Producer
}

func NewProductConsumer(repo repository.IProductRepo, cfg *config.Config, p *producer.Producer) *ProductConsumer {
	return &ProductConsumer{
		repo: repo,
		cfg:  cfg,
		p:    p,
	}
}

// TODO: Idempotent decrease
func (c *ProductConsumer) Process(ctx context.Context, msg kafka.Message) error {
	log.Debug().Msg("Processing payment succeeded message")
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	event, err := message.ParsePaymentSucceededMessage(msg.Value)
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse payment succeeded message")
		return kafkaretry.NewNonRetryableError(err, "Cannot parse message")
	}

	log.Debug().
		Str("event_type", event.EventType).
		Str("event_id", event.EventID).
		Str("source", event.Source).
		Str("version", event.Version).
		Time("OccurredAt", event.OccurredAt).
		Msg("Received event")

	items := event.Message.Items
	if len(items) == 0 {
		log.Warn().Msg("No items in payment message")
		return kafkaretry.NewNonRetryableError(
			fmt.Errorf("empty items"),
			"No items in message")
	}

	ids := make([]int32, 0, len(items))
	itemMap := make(map[uint]int64)
	for _, item := range items {
		ids = append(ids, item.ProductID)
		itemMap[uint(item.ProductID)] = item.Quantity
	}

	products, err := c.repo.FindProductsByIds(ctx, ids)
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch products from DB")
		return err
	}

	for _, product := range products {
		reqQty := itemMap[product.ID]
		if product.Stock < reqQty {
			log.Error().
				Uint("product_id", product.ID).
				Int64("stock", product.Stock).
				Int64("requested", reqQty).
				Msg("Not enough stock for product")

			return kafkaretry.NewNonRetryableError(
				fmt.Errorf("not enough stock for product %d", product.ID),
				"Insufficient stock")
		}
	}

	for _, product := range products {
		reqQty := itemMap[product.ID]
		product.Stock -= reqQty
		_, err := c.repo.UpdateProduct(ctx, product)
		if err != nil {
			log.Error().
				Err(err).
				Uint("product_id", product.ID).
				Msg("Failed to decrease product stock")
			return err
		}
	}

	log.Info().Str("event_id", event.EventID).Msg("Successfully processed payment event and updated stocks")
	return nil
}

func (c *ProductConsumer) MoveToDLQ(ctx context.Context, msg kafka.Message, reason error) {
	dlqMsg := kafkaretry.DLQMessage{
		OriginalMessage: msg.Value,
		Error:           reason.Error(),
		FailedAt:        time.Now(),
	}
	err := c.p.SendMessage(ctx, c.cfg.DLQTOPIC, msg.Key, dlqMsg)
	if err != nil {
		log.Error().Err(err).
			Str("dlq_topic", c.cfg.DLQTOPIC).
			Msg("Failed to publish message to DLQ")
	} else {
		log.Info().Str("dlq_topic", c.cfg.DLQTOPIC).
			Msg("Published message to DLQ")
	}
}
