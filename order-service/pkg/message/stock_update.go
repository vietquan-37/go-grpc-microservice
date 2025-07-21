package message

import (
	"common/kafka/event"
	"encoding/json"
)

type StockUpdateMessage struct {
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
type StockUpdateEvent struct {
	Message StockUpdateMessage
	event.Envelope
}

func ParseStockUpdateMessage(data []byte) (*StockUpdateEvent, error) {
	var envelope event.Envelope
	if err := json.Unmarshal(data, &envelope); err != nil {
		return nil, err
	}
	var msg StockUpdateMessage
	if err := json.Unmarshal(envelope.Payload, &msg); err != nil {
		return nil, err
	}

	return &StockUpdateEvent{
		Envelope: envelope,
		Message:  msg,
	}, nil
}
