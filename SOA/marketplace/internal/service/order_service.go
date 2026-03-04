package service

import (
	"context"
	"log/slog"
	"time"

	"marketplace/internal/domain"
	"marketplace/internal/repository"

	"github.com/google/uuid"
)

type OrderService interface {
	CreateOrder(ctx context.Context, userID string, items []domain.OrderItem, promoCode *string) (*domain.Order, error)
	GetByID(ctx context.Context, id string) (*domain.Order, error)
	ChangeStatus(ctx context.Context, id string, status domain.OrderStatus) error
}

type orderService struct {
	orderRepo   repository.OrderRepository
	productRepo repository.ProductRepository
	promoRepo   repository.PromoRepository
	txManager   repository.TxManager
	logger		*slog.Logger
}

func NewOrderService(
	orderRepo repository.OrderRepository,
	productRepo repository.ProductRepository,
	promoRepo repository.PromoRepository,
	txManager repository.TxManager,
	log *slog.Logger,
) OrderService {
	return &orderService{
		orderRepo:   orderRepo,
		productRepo: productRepo,
		promoRepo:   promoRepo,
		txManager:   txManager,
		logger: log,
	}
}

func (s *orderService) CreateOrder(
	ctx context.Context,
	userID string,
	items []domain.OrderItem,
	promoCode *string,
) (*domain.Order, error) {
	s.logger.Info("creating order")

	if len(items) == 0 {
		s.logger.Error("failed", "error", domain.ErrInvalidInput)
		return nil, domain.ErrInvalidInput
	}

	order := &domain.Order{
		ID:     uuid.NewString(),
		UserID: userID,
		Status: domain.OrderStatusPending,
		Items:  items,
	}

	err := s.txManager.WithinTx(ctx, func(txCtx context.Context) error {

		for i := range order.Items {
			product, err := s.productRepo.GetByID(txCtx, order.Items[i].ProductID)
			if err != nil {
				s.logger.Error("failed to create order", "error", err)
				return domain.ErrNotFound
			}

			if err := product.DecreaseStock(order.Items[i].Quantity); err != nil {
				s.logger.Error("failed to create order", "error", err)
				return err
			}

			if err := s.productRepo.UpdateStock(txCtx, product.ID, product.Stock); err != nil {
				s.logger.Error("failed to create order", "error", err)
				return err
			}

			order.Items[i].Price = product.Price
		}

		order.CalculateTotal()

		if promoCode != nil {
			promo, err := s.promoRepo.GetByCode(txCtx, *promoCode)
			if err != nil {
				s.logger.Error("failed to create order", "error", err)
				return domain.ErrNotFound
			}

			if err := order.ApplyPromo(promo, time.Now()); err != nil {
				s.logger.Error("failed to create order", "error", err)
				return err
			}

			order.PromoID = &promo.ID

			if err := s.promoRepo.IncrementUsage(txCtx, promo.ID); err != nil {
				s.logger.Error("failed to create order", "error", err)
				return err
			}
		}

		if err := s.orderRepo.Create(txCtx, order); err != nil {
			s.logger.Error("failed to create order", "error", err)
			return err
		}

		if err := s.orderRepo.AddItems(txCtx, order.ID, order.Items); err != nil {
			s.logger.Error("failed to create order", "error", err)
			return err
		}

		return nil
	})

	if err != nil {
		s.logger.Error("failed to create order", "error", err)
		return nil, err
	}

	s.logger.Info("order created")
	return order, nil
}

func (s *orderService) GetByID(ctx context.Context, id string) (*domain.Order, error) {
	o, err := s.orderRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("failed to get by id", "error", err)
		return nil, domain.ErrNotFound
	}
	return o, nil
}

func (s *orderService) ChangeStatus(ctx context.Context, id string, status domain.OrderStatus) error {
	order, err := s.orderRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("failed to change status", "error", err)
		return domain.ErrNotFound
	}

	if err := order.ChangeStatus(status); err != nil {
		s.logger.Error("failed to change status", "error", err)
		return err
	}

	return s.orderRepo.UpdateStatus(ctx, id, status)
}
