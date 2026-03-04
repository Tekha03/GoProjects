package service

import (
	"context"
	"time"

	"marketplace/internal/domain"
	"marketplace/internal/repository"
)

type PromoService interface {
	GetByCode(ctx context.Context, code string) (*domain.PromoCode, error)
	Validate(ctx context.Context, code string) (*domain.PromoCode, error)
}

type promoService struct {
	repo repository.PromoRepository
}

func NewPromoService(repo repository.PromoRepository) PromoService {
	return &promoService{repo: repo}
}

func (s *promoService) GetByCode(ctx context.Context, code string) (*domain.PromoCode, error) {
	p, err := s.repo.GetByCode(ctx, code)
	if err != nil {
		return nil, domain.ErrNotFound
	}
	return p, nil
}

func (s *promoService) Validate(ctx context.Context, code string) (*domain.PromoCode, error) {
	promo, err := s.GetByCode(ctx, code)
	if err != nil {
		return nil, err
	}

	if err := promo.CanBeUsed(time.Now()); err != nil {
		return nil, err
	}

	return promo, nil
}
