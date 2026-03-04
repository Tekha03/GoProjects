package service

import (
	"context"

	"github.com/google/uuid"
	"marketplace/internal/domain"
	"marketplace/internal/repository"
)

type ProductService interface {
	Create(ctx context.Context, name, description string, price float64, stock int) (*domain.Product, error)
	GetByID(ctx context.Context, id string) (*domain.Product, error)
	List(ctx context.Context, limit, offset int) ([]domain.Product, error)
}

type productService struct {
	productRepo repository.ProductRepository
}

func NewProductService(productRepo repository.ProductRepository) ProductService {
	return &productService{productRepo: productRepo}
}

func (s *productService) Create(ctx context.Context, name, description string, price float64, stock int) (*domain.Product, error) {
	if price < 0 || stock < 0 {
		return nil, domain.ErrInvalidInput
	}

	p := &domain.Product{
		ID:          uuid.NewString(),
		Name:        name,
		Description: description,
		Price:       price,
		Stock:       stock,
	}

	if err := s.productRepo.Create(ctx, p); err != nil {
		return nil, err
	}

	return p, nil
}

func (s *productService) GetByID(ctx context.Context, id string) (*domain.Product, error) {
	p, err := s.productRepo.GetByID(ctx, id)
	if err != nil {
		return nil, domain.ErrNotFound
	}
	return p, nil
}

func (s *productService) List(ctx context.Context, limit, offset int) ([]domain.Product, error) {
	return s.productRepo.List(ctx, limit, offset)
}
