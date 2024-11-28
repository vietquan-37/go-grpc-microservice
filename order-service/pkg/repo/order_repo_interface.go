package repo

import "github.com/vietquan-37/order-service/pkg/model"

type IOrderRepo interface {
	CreateOrder(order *model.Order) (*model.Order, error)
	GetPendingOrder() (*model.Order, error)
}
