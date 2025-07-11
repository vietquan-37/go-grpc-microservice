package producer

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"
	"time"
)

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(broker ...string) *Producer {
	writer := kafka.Writer{
		Addr:            kafka.TCP(broker...),
		RequiredAcks:    kafka.RequireOne,
		Async:           false,
		MaxAttempts:     5,
		WriteBackoffMin: 300 * time.Millisecond,
		WriteBackoffMax: 1 * time.Second,
		BatchSize:       1,
		BatchTimeout:    1 * time.Nanosecond,
	}
	return &Producer{
		writer: &writer,
	}

}
func NewProducerSafe(brokers []string) *Producer {
	writer := kafka.Writer{
		Addr:            kafka.TCP(brokers...),
		RequiredAcks:    kafka.RequireAll,
		Async:           false,
		Balancer:        &kafka.Hash{},
		WriteBackoffMin: 300 * time.Millisecond,
		WriteBackoffMax: 2 * time.Second,
		BatchSize:       1,
		BatchTimeout:    1 * time.Nanosecond,
	}
	return &Producer{
		writer: &writer,
	}
}
func (p *Producer) SendMessage(ctx context.Context, topic string, key []byte, value interface{}) error {
	valueBytes, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %v", err)
	}
	msg := kafka.Message{
		Topic: topic,
		Key:   key,
		Value: valueBytes,
		Time:  time.Now(),
	}
	err = p.writer.WriteMessages(ctx, msg)
	if err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}
	log.Info().Str("topic", topic).Msg("message sent")
	return nil
}
func (p *Producer) Close() error {
	return p.writer.Close()
}
