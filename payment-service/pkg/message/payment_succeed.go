package message

import (
	"common/kafka/event"
	"encoding/json"
	"time"
)

type PaymentSucceededMessage struct {
	OrderID  int32           `json:"order_id"`
	Customer Customer        `json:"customer"`
	Items    []ItemPurchased `json:"items"`
}

type ItemPurchased struct {
	ProductID   int32   `json:"product_id"`
	ProductName string  `json:"product_name"`
	Quantity    int64   `json:"quantity"`
	Price       float64 `json:"price"`
}
type Customer struct {
	CustomerID    int32  `json:"customer_id"`
	CustomerEmail string `json:"customer_email"`
	CustomerName  string `json:"customer_name"`
}

func NewPaymentEnvelope(source, version string, payload PaymentSucceededMessage) (*event.Envelope, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	return &event.Envelope{
		EventID:    event.GenerateUniqueId("payment.succeed"),
		EventType:  "payment.succeed",
		OccurredAt: time.Now(),
		Source:     source,
		Version:    version,
		Payload:    data,
	}, nil
}
