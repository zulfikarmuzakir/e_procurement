package domain

import (
	"time"
)

type Product struct {
	ID        int64     `json:"id"`
	VendorID  int64     `json:"vendor_id"`
	Name      string    `json:"name" validate:"required"`
	Price     int32     `json:"price" validate:"required,gt=0"`
	Stock     int       `json:"stock" validate:"required,gt=0"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
type ProductWithVendor struct {
	ID          int64     `json:"id"`
	VendorID    int64     `json:"vendor_id"`
	ProductName string    `json:"product_name" validate:"required"`
	Price       int32     `json:"price" validate:"required,gt=0"`
	Stock       int       `json:"stock" validate:"required,gt=0"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	VendorName  string    `json:"vendor_name"`
}

type ProductRepository interface {
	Create(product *Product) error
	GetByID(id int64) (*Product, error)
	Update(product *Product) error
	Delete(id int64) error
	GetAll(name string, limit int, offset int) ([]ProductWithVendor, error)
	GetProductsByVendorID(vendorID int64, limit int, offset int) ([]Product, error)
}

type ProductUsecase interface {
	CreateProduct(product *Product) error
	GetProductByID(id int64) (*Product, error)
	UpdateProduct(product *Product) error
	DeleteProduct(id int64) error
	GetAll(name string, limit int, offset int) ([]ProductWithVendor, error)
	GetProductsByVendorID(vendorID int64, limit int, offset int) ([]Product, error)
}
