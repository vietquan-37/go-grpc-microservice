package handler

import (
	"context"
	"errors"
	"github.com/bufbuild/protovalidate-go"
	"github.com/vietquan-37/order-service/pkg/client"
	"github.com/vietquan-37/order-service/pkg/model"
	"github.com/vietquan-37/order-service/pkg/pb"
	"github.com/vietquan-37/order-service/pkg/repo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
)

type OrderHandler struct {
	pb.UnimplementedOrderServiceServer
	ProductClient client.ProductClient
	Repo          repo.IOrderRepo
	DetailRepo    repo.IOrderDetailRepo
}

func NewOrderHandler(productClient client.ProductClient, repo repo.IOrderRepo, detailRepo repo.IOrderDetailRepo) *OrderHandler {
	return &OrderHandler{
		ProductClient: productClient,
		Repo:          repo,
		DetailRepo:    detailRepo,
	}
}
func (h *OrderHandler) InitialOrder(ctx context.Context, req *pb.InitialOrderRequest) (*emptypb.Empty, error) {

	return nil, status.Errorf(codes.Unimplemented, "method InitialOrder not implemented")
}
func (h *OrderHandler) AddProduct(ctx context.Context, req *pb.AddProductRequest) (*pb.CommonResponse, error) {
	validator, err := protovalidate.New()
	if err != nil {
		panic(err)
	}
	if err := validator.Validate(req); err != nil {
		violation := ErrorResponses(err)
		return nil, invalidArgumentError(violation)
	}
	order, err := h.Repo.GetPendingOrder()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {

		}
		return nil, status.Errorf(codes.Internal, "Error while fetching pending order")
	}
	product, err := h.ProductClient.FindOneProduct(req.GetProductId())
	if err != nil {
		return nil, err
	}
	if product.Stock < req.Stock {
		return nil, status.Errorf(codes.InvalidArgument, "Product stock is insufficient")
	}
	detail := model.OrderDetail{
		OrderId:   int32(order.ID),
		ProductId: req.ProductId,
		Quantity:  req.GetStock(),
		Price:     float64(product.Price),
	}
	err = h.DetailRepo.CreateOrderDetail(detail)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error while adding product to order detail: %v", err)
	}
	return &pb.CommonResponse{
		Message: "add product successfully",
	}, nil
}
