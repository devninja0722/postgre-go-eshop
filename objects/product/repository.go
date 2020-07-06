package product

import (
	"errors"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"github.com/jinzhu/gorm"
	"goshop/dbs"
	"goshop/utils"
)

type Repository interface {
	GetProducts() (*[]Product, error)
	GetProductByID(uuid string) (*Product, error)
	CreateProduct(req *ProductRequest) (*Product, error)
	UpdateProduct(uuid string, req *ProductRequest) (*Product, error)

	beforeCreate(product *Product)
}

type repo struct {
	db *gorm.DB
}

func NewRepository() Repository {
	return &repo{db: dbs.Database}
}

func (r *repo) GetProducts() (*[]Product, error) {
	var products []Product
	if r.db.Find(&products).RecordNotFound() {
		return nil, nil
	}

	return &products, nil

}

func (r *repo) GetProductByID(uuid string) (*Product, error) {
	var product Product
	if r.db.Where("uuid = ?", uuid).Find(&product).RecordNotFound() {
		return nil, errors.New("not found product")
	}

	return &product, nil
}

func (r *repo) CreateProduct(req *ProductRequest) (*Product, error) {
	var product Product
	copier.Copy(&product, &req)

	r.beforeCreate(&product)
	r.db.Create(&product)

	return &product, nil
}

func (r *repo) UpdateProduct(uuid string, req *ProductRequest) (*Product, error) {
	var product Product
	if r.db.Where("uuid = ? ", uuid).First(&product).RecordNotFound() {
		return nil, errors.New("not found product")
	}

	copier.Copy(&product, &req)
	r.db.Save(&product)

	return &product, nil
}

func (r *repo) beforeCreate(product *Product) {
	product.UUID = uuid.New().String()
	product.Code = utils.GenerateCode("P")
}