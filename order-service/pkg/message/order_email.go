package message

import (
	"common/kafka/event"
	"encoding/json"
	"time"
)

type OrderEmailMessage struct {
	OrderID   int32           `json:"order_id"`
	Amount    float64         `json:"amount"`
	OrderDate time.Time       `json:"order_date"`
	Status    string          `json:"status"`
	Customer  Customer        `json:"customer"`
	Items     []ItemPurchased `json:"items"`
}
type Customer struct {
	CustomerID    int32  `json:"customer_id"`
	CustomerEmail string `json:"customer_email"`
	CustomerName  string `json:"customer_name"`
}

func NewOrderEmailEnvelope(source, version string, eventType string, payload OrderEmailMessage) (*event.Envelope, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	return &event.Envelope{
		EventID:    event.GenerateUniqueId(eventType),
		EventType:  eventType,
		OccurredAt: time.Now(),
		Source:     source,
		Version:    version,
		Payload:    data,
	}, nil
}
