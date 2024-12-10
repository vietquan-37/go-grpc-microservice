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
	"gorm.io/gorm"
)

type OrderHandler struct {
	pb.UnimplementedOrderServiceServer
	ProductClient client.ProductClient
	AuthClient    client.AuthClient
	Repo          repo.IOrderRepo
	DetailRepo    repo.IOrderDetailRepo
}

func NewOrderHandler(productClient client.ProductClient, authClient client.AuthClient, repo repo.IOrderRepo, detailRepo repo.IOrderDetailRepo) *OrderHandler {
	return &OrderHandler{
		ProductClient: productClient,
		AuthClient:    authClient,
		Repo:          repo,
		DetailRepo:    detailRepo,
	}
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
	userID, _, err := extractUserMetadata(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to extract user metadata: %v", err)
	}
	user, err := h.AuthClient.GetOneUser(userID)
	if err != nil {
		return nil, err
	}
	product, err := h.ProductClient.FindOneProduct(req.GetProductId())
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

	err = h.Repo.Transaction(func(repo repo.IOrderRepo) error {
		detail, _ := h.DetailRepo.GetOrderDetailByProductId(product.GetId())
		if detail == nil {
			models := &model.OrderDetail{
				OrderId:   int32(order.ID),
				ProductId: req.ProductId,
				Quantity:  req.GetStock(),
				Price:     price,
			}
			err = h.DetailRepo.CreateOrderDetail(models)
			if err != nil {
				return err
			}
		} else {
			detail.Price += price
			detail.Quantity += req.GetStock()

			err = h.DetailRepo.UpdateOrderDetail(detail)
			if err != nil {
				return err
			}
		}
		order.Amount += price
		err = h.Repo.UpdateOrder(order)
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
	detail, err := h.DetailRepo.GetOrderDetailById(req.GetId())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "order detail not found")
		}
		return nil, status.Errorf(codes.Internal, "error while fetching order detail: %v", err)
	}
	err = h.Repo.Transaction(func(repo repo.IOrderRepo) error {
		err = h.DetailRepo.DeleteOrderDetail(detail)
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
