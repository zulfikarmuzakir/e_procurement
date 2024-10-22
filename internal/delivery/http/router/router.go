package router

import (
	"github.com/zulfikarmuzakir/e_procurement/internal/app"
	"github.com/zulfikarmuzakir/e_procurement/internal/delivery/http/handler"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	customMiddleware "github.com/zulfikarmuzakir/e_procurement/internal/delivery/http/middleware"
)

func SetupRouter(app *app.App) *chi.Mux {
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	userHandler := handler.NewUserHandler(app.UserUsecase, app.Logger)
	productHandler := handler.NewProductHandler(app.ProductUsecase, app.Logger)

	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/login", userHandler.Login)
		r.Post("/register-vendor", userHandler.RegisterVendor)
		r.Get("/products", productHandler.GetAllProducts)
		r.Get("/products/{id}", productHandler.GetProductByID)

		//protected routes
		r.Group(func(r chi.Router) {
			r.Use(customMiddleware.JWTAuth(app.JWTAuth))

			r.Get("/users/{id}", userHandler.GetUser)
			r.Put("/users/{id}", userHandler.UpdateUser)
			r.Get("/users/me", userHandler.GetUserMe)

			r.Group(func(r chi.Router) {
				r.Use(customMiddleware.RoleMiddleware("admin"))
				r.Put("/users/{id}/approve", userHandler.ApproveVendor)
				r.Put("/users/{id}/reject", userHandler.RejectVendor)
				r.Delete("/users/{id}", userHandler.DeleteUser)
				r.Get("/vendors", userHandler.GetAllVendor)
			})

			r.Group(func(r chi.Router) {
				r.Use(customMiddleware.RoleMiddleware("vendor"))
				r.Post("/products", productHandler.CreateProduct)
				r.Put("/products/{id}", productHandler.UpdateProduct)
				r.Delete("/products/{id}", productHandler.DeleteProduct)
				r.Get("/my-products", productHandler.GetMyProducts)
			})
		})

	})

	return r
}
