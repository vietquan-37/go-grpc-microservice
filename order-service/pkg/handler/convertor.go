package handler

import (
	"github.com/vietquan-37/order-service/pkg/model"
	"github.com/vietquan-37/order-service/pkg/pb"
)

func convertToCart(order *model.Order, productMap map[int32]string) *pb.UserCartResponse {
	var items []*pb.ItemCart
	for _, detail := range order.OrderDetail {
		name := productMap[detail.ProductId]
		item := convertItemsCart(detail, name)
		items = append(items, item)
	}

	return &pb.UserCartResponse{
		OrderId: int32(order.ID),
		Amount:  float32(order.Amount),
		Status:  string(order.Status),
		Items:   items,
	}
}

func convertItemsCart(detail model.OrderDetail, name string) *pb.ItemCart {
	return &pb.ItemCart{
		Id:          detail.Id,
		ProductId:   detail.ProductId,
		ProductName: name,
		Price:       float32(detail.Price),
		Quantity:    detail.Quantity,
	}
}
func convertItems(detail model.OrderDetail, name string) *pb.Items {
	return &pb.Items{
		Id:          detail.Id,
		ProductId:   detail.ProductId,
		Price:       float32(detail.Price),
		ProductName: name,
		Quantity:    detail.Quantity,
	}
}
func covertToItems(order *model.Order, productMap map[int32]string) []*pb.Items {
	var items []*pb.Items
	for _, orders := range order.OrderDetail {
		name := productMap[orders.ProductId]
		item := convertItems(orders, name)
		items = append(items, item)
	}
	return items
}
