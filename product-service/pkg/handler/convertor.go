package handler

import (
	"github.com/vietquan-37/product-service/pkg/model"
	"github.com/vietquan-37/product-service/pkg/pb"
)

func convertToProduct(req *pb.CreateProductRequest) *model.Product {
	return &model.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
	}
}

func convertToProductResponse(model *model.Product) *pb.ProductResponse {
	return &pb.ProductResponse{
		Id:          int32(model.ID),
		Name:        model.Name,
		Description: model.Description,
		Price:       model.Price,
		Stock:       model.Stock,
	}

}
