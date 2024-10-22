package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/zulfikarmuzakir/e_procurement/internal/delivery/http/middleware"
	"github.com/zulfikarmuzakir/e_procurement/internal/domain"
	"github.com/zulfikarmuzakir/e_procurement/internal/domain/converter"
	"github.com/zulfikarmuzakir/e_procurement/pkg/errors"
	"github.com/zulfikarmuzakir/e_procurement/pkg/validator"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type UserHandler struct {
	UserUsecase domain.UserUsecase
	Logger      *zap.Logger
}

func NewUserHandler(u domain.UserUsecase, logger *zap.Logger) *UserHandler {
	return &UserHandler{
		UserUsecase: u,
		Logger:      logger,
	}
}

func (h *UserHandler) RegisterVendor(w http.ResponseWriter, r *http.Request) {
	var user domain.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		h.Logger.Error("Failed to decode request body", zap.Error(err))
		h.sendErrorResponse(w, errors.NewAppError(err, "Invalid request body", http.StatusBadRequest))
		return
	}

	if err := validator.ValidateStruct(user); err != nil {
		h.Logger.Error("Invalid user data", zap.Error(err))
		h.sendErrorResponse(w, errors.NewAppError(err, "Invalid user data", http.StatusBadRequest))
		return
	}

	user.Role = "vendor"
	user.Status = "pending"

	if err := h.UserUsecase.Register(&user); err != nil {
		h.sendErrorResponse(w, err)
		return
	}

	h.Logger.Info("Vendor registered successfully", zap.String("email", user.Email))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Vendor registered successfully",
	})
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var loginRequest struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	}

	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		h.Logger.Error("Failed to decode request body", zap.Error(err))
		h.sendErrorResponse(w, errors.NewAppError(err, "Invalid request body", http.StatusBadRequest))
		return
	}

	if err := validator.ValidateStruct(loginRequest); err != nil {
		h.Logger.Error("Invalid login data", zap.Error(err))
		h.sendErrorResponse(w, errors.NewAppError(err, "Invalid login data", http.StatusBadRequest))
		return
	}

	accessToken, refreshToken, err := h.UserUsecase.Login(loginRequest.Email, loginRequest.Password)
	if err != nil {
		h.sendErrorResponse(w, err)
		return
	}

	user, err := h.UserUsecase.GetByEmail(loginRequest.Email)
	if err != nil {
		h.sendErrorResponse(w, err)
		return
	}

	userResponse := converter.UserToUserResponse(user)

	h.Logger.Info("User logged in successfully", zap.String("email", loginRequest.Email))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"data":          userResponse,
	})
}

func (h *UserHandler) ApproveVendor(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)

	if err := h.UserUsecase.ApproveVendor(id); err != nil {
		h.sendErrorResponse(w, err)
		return
	}

	h.Logger.Info("Vendor approved successfully", zap.Int64("user_id", id))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Vendor approved successfully",
	})
}

func (h *UserHandler) RejectVendor(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)

	if err := h.UserUsecase.RejectVendor(id); err != nil {
		h.sendErrorResponse(w, err)
		return
	}

	h.Logger.Info("Vendor rejected successfully", zap.Int64("user_id", id))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Vendor rejected successfully",
	})
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)

	user, err := h.UserUsecase.GetByID(id)
	if err != nil {
		h.sendErrorResponse(w, err)
		return
	}

	h.Logger.Info("User retrieved successfully", zap.Int64("user_id", id))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) GetAllVendor(w http.ResponseWriter, r *http.Request) {
	users, err := h.UserUsecase.GetAllByRole("vendor")
	if err != nil {
		h.sendErrorResponse(w, err)
		return
	}

	h.Logger.Info("Get all vendors successfully")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)

	var user domain.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		h.Logger.Error("Failed to decode request body", zap.Error(err))
		h.sendErrorResponse(w, errors.NewAppError(err, "Invalid request body", http.StatusBadRequest))
		return
	}

	user.ID = id

	if err := h.UserUsecase.Update(&user); err != nil {
		h.sendErrorResponse(w, err)
		return
	}

	h.Logger.Info("User updated successfully", zap.Int64("user_id", id))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "User updated successfully",
	})
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)

	// Get the ID of the user making the request
	currentUserID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		h.Logger.Error("Failed to get user ID from context")
		h.sendErrorResponse(w, errors.NewAppError(nil, "Unauthorized", http.StatusUnauthorized))
		return
	}

	// Check if the user is trying to delete themselves
	if id == currentUserID {
		h.Logger.Warn("User attempted to delete their own account", zap.Int64("user_id", id))
		h.sendErrorResponse(w, errors.NewAppError(nil, "You cannot delete your own account", http.StatusForbidden))
		return
	}

	if err := h.UserUsecase.Delete(id); err != nil {
		h.sendErrorResponse(w, err)
		return
	}

	h.Logger.Info("User deleted successfully", zap.Int64("user_id", id))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "User deleted successfully",
	})
}

func (h *UserHandler) GetUserMe(w http.ResponseWriter, r *http.Request) {
	currentUserID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		h.Logger.Error("Failed to get user ID from context")
		h.sendErrorResponse(w, errors.NewAppError(nil, "Unauthorized", http.StatusUnauthorized))
		return
	}

	user, err := h.UserUsecase.GetByID(currentUserID)
	if err != nil {
		h.sendErrorResponse(w, err)
		return
	}

	h.Logger.Info("User me retrieved successfully", zap.Int64("user_id", currentUserID))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

// New feature: Refresh token endpoint
func (h *UserHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var refreshRequest struct {
		RefreshToken string `json:"refresh_token" validate:"required"`
	}

	if err := json.NewDecoder(r.Body).Decode(&refreshRequest); err != nil {
		h.Logger.Error("Failed to decode request body", zap.Error(err))
		h.sendErrorResponse(w, errors.NewAppError(err, "Invalid request body", http.StatusBadRequest))
		return
	}

	if err := validator.ValidateStruct(refreshRequest); err != nil {
		h.Logger.Error("Invalid refresh token data", zap.Error(err))
		h.sendErrorResponse(w, errors.NewAppError(err, "Invalid refresh token data", http.StatusBadRequest))
		return
	}

	newAccessToken, err := h.UserUsecase.RefreshToken(refreshRequest.RefreshToken)
	if err != nil {
		h.sendErrorResponse(w, err)
		return
	}

	h.Logger.Info("Access token refreshed successfully")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"access_token": newAccessToken,
	})
}

func (h *UserHandler) sendErrorResponse(w http.ResponseWriter, err error) {
	appErr, ok := err.(*errors.AppError)
	if !ok {
		appErr = errors.NewAppError(err, "Internal server error", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(appErr.Code)
	json.NewEncoder(w).Encode(map[string]string{
		"error":   appErr.Error(),
		"message": appErr.Message,
	})
}
