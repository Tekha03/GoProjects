package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"

	"marketplace/internal/auth"
	appmiddleware "marketplace/internal/transport/middleware"
	"marketplace/internal/service"
	"marketplace/internal/transport/handlers"
)

type Dependencies struct {
	AuthService    service.AuthService
	OrderService   service.OrderService
	ProductService service.ProductService
	PromoService   service.PromoService
	JWTManager     *auth.JWTManager
}

func NewRouter(deps Dependencies) http.Handler {
	r := chi.NewRouter()

	r.Use(chimiddleware.Recoverer)
	r.Use(appmiddleware.RequestID)
	r.Use(appmiddleware.Logging)

	authHandler := handlers.NewAuthHandler(deps.AuthService, deps.JWTManager)
	orderHandler := handlers.NewOrderHandler(deps.OrderService)
	productHandler := handlers.NewProductHandler(deps.ProductService)
	promoHandler := handlers.NewPromoHandler(deps.PromoService)

	r.Route("/api/v1", func(r chi.Router) {

		r.Group(func(r chi.Router) {

			r.Post("/auth/register", authHandler.Register)
			r.Post("/auth/login", authHandler.Login)

			r.Get("/products", productHandler.List)
			r.Get("/products/{id}", productHandler.GetByID)

			r.Get("/promos/{code}", promoHandler.GetByCode)
			r.Get("/promos/{code}/validate", promoHandler.Validate)
		})

		r.Group(func(r chi.Router) {

			r.Use(appmiddleware.Auth(deps.JWTManager))

			r.Post("/orders", orderHandler.Create)
			r.Get("/orders/{id}", orderHandler.GetByID)
			r.Patch("/orders/{id}/status", orderHandler.ChangeStatus)

			r.With(appmiddleware.RequireRole("admin")).
				Post("/products", productHandler.Create)
		})
	})

	return r
}
