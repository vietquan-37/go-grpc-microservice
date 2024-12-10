package handler

import (
	"context"
	"fmt"
	"github.com/vietquan-37/order-service/pkg/enum"
	"github.com/vietquan-37/order-service/pkg/model"
	"google.golang.org/grpc/metadata"
	"strconv"
)

func createPendingOrder(userId int32) *model.Order {
	return &model.Order{
		Amount: 0,
		Status: enum.PENDING,
		UserId: uint(userId),
	}
}
func extractUserMetadata(ctx context.Context) (userID int32, role string, err error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return 0, "", fmt.Errorf("metadata not found in context")
	}

	idValues := md.Get("id")

	id, err := strconv.Atoi(idValues[0])
	if err != nil {
		return 0, "", fmt.Errorf("invalid id format: %v", err)
	}

	roleValues := md.Get("role")
	if len(roleValues) == 0 {
		return 0, "", fmt.Errorf("role metadata is missing")
	}
	fmt.Printf("ID Values: %v\n", id)           // In giá trị của id để kiểm tra
	fmt.Printf("Role Values: %v\n", roleValues) // In giá trị của role để kiểm tra

	return int32(id), roleValues[0], nil
}
