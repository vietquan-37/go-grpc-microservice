package handler

import (
	commonclient "common/client"
	"context"
	"errors"
	"github.com/vietquan-37/order-service/pkg/client"
	"github.com/vietquan-37/order-service/pkg/model"
	"github.com/vietquan-37/order-service/pkg/pb"
	"github.com/vietquan-37/order-service/pkg/repo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type OrderHandler struct {
	pb.UnimplementedOrderServiceServer
	ProductClient client.ProductClient
	AuthClient    commonclient.AuthClient
	Repo          repo.IOrderRepo
}

func NewOrderHandler(productClient client.ProductClient, authClient commonclient.AuthClient, repo repo.IOrderRepo) *OrderHandler {
	return &OrderHandler{
		ProductClient: productClient,
		AuthClient:    authClient,
		Repo:          repo,
	}
}

func (h *OrderHandler) AddProduct(ctx context.Context, req *pb.AddProductRequest) (*pb.CommonResponse, error) {

	userID, _, err := extractUserMetadata(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to extract user metadata: %v", err)
	}
	user, err := h.AuthClient.GetOneUser(userID)
	if err != nil {
		return nil, err
	}
	product, err := h.ProductClient.FindOneProduct(ctx, req.GetProductId())
	if err != nil {
		return nil, err
	}
	if product.Stock < req.Stock {
		return nil, status.Errorf(codes.InvalidArgument, "Product stock is insufficient")
	}
	order, _ := h.Repo.GetPendingOrder(user.UserId)
	if order == nil {
		order, err = h.Repo.CreateOrder(createPendingOrder(user.UserId))
		if err != nil {
			return nil, status.Errorf(codes.Internal, "error while creating order: %v", err)
		}
	}
	price := float64(product.Price * float32(req.GetStock()))
	detail, _ := h.Repo.GetOrderDetailByProductId(product.GetId())
	err = h.Repo.Transaction(func(repo repo.IOrderRepo) error {
		if detail == nil {
			models := &model.OrderDetail{
				OrderId:   int32(order.ID),
				ProductId: req.ProductId,
				Quantity:  req.GetStock(),
				Price:     price,
			}
			err = repo.CreateOrderDetail(models)
			if err != nil {
				return err
			}
		} else {
			detail.Price += price
			detail.Quantity += req.GetStock()

			err = repo.UpdateOrderDetail(detail)
			if err != nil {
				return err
			}
		}
		order.Amount += price
		//rollback
		err = repo.UpdateOrder(order)
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
	userID, _, err := extractUserMetadata(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to extract user metadata: %v", err)
	}
	detail, err := h.Repo.GetOrderDetailById(req.GetId())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "order detail not found")
		}
		return nil, status.Errorf(codes.Internal, "error while fetching order detail: %v", err)
	}
	err = h.Repo.Transaction(func(repo repo.IOrderRepo) error {
		err = repo.DeleteOrderDetail(detail)
		if err != nil {
			return err
		}
		order, err := h.Repo.GetPendingOrder(userID)
		if err != nil {
			return err
		}
		order.Amount -= detail.Price
		err = h.Repo.UpdateOrder(order)
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
	userID, _, err := extractUserMetadata(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to extract user metadata: %v", err)
	}
	user, err := h.AuthClient.GetOneUser(userID)
	if err != nil {
		return nil, err
	}
	order, err := h.Repo.GetPendingOrder(user.UserId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "order not found")
		}
		return nil, status.Errorf(codes.Internal, "error while fetching order: %v", err)
	}
	return convertToCart(order), nil
}
