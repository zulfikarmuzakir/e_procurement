package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/zulfikarmuzakir/e_procurement/config"
	"github.com/zulfikarmuzakir/e_procurement/internal/app"
	"github.com/zulfikarmuzakir/e_procurement/internal/delivery/http/router"
	"github.com/zulfikarmuzakir/e_procurement/internal/repository/postgres"
	"github.com/zulfikarmuzakir/e_procurement/internal/usecase"
	"github.com/zulfikarmuzakir/e_procurement/pkg/auth"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

func main() {
	// load config
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize Logger
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	// Connect to database
	dbUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)
	db, err := pgxpool.New(context.Background(), dbUrl)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer db.Close()

	jwtAuth := auth.NewJWTAuth(cfg.JWTSecret, cfg.JWTSecret)

	userRepo := postgres.NewUserRepository(db)
	userUsecase := usecase.NewUserUsecase(userRepo, jwtAuth, logger)

	productRepo := postgres.NewProductRepository(db)
	productUsecase := usecase.NewProductUsecase(productRepo, logger)

	app := app.NewApp(userUsecase, productUsecase, jwtAuth, logger)

	r := router.SetupRouter(app)

	logger.Info("Starting server", zap.String("port", "8080"))
	err = http.ListenAndServe(":8080", r)
	if err != nil {
		logger.Fatal("Failed to start server: %v", zap.Error(err))
	}
}
