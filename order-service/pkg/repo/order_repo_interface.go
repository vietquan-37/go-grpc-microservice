package repo

import (
	"context"
	"github.com/vietquan-37/order-service/pkg/model"
)

type IOrderRepo interface {
	CreateOrder(ctx context.Context, order *model.Order) (*model.Order, error)
	GetPendingOrder(ctx context.Context, Id int32) (*model.Order, error)
	UpdateOrder(ctx context.Context, order *model.Order) error
	GetOrderById(ctx context.Context, orderId int32) (*model.Order, error)
	CreateOrderDetail(ctx context.Context, detail *model.OrderDetail) (error error)
	GetOrderDetailByProductId(ctx context.Context, productId int32) (*model.OrderDetail, error)
	DeleteOrderDetail(ctx context.Context, model *model.OrderDetail) (err error)
	GetOrderDetailById(ctx context.Context, id int32) (*model.OrderDetail, error)
	UpdateOrderDetail(ctx context.Context, model *model.OrderDetail) (err error)
	Transaction(fn func(repo IOrderRepo) error) error
}
