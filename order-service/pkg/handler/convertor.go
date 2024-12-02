package handler

import (
	"github.com/vietquan-37/order-service/pkg/model"
	"github.com/vietquan-37/order-service/pkg/pb"
)

func convertToCart(order *model.Order) *pb.UserCartResponse {
	var items []*pb.ItemCart
	for _, orders := range order.OrderDetail {
		item := convertItems(orders)
		items = append(items, item)
	}

	return &pb.UserCartResponse{
		OrderId: int32(order.ID),
		Amount:  float32(order.Amount),
		Status:  string(order.Status),
		Items:   items,
	}
}
func convertItems(detail model.OrderDetail) *pb.ItemCart {
	return &pb.ItemCart{
		Id:        detail.Id,
		ProductId: detail.ProductId,
		Price:     float32(detail.Price),
		Quantity:  detail.Quantity,
	}
}
