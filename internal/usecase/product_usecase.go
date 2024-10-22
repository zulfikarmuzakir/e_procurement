package usecase

import (
	"net/http"

	"github.com/zulfikarmuzakir/e_procurement/internal/domain"
	"github.com/zulfikarmuzakir/e_procurement/pkg/errors"

	"go.uber.org/zap"
)

type productUsecase struct {
	productRepo domain.ProductRepository
	logger      *zap.Logger
}

func NewProductUsecase(productRepo domain.ProductRepository, logger *zap.Logger) domain.ProductUsecase {
	return &productUsecase{productRepo: productRepo, logger: logger}
}

// CreateProduct implements domain.ProductUsecase.
func (p *productUsecase) CreateProduct(product *domain.Product) error {
	p.logger.Debug("CreateProduct function called", zap.String("product", product.Name))

	if err := p.productRepo.Create(product); err != nil {
		p.logger.Error("Failed to create product", zap.Error(err))
		return errors.NewAppError(err, "Failed to create product", http.StatusInternalServerError)
	}

	p.logger.Info("Product created successfully", zap.String("product", product.Name))
	return nil
}

// GetAll implements domain.ProductUsecase.
func (p *productUsecase) GetAll(name string, limit int, offset int) ([]domain.ProductWithVendor, error) {
	p.logger.Debug("GetAll function called", zap.String("name", name), zap.Int("limit", limit), zap.Int("offset", offset))

	if limit <= 0 {
		limit = 10
	}

	if offset < 0 {
		offset = 0
	}

	products, err := p.productRepo.GetAll(name, limit, offset)
	if err != nil {
		p.logger.Error("Failed to get products", zap.Error(err))
		return nil, errors.NewAppError(err, "Failed to get products", http.StatusInternalServerError)
	}

	if len(products) == 0 {
		p.logger.Info("No products found")
		return []domain.ProductWithVendor{}, nil
	}

	p.logger.Info("Products retrieved successfully", zap.Int("count", len(products)))
	return products, nil
}

// GetProductsByVendorID implements domain.ProductUsecase.
func (p *productUsecase) GetProductsByVendorID(vendorID int64, limit int, offset int) ([]domain.Product, error) {
	p.logger.Debug("GetProductsByVendorID function called", zap.Int64("vendorID", vendorID), zap.Int("limit", limit), zap.Int("offset", offset))

	products, err := p.productRepo.GetProductsByVendorID(vendorID, limit, offset)
	if err != nil {
		p.logger.Error("Failed to get products by vendor ID", zap.Error(err))
		return nil, errors.NewAppError(err, "Failed to get products by vendor ID", http.StatusInternalServerError)
	}

	p.logger.Info("Products retrieved successfully", zap.Int("count", len(products)))
	return products, nil
}

// DeleteProduct implements domain.ProductUsecase.
func (p *productUsecase) DeleteProduct(id int64) error {
	p.logger.Debug("DeleteProduct function called", zap.Int64("id", id))

	if err := p.productRepo.Delete(id); err != nil {
		p.logger.Error("Failed to delete product", zap.Error(err))
		return errors.NewAppError(err, "Failed to delete product", http.StatusInternalServerError)
	}

	p.logger.Info("Product deleted successfully", zap.Int64("id", id))
	return nil
}

// GetProductByID implements domain.ProductUsecase.
func (p *productUsecase) GetProductByID(id int64) (*domain.Product, error) {
	p.logger.Debug("GetProductByID function called", zap.Int64("id", id))

	product, err := p.productRepo.GetByID(id)
	if err != nil {
		p.logger.Error("Failed to get product by ID", zap.Error(err))
		return nil, errors.NewAppError(err, "Failed to get product by ID", http.StatusInternalServerError)
	}

	p.logger.Info("Product retrieved successfully", zap.Int64("id", id))
	return product, nil
}

// UpdateProduct implements domain.ProductUsecase.
func (p *productUsecase) UpdateProduct(product *domain.Product) error {
	p.logger.Debug("UpdateProduct function called", zap.String("product", product.Name))

	// check product is exist
	if product.ID == 0 {
		p.logger.Error("Product ID is required")
		return errors.NewAppError(nil, "Product ID is required", http.StatusBadRequest)
	}

	existingProduct, err := p.productRepo.GetByID(product.ID)
	if err != nil {
		p.logger.Error("Failed to get product by ID", zap.Error(err))
		return errors.NewAppError(err, "Failed to get product by ID", http.StatusInternalServerError)
	}

	if product.VendorID != existingProduct.VendorID {
		p.logger.Error("Product does not belong to the vendor")
		return errors.NewAppError(nil, "Product does not belong to the vendor", http.StatusForbidden)
	}

	if err := p.productRepo.Update(product); err != nil {
		p.logger.Error("Failed to update product", zap.Error(err))
		return errors.NewAppError(err, "Failed to update product", http.StatusInternalServerError)
	}

	p.logger.Info("Product updated successfully", zap.String("product", product.Name))
	return nil
}
