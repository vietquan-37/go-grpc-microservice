package repository

import (
	"context"
	"github.com/vietquan-37/product-service/pkg/model"
)

type IProductRepo interface {
	CreateProduct(context.Context, *model.Product) (*model.Product, error)
	FindProduct(context.Context, int32) (*model.Product, error)
	UpdateProduct(context.Context, *model.Product) (*model.Product, error)
	DeleteProduct(context.Context, int32) error
	FindAllProducts(context.Context) ([]*model.Product, error)
}
