package repo

import "github.com/vietquan-37/order-service/pkg/model"

type IOrderDetailRepo interface {
	CreateOrderDetail(detail *model.OrderDetail) (error error)
	GetOrderDetailByProductId(productId int32) (*model.OrderDetail, error)
	DeleteOrderDetail(model *model.OrderDetail) (err error)
	GetOrderDetailById(id int32) (*model.OrderDetail, error)
	UpdateOrderDetail(model *model.OrderDetail) (err error)
}
