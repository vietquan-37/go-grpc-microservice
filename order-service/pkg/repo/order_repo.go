package repo

import (
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
func (r *OrderRepo) CreateOrder(order *model.Order) (*model.Order, error) {
	err := r.DB.Create(&order).Error
	if err != nil {
		return nil, err
	}
	return order, nil
}
func (r *OrderRepo) GetPendingOrder(Id int32) (order *model.Order, err error) {
	err = r.DB.Where("status = ? AND user_id = ?", enum.PENDING, Id).Preload("OrderDetail").First(&order).Error
	if err != nil {
		return nil, err
	}
	return order, nil
}
func (r *OrderRepo) UpdateOrder(order *model.Order) error {
	err := r.DB.Save(&order).Error
	if err != nil {
		return err
	}
	return nil
}
func (r *OrderRepo) GetOrderById(orderId int32) (order *model.Order, err error) {
	err = r.DB.Where("id=?", orderId).First(&order).Error
	if err != nil {
		return nil, err
	}
	return order, nil
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
