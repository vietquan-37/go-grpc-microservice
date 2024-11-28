package repo

import (
	"github.com/vietquan-37/order-service/pkg/enum"
	"github.com/vietquan-37/order-service/pkg/model"
	"gorm.io/gorm"
)

type OrderDetailRepo struct {
	DB *gorm.DB
}

func NewOrderDetailRepo(db *gorm.DB) IOrderDetailRepo {
	return &OrderDetailRepo{
		DB: db,
	}
}
func (r *OrderDetailRepo) CreateOrderDetail(detail model.OrderDetail) (error error) {
	err := r.DB.Create(&detail).Error
	if err != nil {
		return err
	}
	return nil
}
func (r *OrderDetailRepo) GetOrderDetailByProductId(productId int32) (*model.OrderDetail, error) {
	var orderDetail model.OrderDetail
	err := r.DB.Where("product_id = ?", productId).First(&orderDetail).Error
	if err != nil {
		return nil, err
	}
	return &orderDetail, nil
}
func (r *OrderDetailRepo) GetPendingOrder() (order *model.Order, err error) {
	err = r.DB.Where("status=?", enum.PENDING).First(&order).Error
	if err != nil {
		return nil, err
	}
	return order, nil
}
