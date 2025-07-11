package handler

import (
	commonclient "common/client"
	"common/extract"
	"context"
	"errors"
	"github.com/vietquan-37/order-service/pkg/client"
	"github.com/vietquan-37/order-service/pkg/model"
	"github.com/vietquan-37/order-service/pkg/pb"
	"github.com/vietquan-37/order-service/pkg/repository"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type OrderHandler struct {
	pb.UnimplementedOrderServiceServer
	ProductClient *client.ProductClient
	AuthClient    *commonclient.AuthClient
	Repo          repository.IOrderRepo
}

func NewOrderHandler(productClient *client.ProductClient, authClient *commonclient.AuthClient, repo repository.IOrderRepo) *OrderHandler {
	return &OrderHandler{
		ProductClient: productClient,
		AuthClient:    authClient,
		Repo:          repo,
	}
}

func (h *OrderHandler) AddProduct(ctx context.Context, req *pb.AddProductRequest) (*pb.CommonResponse, error) {

	metadata, err := extract.UsersMetadata(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to extract user mtdt: %v", err)
	}

	product, err := h.ProductClient.FindOneProduct(ctx, req.GetProductId())
	if err != nil {
		return nil, err
	}
	if product.Stock < req.Stock {
		return nil, status.Errorf(codes.InvalidArgument, "Product stock is insufficient")
	}
	order, _ := h.Repo.GetPendingOrder(ctx, metadata.User.UserId)
	if order == nil {
		order, err = h.Repo.CreateOrder(ctx, createPendingOrder(metadata.User.UserId))
		if err != nil {
			return nil, status.Errorf(codes.Internal, "error while creating order: %v", err)
		}
	}
	price := float64(product.Price * float32(req.GetStock()))
	detail, _ := h.Repo.GetOrderDetailByProductId(ctx, product.GetId())
	if detail != nil && detail.Quantity+req.GetStock() > product.Stock {
		return nil, status.Errorf(codes.InvalidArgument, "Product stock is insufficient")
	}
	err = h.Repo.Transaction(func(repo repository.IOrderRepo) error {
		if detail == nil {
			models := &model.OrderDetail{
				OrderId:   int32(order.ID),
				ProductId: req.ProductId,
				Quantity:  req.GetStock(),
				Price:     price,
			}
			err = repo.CreateOrderDetail(ctx, models)
			if err != nil {
				return err
			}
		} else {
			detail.Price += price
			detail.Quantity += req.GetStock()

			err = repo.UpdateOrderDetail(ctx, detail)
			if err != nil {
				return err
			}
		}
		order.Amount += price
		err = repo.UpdateOrder(ctx, order)
		if err != nil {
			return err
		}
		return nil
	},
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error while updating order: %v", err)
	}

	return &pb.CommonResponse{
		Message: "add product successfully",
	}, nil
}
func (h *OrderHandler) DeleteDetail(ctx context.Context, req *pb.DeleteDetailRequest) (*pb.CommonResponse, error) {
	metadata, err := extract.UsersMetadata(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to extract user mtdt: %v", err)
	}
	detail, err := h.Repo.GetOrderDetailById(ctx, req.GetId())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "order detail not found")
		}
		return nil, status.Errorf(codes.Internal, "error while fetching order detail: %v", err)
	}
	err = h.Repo.Transaction(func(repo repository.IOrderRepo) error {
		err = repo.DeleteOrderDetail(ctx, detail)
		if err != nil {
			return err
		}
		order, err := h.Repo.GetPendingOrder(ctx, metadata.User.UserId)
		if err != nil {
			return err
		}
		order.Amount -= detail.Price
		err = h.Repo.UpdateOrder(ctx, order)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error while deleting detail: %v", err)
	}

	return &pb.CommonResponse{
		Message: "Delete detail successfully",
	}, nil
}
func (h *OrderHandler) GetUserCart(ctx context.Context, req *pb.UserCartRequest) (*pb.UserCartResponse, error) {
	metadata, err := extract.UsersMetadata(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to extract user mtdt: %v", err)
	}
	order, err := h.Repo.GetPendingOrder(ctx, metadata.User.UserId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "order not found")
		}
		return nil, status.Errorf(codes.Internal, "error while fetching order: %v", err)
	}
	return convertToCart(order), nil
}
