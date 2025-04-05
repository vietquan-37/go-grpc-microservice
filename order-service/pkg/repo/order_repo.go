package repo

import (
	"context"
	"github.com/vietquan-37/order-service/pkg/enum"
	"github.com/vietquan-37/order-service/pkg/model"
	"gorm.io/gorm"
)

type OrderRepo struct {
	DB *gorm.DB
}

func NewOrderRepo(db *gorm.DB) IOrderRepo {
	return &OrderRepo{
		DB: db,
	}
}
func (r *OrderRepo) CreateOrder(ctx context.Context, order *model.Order) (*model.Order, error) {
	err := r.DB.WithContext(ctx).Create(&order).Error
	if err != nil {
		return nil, err
	}
	return order, nil
}
func (r *OrderRepo) GetPendingOrder(ctx context.Context, Id int32) (order *model.Order, err error) {
	err = r.DB.WithContext(ctx).Where("status = ? AND user_id = ?", enum.PENDING, Id).Preload("OrderDetail").First(&order).Error
	if err != nil {
		return nil, err
	}
	return order, nil
}
func (r *OrderRepo) UpdateOrder(ctx context.Context, order *model.Order) error {
	err := r.DB.WithContext(ctx).Save(&order).Error
	if err != nil {
		return err
	}
	return nil
}
func (r *OrderRepo) GetOrderById(ctx context.Context, orderId int32) (order *model.Order, err error) {
	err = r.DB.WithContext(ctx).Where("id=?", orderId).First(&order).Error
	if err != nil {
		return nil, err
	}

	return order, nil

}
func (r *OrderRepo) CreateOrderDetail(ctx context.Context, detail *model.OrderDetail) (error error) {
	err := r.DB.WithContext(ctx).Create(&detail).Error
	if err != nil {
		return err
	}
	return nil
}
func (r *OrderRepo) GetOrderDetailByProductId(ctx context.Context, productId int32) (*model.OrderDetail, error) {
	var orderDetail model.OrderDetail
	err := r.DB.WithContext(ctx).Where("product_id = ?", productId).First(&orderDetail).Error
	if err != nil {
		return nil, err
	}
	return &orderDetail, nil
}
func (r *OrderRepo) DeleteOrderDetail(ctx context.Context, model *model.OrderDetail) (err error) {

	err = r.DB.WithContext(ctx).Unscoped().Delete(&model).Error
	if err != nil {
		return err
	}
	return nil
}
func (r *OrderRepo) GetOrderDetailById(ctx context.Context, Id int32) (orderDetail *model.OrderDetail, err error) {
	err = r.DB.WithContext(ctx).Where("id = ?", Id).Model(&orderDetail).First(&orderDetail).Error
	if err != nil {
		return nil, err
	}
	return orderDetail, nil
}
func (r *OrderRepo) UpdateOrderDetail(ctx context.Context, detail *model.OrderDetail) (err error) {
	err = r.DB.WithContext(ctx).Save(&detail).Error
	if err != nil {
		return err
	}
	return nil
}
func (r *OrderRepo) Transaction(fn func(repo IOrderRepo) error) error {
	tx := r.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	repo := NewOrderRepo(tx)
	err := fn(repo)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}
