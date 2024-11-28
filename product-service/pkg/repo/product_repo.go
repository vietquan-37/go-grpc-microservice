package repo

import (
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
func (repo *ProductRepo) CreateProduct(model *model.Product) (*model.Product, error) {
	err := repo.DB.Create(&model).Error
	if err != nil {
		return nil, err
	}
	return model, nil
}
func (repo *ProductRepo) FindProduct(Id int32) (*model.Product, error) {
	var product model.Product
	err := repo.DB.Where("id=?", Id).First(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}
func (repo *ProductRepo) UpdateProduct(model *model.Product) (*model.Product, error) {
	err := repo.DB.Save(model).Error
	if err != nil {
		return nil, err
	}
	return model, nil
}
func (repo *ProductRepo) DeleteProduct(Id int32) error {
	err := repo.DB.Where("id=?", Id).Delete(&model.Product{}).Error
	if err != nil {
		return err
	}
	return nil
}
func (repo *ProductRepo) FindAllProducts() ([]*model.Product, error) {
	var products []*model.Product
	err := repo.DB.Find(&products).Error
	if err != nil {
		return nil, err
	}
	return products, nil
}
