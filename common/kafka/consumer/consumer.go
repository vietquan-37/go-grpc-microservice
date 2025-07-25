package consumer

import (
	"context"
	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"
)

type MessageHandler func(ctx context.Context, key, value []byte) error
type Consumer struct {
	reader  *kafka.Reader
	handler MessageHandler
}

func NewConsumer(brokers []string, topic string, groupId string, handler MessageHandler) *Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		Topic:    topic,
		GroupID:  groupId,
		MinBytes: 10e3,
		MaxBytes: 10e6,
	})
	return &Consumer{
		reader:  reader,
		handler: handler,
	}
}
func (c *Consumer) StartCommitFirst(ctx context.Context) error {
	log.Info().
		Str("topic", c.reader.Config().Topic).
		Str("group_id", c.reader.Config().Topic).
		Msg("Starting consumer")
	for {
		select {
		case <-ctx.Done():
			return c.reader.Close()
		default:
			m, err := c.reader.ReadMessage(ctx)
			if err != nil {
				log.Error().Err(err).Msg("Error reading message")
				continue
			}
			log.Debug().
				Str("topic", m.Topic).
				Int("partition", m.Partition).
				Int64("offset", m.Offset).
				Str("key", string(m.Key)).
				Msg("received message")
			if err := c.handler(ctx, m.Key, m.Value); err != nil {
				log.Error().Err(err).Msg("Error handling message")
			}

		}
	}
}
func (c *Consumer) Close() error {
	return c.reader.Close()
}

func (c *Consumer) Start(ctx context.Context) error {
	log.Info().
		Str("topic", c.reader.Config().Topic).
		Str("group_id", c.reader.Config().Topic).
		Msg("Starting consumer")
	for {
		select {
		case <-ctx.Done():
			return c.reader.Close()
		default:
			m, err := c.reader.FetchMessage(ctx)
			if err != nil {
				log.Error().Err(err).Msg("Error reading message")
				continue
			}
			log.Debug().
				Str("topic", m.Topic).
				Int("partition", m.Partition).
				Int64("offset", m.Offset).
				Str("key", string(m.Key)).
				Msg("received message")
			if err := c.handler(ctx, m.Key, m.Value); err != nil {
				log.Error().Err(err).Msg("Error handling message")
				continue
			}
			if err := c.reader.CommitMessages(ctx, m); err != nil {
				log.Error().Err(err).Msg("Error committing message")

			}

		}
	}
}
