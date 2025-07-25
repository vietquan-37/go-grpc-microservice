package consumer

import (
	"common/cache"
	"common/kafka/producer"
	kafkaretry "common/kafka/retry"
	"context"
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"
	"github.com/vietquan-37/product-service/pkg/config"
	"github.com/vietquan-37/product-service/pkg/message"
	"github.com/vietquan-37/product-service/pkg/repository"
	"time"
)

type ProductConsumer struct {
	repo  repository.IProductRepo
	cfg   *config.Config
	p     *producer.Producer
	redis cache.Client
}

func NewProductConsumer(repo repository.IProductRepo, cfg *config.Config, p *producer.Producer, redis cache.Client) *ProductConsumer {
	return &ProductConsumer{
		repo:  repo,
		cfg:   cfg,
		p:     p,
		redis: redis,
	}
}

// TODO: Idempotent decrease
func (c *ProductConsumer) Process(ctx context.Context, msg kafka.Message) error {
	log.Debug().Msg("Processing payment succeeded message")
	ctx, cancel := context.WithTimeout(ctx, 4*time.Second)
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
	// TODO: check process yet by cache DONE
	eventKey := fmt.Sprintf("processed_event:%s", event.EventID)
	processed, err := c.redis.Get(ctx, eventKey)
	if err != nil && !errors.Is(err, cache.ErrorCacheMiss) {
		log.Error().Err(err).Str("event_key", eventKey).Msg("Failed to get event")
		return kafkaretry.NewRetryableError(err, "Cannot get event")
	}
	if processed != nil {
		log.Info().Type("msg", event).Msg("Event already processed")
		return nil
	}
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

			err = c.markEventAsProcessed(ctx, event.EventID, "fail")
			if err != nil {
				log.Error().Err(err).Msg("Failed to mark event as processed")
			}
			//pub to order
			payload, err := message.NewStockUpdateEnvelope(
				"product-service",
				"1",
				"update.failed",
				message.StockUpdateMessage{OrderID: event.Message.OrderID})
			if err != nil {
				log.Error().Err(err).Msg("Failed to create stock update message")
				return err
			}
			err = c.p.SendMessage(context.Background(), c.cfg.OrderTopic, nil, payload)
			if err != nil {
				log.Error().Err(err).Msg("Failed to send stock update message")
			}
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
	//mark as completed
	err = c.markEventAsProcessed(ctx, event.EventID, "success")
	if err != nil {
		log.Error().Err(err).Msg("Failed to mark event as processed")
	}
	//pub to order
	payload, err := message.NewStockUpdateEnvelope(
		"product-service",
		"1",
		"update.success",
		message.StockUpdateMessage{OrderID: event.Message.OrderID, Customer: event.Message.Customer, Items: items})
	if err != nil {
		log.Error().Err(err).Msg("Failed to create stock update message")
		return err
	}
	err = c.p.SendMessage(context.Background(), c.cfg.OrderTopic, nil, payload)
	if err != nil {
		log.Error().Err(err).Msg("Failed to send stock update message")
		return err
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
func (c *ProductConsumer) markEventAsProcessed(ctx context.Context, eventID, status string) error {
	eventKey := fmt.Sprintf("processed_event:%s", eventID)
	value := fmt.Sprintf("processed_at_%d_status_%s", time.Now().Unix(), status)
	ttl := 1 * time.Hour

	return c.redis.Set(ctx, eventKey, value, ttl)
}
