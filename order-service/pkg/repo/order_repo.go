package repo

import (
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
