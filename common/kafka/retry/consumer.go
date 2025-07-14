package kafka_retry

import (
	"context"
	"github.com/rs/zerolog/log"

	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/segmentio/kafka-go"
)

type ProcessRetryHandler interface {
	Process(context.Context, kafka.Message) error
	MoveToDLQ(context.Context, kafka.Message)
}

type ConsumerWithRetry struct {
	reader      *kafka.Reader
	handler     ProcessRetryHandler
	maxRetries  int
	backoffFunc func() backoff.BackOff
	workerCount int
	retryQueue  chan kafka.Message
}

func NewConsumerWithRetry(reader *kafka.Reader, handler ProcessRetryHandler, maxRetries int, backoffFunc func() backoff.BackOff, workerCount int) *ConsumerWithRetry {
	if workerCount <= 0 {
		workerCount = 1
	}

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

func (c *ConsumerWithRetry) handleRetry(ctx context.Context, msg kafka.Message, workerID int) {
	retries := 0
	bo := c.backoffFunc()
	bo.Reset()

	for retries < c.maxRetries {
		if ctx.Err() != nil {
			log.Error().Msgf("Worker %d: context canceled during retry", workerID)
			return
		}

		retries++
		log.Info().Msgf("Worker %d: Retry %d for message %s", workerID, retries, string(msg.Key))

		if err := c.handler.Process(ctx, msg); err != nil {
			log.Error().Msgf("Worker %d: Processing error: %v", workerID, err)
			time.Sleep(bo.NextBackOff())
		} else {
			log.Info().Msgf("Worker %d: Successfully processed message %s", workerID, string(msg.Key))
			if err := c.reader.CommitMessages(ctx, msg); err != nil {
				log.Error().Msgf("Worker %d: Commit failed after retry: %v", workerID, err)
			}
			return
		}
	}

	log.Info().Msgf("Worker %d: Max retries exceeded for %s, moving to DLQ", workerID, string(msg.Key))
	c.handler.MoveToDLQ(ctx, msg)
}

// TODO: Phân biệt các lỗi có thể retry hoặc ko
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
				log.Printf("Processing failed for %s, adding to retry queue", string(msg.Key))
				select {
				case c.retryQueue <- msg:
				default:
					log.Error().Msg("Retry queue full, moving message to DLQ")
					c.handler.MoveToDLQ(ctx, msg)
				}
				continue
			}

			if err := c.reader.CommitMessages(ctx, msg); err != nil {
				log.Error().Msgf("Commit failed: %v", err)
			}
		}
	}
}
