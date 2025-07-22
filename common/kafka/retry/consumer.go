package kafka_retry

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"

	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/segmentio/kafka-go"
)

type ProcessRetryHandler interface {
	Process(context.Context, kafka.Message) error
	MoveToDLQ(context.Context, kafka.Message, error)
}

type ConsumerWithRetry struct {
	reader      *kafka.Reader
	handler     ProcessRetryHandler
	maxRetries  int
	backoffFunc func() backoff.BackOff
	workerCount int
	retryQueue  chan kafka.Message
}

func NewConsumerWithRetry(brokers []string, topic string, groupId string, handler ProcessRetryHandler, maxRetries int, backoffFunc func() backoff.BackOff, workerCount int) *ConsumerWithRetry {
	if workerCount <= 0 {
		workerCount = 1
	}
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        brokers,
		Topic:          topic,
		CommitInterval: 0,
		GroupID:        groupId,
		MinBytes:       10e3,
		MaxBytes:       10e6,
	})
	return &ConsumerWithRetry{
		reader:      reader,
		handler:     handler,
		maxRetries:  maxRetries,
		backoffFunc: backoffFunc,
		workerCount: workerCount,
		retryQueue:  make(chan kafka.Message, 1000),
	}
}

func (c *ConsumerWithRetry) Start(ctx context.Context) error {
	c.startWorkers(ctx)
	return c.consumeLoop(ctx)
}

func (c *ConsumerWithRetry) startWorkers(ctx context.Context) {
	for i := 0; i < c.workerCount; i++ {
		go c.retryWorker(ctx, i)
	}
}

func (c *ConsumerWithRetry) retryWorker(ctx context.Context, workerID int) {
	for {
		select {
		case <-ctx.Done():
			log.Error().Msgf("Worker %d: context canceled, exiting", workerID)
			return
		case msg, ok := <-c.retryQueue:
			if !ok {
				log.Error().Msgf("Worker %d: retry queue closed", workerID)
				return
			}
			c.handleRetry(ctx, msg, workerID)
		}
	}
}

func (c *ConsumerWithRetry) handleRetry(ctx context.Context, msg kafka.Message, id int) {
	retries := 0
	bo := c.backoffFunc()
	bo.Reset()

	for retries < c.maxRetries {
		if ctx.Err() != nil {
			return
		}
		retries++
		log.Info().Msgf("Worker %d retry attempt %d for key=%s", id, retries, string(msg.Key))

		err := c.handler.Process(ctx, msg)
		if err == nil {

			if commitErr := c.reader.CommitMessages(ctx, msg); commitErr != nil {
				log.Error().Err(commitErr).Msg("Failed commit after retry success")
			}
			return
		}

		if !ShouldRetry(err) {
			log.Info().Msgf("Non-retryable error after retries for key=%s: %v", string(msg.Key), err)
			c.handler.MoveToDLQ(ctx, msg, err)
			return
		}

		log.Error().Err(err).Msg("Retryable error, backing off before next attempt")
		time.Sleep(bo.NextBackOff())
	}

	log.Info().Msgf("Max retries exceeded for key=%s, moving to DLQ", string(msg.Key))
	c.handler.MoveToDLQ(ctx, msg, fmt.Errorf("max retries exceeded"))
}

func (c *ConsumerWithRetry) consumeLoop(ctx context.Context) error {
	defer close(c.retryQueue)
	for {
		select {
		case <-ctx.Done():
			return c.reader.Close()
		default:
			msg, err := c.reader.FetchMessage(ctx)
			if err != nil {
				log.Error().Msgf("Kafka fetch error: %v", err)
				continue
			}

			if err := c.handler.Process(ctx, msg); err != nil {
				if ShouldRetry(err) {
					log.Printf("Processing failed for %s, adding to retry queue", string(msg.Key))
					select {
					case c.retryQueue <- msg:
					default:
						log.Error().Msg("Retry queue full, moving message to DLQ")
						c.handler.MoveToDLQ(ctx, msg, err)
					}
				} else {
					c.handler.MoveToDLQ(ctx, msg, err)
					if commitErr := c.reader.CommitMessages(ctx, msg); commitErr != nil {
						log.Error().Err(commitErr).Msg("Failed commit after DLQ move")
					}
				}
				continue
			}

			if err := c.reader.CommitMessages(ctx, msg); err != nil {
				log.Error().Msgf("Commit failed: %v", err)
			}
		}
	}
}

type DLQMessage struct {
	OriginalMessage json.RawMessage `json:"original_message"`
	Error           string          `json:"error"`
	FailedAt        time.Time       `json:"failed_at"`
}
