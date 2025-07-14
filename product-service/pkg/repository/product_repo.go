package repository

import (
	"context"
	"github.com/vietquan-37/product-service/pkg/model"
	"gorm.io/gorm"
)

type ProductRepo struct {
	DB *gorm.DB
}

func NewProductRepo(db *gorm.DB) IProductRepo {
	return &ProductRepo{
		DB: db,
	}
}
func (repo *ProductRepo) CreateProduct(ctx context.Context, model *model.Product) (*model.Product, error) {
	err := repo.DB.WithContext(ctx).Create(&model).Error
	if err != nil {
		return nil, err
	}
	return model, nil
}
func (repo *ProductRepo) FindProduct(ctx context.Context, Id int32) (*model.Product, error) {
	var product model.Product
	err := repo.DB.WithContext(ctx).Where("id=?", Id).First(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}
func (repo *ProductRepo) UpdateProduct(ctx context.Context, model *model.Product) (*model.Product, error) {
	err := repo.DB.WithContext(ctx).Save(model).Error
	if err != nil {
		return nil, err
	}
	return model, nil
}
func (repo *ProductRepo) DeleteProduct(ctx context.Context, Id int32) error {
	err := repo.DB.WithContext(ctx).Where("id = ?", Id).Delete(&model.Product{}).Error
	if err != nil {
		return err
	}
	return nil
}

func (repo *ProductRepo) FindAllProducts(ctx context.Context) ([]*model.Product, error) {
	var products []*model.Product
	err := repo.DB.WithContext(ctx).Find(&products).Error
	if err != nil {
		return nil, err
	}
	return products, nil
}
func (repo *ProductRepo) FindProductsByIds(ctx context.Context, ids []int32) ([]*model.Product, error) {
	var products []*model.Product
	err := repo.DB.WithContext(ctx).Where("id in (?)", ids).Find(&products).Error
	if err != nil {
		return nil, err
	}
	return products, nil
}
