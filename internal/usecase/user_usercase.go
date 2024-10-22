package usecase

import (
	"net/http"

	"github.com/zulfikarmuzakir/e_procurement/internal/domain"
	"github.com/zulfikarmuzakir/e_procurement/pkg/auth"
	"github.com/zulfikarmuzakir/e_procurement/pkg/errors"
	"github.com/zulfikarmuzakir/e_procurement/pkg/hash"

	"go.uber.org/zap"
)

type userUsecase struct {
	userRepo domain.UserRepository
	jwtAuth  *auth.JWTAuth
	logger   *zap.Logger
}

func NewUserUsecase(userRepo domain.UserRepository, jwtAuth *auth.JWTAuth, logger *zap.Logger) domain.UserUsecase {
	return &userUsecase{
		userRepo: userRepo,
		jwtAuth:  jwtAuth,
		logger:   logger,
	}
}

func (u *userUsecase) Register(user *domain.User) error {
	u.logger.Debug("Register function called", zap.String("email", user.Email), zap.String("password", user.Password))

	hashedPassword, err := hash.HashPassword(user.Password)
	if err != nil {
		u.logger.Error("Failed to hash password", zap.Error(err))
		return errors.NewAppError(err, "Failed to process password", http.StatusInternalServerError)
	}

	user.Password = hashedPassword

	if err := u.userRepo.Create(user); err != nil {
		u.logger.Error("Failed to create user", zap.Error(err))
		return errors.NewAppError(err, err.Error(), http.StatusInternalServerError)
	}

	u.logger.Info("User registered successfully", zap.String("email", user.Email))
	return nil
}

func (u *userUsecase) Login(email, password string) (string, string, error) {
	u.logger.Debug("Login function called", zap.String("email", email), zap.String("password", password))

	user, err := u.userRepo.GetByEmail(email)
	if err != nil {
		u.logger.Warn("Login attempt with non-existent email", zap.String("email", email))
		return "", "", errors.NewAppError(errors.ErrInvalidCredentials, "Invalid email or password", http.StatusUnauthorized)
	}

	// check status active or not
	if user.Status != "active" {
		u.logger.Warn("Login attempt with non-active user", zap.String("email", email))
		return "", "", errors.NewAppError(errors.ErrUserNotActive, "trying to login with non-active user", http.StatusUnauthorized)
	}

	u.logger.Debug("User retrieved from database", zap.String("storedPassword", user.Password))

	err = hash.CheckPasswordHash(password, user.Password)
	if err != nil {
		u.logger.Warn("Login attempt with incorrect password",
			zap.String("email", email),
			zap.String("inputPassword", password),
			zap.String("storedHash", user.Password),
			zap.Error(err))
		return "", "", errors.NewAppError(errors.ErrInvalidCredentials, "Invalid email or password", http.StatusUnauthorized)
	}

	accessToken, err := u.jwtAuth.GenerateAccessToken(user)
	if err != nil {
		u.logger.Error("Failed to generate access token", zap.Error(err))
		return "", "", errors.NewAppError(err, "Failed to generate access token", http.StatusInternalServerError)
	}

	refreshToken, err := u.jwtAuth.GenerateRefreshToken(user)
	if err != nil {
		u.logger.Error("Failed to generate refresh token", zap.Error(err))
		return "", "", errors.NewAppError(err, "Failed to generate refresh token", http.StatusInternalServerError)
	}

	u.logger.Info("User logged in successfully", zap.String("email", email))
	return accessToken, refreshToken, nil
}

// New feature: Refresh token
func (u *userUsecase) RefreshToken(refreshToken string) (string, error) {
	claims, err := u.jwtAuth.ValidateToken(refreshToken, false)
	if err != nil {
		u.logger.Warn("Invalid refresh token", zap.Error(err))
		return "", errors.NewAppError(err, "Invalid refresh token", http.StatusUnauthorized)
	}

	user, err := u.userRepo.GetByID(claims.UserID)
	if err != nil {
		u.logger.Warn("Failed to get user", zap.Error(err), zap.Int64("user_id", claims.UserID))
		return "", errors.NewAppError(errors.ErrUserNotFound, "User not found", http.StatusNotFound)
	}

	newAccessToken, err := u.jwtAuth.GenerateAccessToken(user)
	if err != nil {
		u.logger.Error("Failed to generate new access token", zap.Error(err))
		return "", errors.NewAppError(err, "Failed to generate new access token", http.StatusInternalServerError)
	}

	u.logger.Info("Access token refreshed successfully", zap.Int64("user_id", user.ID))
	return newAccessToken, nil
}

func (u *userUsecase) GetByID(id int64) (*domain.User, error) {
	user, err := u.userRepo.GetByID(id)
	if err != nil {
		u.logger.Warn("Failed to get user", zap.Error(err), zap.Int64("user_id", id))
		return nil, errors.NewAppError(errors.ErrUserNotFound, "User not found", http.StatusNotFound)
	}
	return user, nil
}

func (u *userUsecase) GetByEmail(email string) (*domain.User, error) {
	user, err := u.userRepo.GetByEmail(email)
	if err != nil {
		u.logger.Warn("Failed to get user", zap.Error(err), zap.String("email", email))
		return nil, errors.NewAppError(errors.ErrUserNotFound, "User not found", http.StatusNotFound)
	}
	return user, nil
}

func (u *userUsecase) GetAllByRole(role string) ([]*domain.User, error) {
	users, err := u.userRepo.GetAllByRole(role)
	if err != nil {
		u.logger.Warn("Failed to get users", zap.Error(err), zap.String("role", role))
		return nil, errors.NewAppError(errors.ErrUserNotFound, "Users not found", http.StatusNotFound)
	}
	return users, nil
}

func (u *userUsecase) Update(user *domain.User) error {
	if err := u.userRepo.Update(user); err != nil {
		u.logger.Error("Failed to update user", zap.Error(err), zap.Int64("user_id", user.ID))
		return errors.NewAppError(err, "Failed to update user", http.StatusInternalServerError)
	}
	u.logger.Info("User updated successfully", zap.Int64("user_id", user.ID))
	return nil
}

func (u *userUsecase) ApproveVendor(id int64) error {
	user, err := u.userRepo.GetByID(id)
	if err != nil {
		u.logger.Error("Failed to get user", zap.Error(err), zap.Int64("user_id", id))
		return errors.NewAppError(err, "Failed to get user", http.StatusInternalServerError)
	}
	user.Status = "active"
	if err := u.userRepo.Update(user); err != nil {
		u.logger.Error("Failed to update user status", zap.Error(err), zap.Int64("user_id", id))
		return errors.NewAppError(err, "Failed to update user status", http.StatusInternalServerError)
	}
	u.logger.Info("Vendor approved successfully", zap.Int64("user_id", id))
	return nil
}

func (u *userUsecase) RejectVendor(id int64) error {
	user, err := u.userRepo.GetByID(id)
	if err != nil {
		u.logger.Error("Failed to get user", zap.Error(err), zap.Int64("user_id", id))
		return errors.NewAppError(err, "Failed to get user", http.StatusInternalServerError)
	}
	user.Status = "rejected"
	if err := u.userRepo.Update(user); err != nil {
		u.logger.Error("Failed to update user status", zap.Error(err), zap.Int64("user_id", id))
		return errors.NewAppError(err, "Failed to update user status", http.StatusInternalServerError)
	}
	u.logger.Info("Vendor rejected successfully", zap.Int64("user_id", id))
	return nil
}

func (u *userUsecase) Delete(id int64) error {
	if err := u.userRepo.Delete(id); err != nil {
		u.logger.Error("Failed to delete user", zap.Error(err), zap.Int64("user_id", id))
		return errors.NewAppError(err, "Failed to delete user", http.StatusInternalServerError)
	}
	u.logger.Info("User deleted successfully", zap.Int64("user_id", id))
	return nil
}
