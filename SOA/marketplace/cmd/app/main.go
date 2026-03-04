package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	_ "github.com/lib/pq"

	"marketplace/internal/auth"
	"marketplace/internal/config"
	"marketplace/internal/logger"
	"marketplace/internal/repository"
	"marketplace/internal/service"
	"marketplace/internal/transport"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config error: %v", err)
	}

	logg := logger.New("info", "json")

	logg.Info("starting application",
		"env", cfg.AppEnv,
		"port", cfg.Port,
	)

	dsn := buildDSN(cfg)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		logg.Error("failed to open db", "error", err)
		os.Exit(1)
	}

	if err := db.Ping(); err != nil {
		logg.Error("failed to ping db", "error", err)
		os.Exit(1)
	}

	logg.Info("database connected")

	txManager := repository.NewTxManager(db)

	userRepo := repository.NewUserRepository(db)
	productRepo := repository.NewProductRepository(db)
	orderRepo := repository.NewOrderRepository(db)
	promoRepo := repository.NewPromoRepository(db)

	jwtManager := auth.NewJWTManager(
		cfg.JWT.Secret,
		"marketplace",
		cfg.JWT.AccessTTL,
		cfg.JWT.RefreshTTL,
	)

	hasher := auth.NewBcryptHasher(12)

	authService := service.NewAuthService(userRepo, hasher)
	productService := service.NewProductService(productRepo)
	orderService := service.NewOrderService(
		orderRepo,
		productRepo,
		promoRepo,
		txManager,
		logg,
	)
	promoService := service.NewPromoService(promoRepo)

	r := router.NewRouter(router.Dependencies{
		AuthService:    authService,
		OrderService:   orderService,
		ProductService: productService,
		PromoService:   promoService,
		JWTManager:     jwtManager,
	})

	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		logg.Info("http server started", "addr", server.Addr)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logg.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	logg.Info("shutting down server")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logg.Error("graceful shutdown failed", "error", err)
	} else {
		logg.Info("server stopped gracefully")
	}

	if err := db.Close(); err != nil {
		logg.Error("db close error", "error", err)
	}
}

func buildDSN(cfg *config.Config) string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.User,
		cfg.DB.Password,
		cfg.DB.Name,
		cfg.DB.SSLMode,
	)
}
