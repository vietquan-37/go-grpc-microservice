package handler

import (
	"context"
	"errors"
	"github.com/vietquan-37/product-service/pkg/pb"
	"github.com/vietquan-37/product-service/pkg/repo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
	"sync"
)

var wait sync.WaitGroup

type ProductHandler struct {
	pb.UnimplementedProductServiceServer
	Repo repo.IProductRepo
}

func NewProductHandler(repo repo.IProductRepo) *ProductHandler {
	return &ProductHandler{
		Repo: repo,
	}
}
func (h *ProductHandler) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.ProductResponse, error) {

	model := convertToProduct(req)
	product, err := h.Repo.CreateProduct(model)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error while creating product : %s", err.Error())
	}
	rsp := convertToProductResponse(product)
	return rsp, nil
}
func (h *ProductHandler) UpdateProduct(ctx context.Context, req *pb.UpdateProductRequest) (*pb.ProductResponse, error) {
	product, err := h.Repo.FindProduct(req.GetId())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "product not found")
		}
		return nil, status.Errorf(codes.Internal, "error while finding product : %s", err.Error())
	}
	product.Price = req.GetPrice()
	product.Stock = req.GetStock()
	product, err = h.Repo.UpdateProduct(product)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error while updating product : %s", err.Error())
	}
	rsp := convertToProductResponse(product)
	return rsp, nil
}
func (h *ProductHandler) DecreaseProductStock(ctx context.Context, req *pb.DecreaseStockRequest) (*pb.CommonResponse, error) {
	product, err := h.Repo.FindProduct(req.GetProductId())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "product not found")
		}
		return nil, status.Errorf(codes.Internal, "error while finding product : %s", err.Error())
	}
	if product.Stock < req.Quantity {
		return nil, status.Errorf(codes.InvalidArgument, "product %s stock is insufficient", product.Name)
	}
	product.Stock -= req.Quantity
	product, err = h.Repo.UpdateProduct(product)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error while updating product : %s", err.Error())
	}
	return &pb.CommonResponse{
		Message: "Product stock updated",
	}, nil
}
func (h *ProductHandler) FindOneProduct(ctx context.Context, req *pb.ProductRequest) (*pb.ProductResponse, error) {
	//time.Sleep(5 * time.Second)
	product, err := h.Repo.FindProduct(req.GetId())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "product not found")
		}
		return nil, status.Errorf(codes.Internal, "error while finding product : %s", err.Error())
	}
	rsp := convertToProductResponse(product)
	//wait.Add(1)
	//go func() {
	//	defer wait.Done()
	//	time.Sleep(5 * time.Second)
	//	log.Println("email have been sent")
	//}()
	return rsp, nil
}
func (h *ProductHandler) FindAllProduct(ctx context.Context, req *emptypb.Empty) (*pb.ProductResponseList, error) {
	products, err := h.Repo.FindAllProducts()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error while finding all products : %s", err.Error())
	}
	var productList []*pb.ProductResponse
	for _, product := range products {
		product := convertToProductResponse(product)
		productList = append(productList, product)
	}
	return &pb.ProductResponseList{
		Products: productList,
	}, nil
}
func (h *ProductHandler) DeleteProduct(context.Context, *pb.ProductRequest) (*pb.CommonResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteProduct not implemented")
}
