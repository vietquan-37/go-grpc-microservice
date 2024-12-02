package repo

import (
	"github.com/vietquan-37/order-service/pkg/model"
)

type IOrderRepo interface {
	CreateOrder(order *model.Order) (*model.Order, error)
	GetPendingOrder(Id int32) (*model.Order, error)
	UpdateOrder(order *model.Order) error
	GetOrderById(orderId int32) (*model.Order, error)
	Transaction(fn func(repo IOrderRepo) error) error
}
