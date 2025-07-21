package message

import (
	"common/kafka/event"
	"encoding/json"
	"time"
)

type StockUpdateMessage struct {
	OrderID  int32           `json:"order_id"`
	Customer Customer        `json:"customer"`
	Items    []ItemPurchased `json:"items"`
}

func NewStockUpdateEnvelope(source, version string, eventType string, payload StockUpdateMessage) (*event.Envelope, error) {
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
