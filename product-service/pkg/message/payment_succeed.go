package message

import (
	"common/kafka/event"
	"encoding/json"
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

type PaymentSucceededEvent struct {
	Message PaymentSucceededMessage
	event.Envelope
}

func ParsePaymentSucceededMessage(data []byte) (*PaymentSucceededEvent, error) {
	var envelope event.Envelope
	if err := json.Unmarshal(data, &envelope); err != nil {
		return nil, err
	}
	var paymentMsg PaymentSucceededMessage
	if err := json.Unmarshal(envelope.Payload, &paymentMsg); err != nil {
		return nil, err
	}

	return &PaymentSucceededEvent{
		Envelope: envelope,
		Message:  paymentMsg,
	}, nil
}
