package repo

import "github.com/vietquan-37/product-service/pkg/model"

type IProductRepo interface {
	CreateProduct(*model.Product) (*model.Product, error)
	FindProduct(int32) (*model.Product, error)
	UpdateProduct(*model.Product) (*model.Product, error)
	DeleteProduct(int32) error
	FindAllProducts() ([]*model.Product, error)
}
