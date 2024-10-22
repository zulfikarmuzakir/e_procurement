package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/zulfikarmuzakir/e_procurement/internal/delivery/http/middleware"
	"github.com/zulfikarmuzakir/e_procurement/internal/domain"
	"github.com/zulfikarmuzakir/e_procurement/pkg/errors"
	"github.com/zulfikarmuzakir/e_procurement/pkg/validator"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type ProductHandler struct {
	ProductUsecase domain.ProductUsecase
	Logger         *zap.Logger
}

func NewProductHandler(productUsecase domain.ProductUsecase, logger *zap.Logger) *ProductHandler {
	return &ProductHandler{
		ProductUsecase: productUsecase,
		Logger:         logger,
	}
}

func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product domain.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		h.Logger.Error("Failed to decode request body", zap.Error(err))
		h.sendErrorResponse(w, errors.NewAppError(err, "Invalid request body", http.StatusBadRequest))
		return
	}

	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		h.Logger.Error("Failed to get user ID from context")
		h.sendErrorResponse(w, errors.NewAppError(nil, "Unauthorized", http.StatusUnauthorized))
		return
	}

	product.VendorID = userID

	if err := validator.ValidateStruct(product); err != nil {
		h.sendValidationErrorResponse(w, err)
		return
	}

	if err := h.ProductUsecase.CreateProduct(&product); err != nil {
		h.sendErrorResponse(w, err)
		return
	}

	h.Logger.Info("Product created successfully", zap.String("product", product.Name), zap.Int64("vendorID", product.VendorID))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Product created successfully",
		"data":    product,
	})
}

func (h *ProductHandler) GetAllProducts(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	limit, _ := strconv.ParseInt(r.URL.Query().Get("limit"), 10, 64)
	if limit == 0 {
		limit = 10
	}
	offset, _ := strconv.ParseInt(r.URL.Query().Get("offset"), 10, 64)

	if name == "" {
		fmt.Println("name is empty")
		name = "%"
	}

	products, err := h.ProductUsecase.GetAll(name, int(limit), int(offset))
	if err != nil {
		h.sendErrorResponse(w, err)
		return
	}

	h.Logger.Info("Products retrieved successfully", zap.Int("count", len(products)))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Products retrieved successfully",
		"data":    products,
	})
}

func (h *ProductHandler) GetMyProducts(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		h.Logger.Error("Failed to get user ID from context")
		h.sendErrorResponse(w, errors.NewAppError(nil, "Unauthorized", http.StatusUnauthorized))
		return
	}

	limit, _ := strconv.ParseInt(r.URL.Query().Get("limit"), 10, 64)
	if limit == 0 {
		limit = 10
	}
	offset, _ := strconv.ParseInt(r.URL.Query().Get("offset"), 10, 64)

	products, err := h.ProductUsecase.GetProductsByVendorID(userID, int(limit), int(offset))
	if err != nil {
		h.sendErrorResponse(w, err)
		return
	}

	h.Logger.Info("My products retrieved successfully", zap.Int("count", len(products)))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "My products retrieved successfully",
		"data":    products,
	})
}

func (h *ProductHandler) GetProductByID(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)

	product, err := h.ProductUsecase.GetProductByID(id)
	if err != nil {
		h.sendErrorResponse(w, err)
		return
	}

	h.Logger.Info("Product retrieved successfully", zap.Int64("product_id", id))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Product retrieved successfully",
		"data":    product,
	})
}

func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)

	var product domain.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		h.Logger.Error("Failed to decode request body", zap.Error(err))
		h.sendErrorResponse(w, errors.NewAppError(err, "Invalid request body", http.StatusBadRequest))
		return
	}

	// validate the product
	if err := validator.ValidateStruct(product); err != nil {
		h.sendValidationErrorResponse(w, err)
		return
	}

	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		h.Logger.Error("Failed to get user ID from context")
		h.sendErrorResponse(w, errors.NewAppError(nil, "Unauthorized", http.StatusUnauthorized))
		return
	}

	product.ID = id
	product.VendorID = userID

	if err := h.ProductUsecase.UpdateProduct(&product); err != nil {
		h.sendErrorResponse(w, err)
		return
	}

	h.Logger.Info("Product updated successfully", zap.Int64("product_id", id))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Product updated successfully",
		"data":    product,
	})
}

func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)

	if err := h.ProductUsecase.DeleteProduct(id); err != nil {
		h.sendErrorResponse(w, err)
		return
	}

	h.Logger.Info("Product deleted successfully", zap.Int64("product_id", id))
	w.WriteHeader(http.StatusNoContent)
}

func (h *ProductHandler) sendValidationErrorResponse(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)

	validationErrors := validator.GetValidationErrors(err)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": "Validation failed",
		"data":  validationErrors,
	})
}

func (h *ProductHandler) sendErrorResponse(w http.ResponseWriter, err error) {
	fmt.Println(err)
	appErr, ok := err.(*errors.AppError)
	if !ok {
		appErr = errors.NewAppError(err, "Internal server error", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(appErr.Code)
	json.NewEncoder(w).Encode(map[string]string{"error": appErr.Message})
}
