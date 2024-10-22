package postgres

import (
	"context"

	"github.com/zulfikarmuzakir/e_procurement/internal/domain"
	postgres "github.com/zulfikarmuzakir/e_procurement/internal/repository/postgres/sqlc"

	"github.com/jackc/pgx/v5/pgxpool"
)

type productRepository struct {
	q *postgres.Queries
}

func NewProductRepository(db *pgxpool.Pool) domain.ProductRepository {
	return &productRepository{q: postgres.New(db)}
}

// CreateProduct implements domain.ProductRepository.
func (p *productRepository) Create(product *domain.Product) error {
	ctx := context.Background()
	_, err := p.q.CreateProduct(ctx, postgres.CreateProductParams{
		VendorID: int32(product.VendorID),
		Name:     product.Name,
		Price:    product.Price,
		Stock:    int32(product.Stock),
	})

	return err
}

func (p *productRepository) GetAll(name string, limit int, offset int) ([]domain.ProductWithVendor, error) {
	ctx := context.Background()
	products, err := p.q.GetProductsWithVendor(ctx, postgres.GetProductsWithVendorParams{
		Column1: name,
		Limit:   int32(limit),
		Offset:  int32(offset),
	})

	if err != nil {
		return nil, err
	}

	var domainProducts []domain.ProductWithVendor
	for _, product := range products {
		domainProducts = append(domainProducts, domain.ProductWithVendor{
			ID:          int64(product.ID),
			VendorID:    int64(product.VendorID),
			ProductName: product.ProductName,
			Price:       product.Price,
			Stock:       int(product.Stock),
			VendorName:  product.VendorName.String,
		})
	}

	return domainProducts, nil
}

func (p *productRepository) GetProductsByVendorID(vendorID int64, limit int, offset int) ([]domain.Product, error) {
	ctx := context.Background()
	products, err := p.q.GetProductsByVendorID(ctx, postgres.GetProductsByVendorIDParams{
		VendorID: int32(vendorID),
		Limit:    int32(limit),
		Offset:   int32(offset),
	})

	if err != nil {
		return nil, err
	}

	var domainProducts []domain.Product
	for _, product := range products {
		domainProducts = append(domainProducts, domain.Product{
			ID:       int64(product.ID),
			VendorID: int64(product.VendorID),
			Name:     product.Name,
			Price:    product.Price,
			Stock:    int(product.Stock),
		})
	}

	return domainProducts, nil
}

// DeleteProduct implements domain.ProductRepository.
func (p *productRepository) Delete(id int64) error {
	ctx := context.Background()
	err := p.q.DeleteProduct(ctx, int32(id))

	return err
}

// GetProductByID implements domain.ProductRepository.
func (p *productRepository) GetByID(id int64) (*domain.Product, error) {
	ctx := context.Background()
	product, err := p.q.GetProductByID(ctx, int32(id))

	return &domain.Product{
		ID:       int64(product.ID),
		VendorID: int64(product.VendorID),
		Name:     product.Name,
		Price:    product.Price,
		Stock:    int(product.Stock),
	}, err
}

// UpdateProduct implements domain.ProductRepository.
func (p *productRepository) Update(product *domain.Product) error {
	ctx := context.Background()
	err := p.q.UpdateProduct(ctx, postgres.UpdateProductParams{
		ID:    int32(product.ID),
		Name:  product.Name,
		Price: product.Price,
		Stock: int32(product.Stock),
	})

	return err
}
