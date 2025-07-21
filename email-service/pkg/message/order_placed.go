package message

import (
	"common/kafka/event"
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

type OrderEmailEvent struct {
	Message OrderEmailMessage
	event.Envelope
}
type ItemPurchased struct {
	ProductID   int32   `json:"product_id"`
	ProductName string  `json:"product_name"`
	Quantity    int64   `json:"quantity"`
	Price       float64 `json:"price"`
}
