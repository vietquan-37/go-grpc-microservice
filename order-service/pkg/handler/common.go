package handler

import (
	"github.com/vietquan-37/order-service/pkg/enum"
	"github.com/vietquan-37/order-service/pkg/model"
)

func createPendingOrder(userId int32) *model.Order {
	return &model.Order{
		Amount: 0,
		Status: enum.PENDING,
		UserId: uint(userId),
	}
}
