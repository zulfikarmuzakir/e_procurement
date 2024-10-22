package app

import (
	"github.com/zulfikarmuzakir/e_procurement/internal/domain"
	"github.com/zulfikarmuzakir/e_procurement/pkg/auth"

	"go.uber.org/zap"
)

type App struct {
	UserUsecase    domain.UserUsecase
	ProductUsecase domain.ProductUsecase
	JWTAuth        *auth.JWTAuth
	Logger         *zap.Logger
}

func NewApp(userUsecase domain.UserUsecase, productUsecase domain.ProductUsecase, jwtAuth *auth.JWTAuth, logger *zap.Logger) *App {
	return &App{UserUsecase: userUsecase, ProductUsecase: productUsecase, JWTAuth: jwtAuth, Logger: logger}
}

// You can add more methods here if needed, such as initialization or shutdown procedures
func (a *App) Initialize() error {
	// Perform any necessary initialization
	a.Logger.Info("Initializing application")
	return nil
}

func (a *App) Shutdown() error {
	// Perform any necessary cleanup
	a.Logger.Info("Shutting down application")
	return nil
}
