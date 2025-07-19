package consumer

import (
	"common/kafka/producer"
	kafka_retry "common/kafka/retry"
	"github.com/vietquan-37/order-service/pkg/repository"
)

type OrderConsumer struct {
	repo repository.IOrderRepo
	p    *producer.Producer
	c    *kafka_retry.ConsumerWithRetry
}

func NewOrderConsumer(repo repository.IOrderRepo, p *producer.Producer, c *kafka_retry.ConsumerWithRetry) *OrderConsumer {
	return &OrderConsumer{
		repo: repo,
		p:    p,
		c:    c,
	}
}

//func (c *OrderConsumer) Process(ctx context.Context, msg kafka.Message) error {
//	//log.Debug().Msg("Processing message")
//	//event, err := message.ParseUserCreatedMessage(value)
//	//if err != nil {
//	//	log.Error().Err(err).Msg("failed to parse message")
//	//}
//	//log.Debug().
//	//	Str("event_type", event.EventType).
//	//	Str("event_id", event.EventID).
//	//	Str("source", event.Source).
//	//	Str("version", event.Version).
//	//	Time("OccurredAt", event.OccurredAt)
//
//}
//func (c *OrderConsumer) MoveToDLQ(ctx context.Context, msg kafka.Message) {
//
//}
