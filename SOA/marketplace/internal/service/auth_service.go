package service

import (
	"context"
	"log/slog"

	"marketplace/internal/domain"
	"marketplace/internal/repository"

	"github.com/google/uuid"
)

type PasswordHasher interface {
	Hash(password string) (string, error)
	Compare(hash, password string) error
}

type AuthService interface {
	Register(ctx context.Context, email, password string) (*domain.User, error)
	Login(ctx context.Context, email, password string) (*domain.User, error)
}

type authService struct {
	userRepo repository.UserRepository
	hasher   PasswordHasher
	logger	 *slog.Logger
}

func NewAuthService(userRepo repository.UserRepository, hasher PasswordHasher) AuthService {
	return &authService{
		userRepo: userRepo,
		hasher:   hasher,
	}
}

func (s *authService) Register(ctx context.Context, email, password string) (*domain.User, error) {
	s.logger.Info("registering user")

	if email == "" || password == "" {
		return nil, domain.ErrInvalidInput
	}

	_, err := s.userRepo.GetByEmail(ctx, email)
	if err == nil {
		s.logger.Error("failed", "error", err)
		return nil, domain.ErrEmailAlreadyExists
	}

	hash, err := s.hasher.Hash(password)
	if err != nil {
		s.logger.Error("failed", "error", err)
		return nil, domain.ErrInternal
	}

	user := &domain.User{
		ID:           uuid.NewString(),
		Email:        email,
		PasswordHash: hash,
		Role:         domain.RoleUser,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		s.logger.Error("failed", "error", err)
		return nil, err
	}

	return user, nil
}

func (s *authService) Login(ctx context.Context, email, password string) (*domain.User, error) {
	s.logger.Info("logging user")

	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		s.logger.Error("failed", "error", err)
		return nil, domain.ErrUnauthorized
	}

	if err := s.hasher.Compare(user.PasswordHash, password); err != nil {
		s.logger.Error("failed", "error", err)
		return nil, domain.ErrUnauthorized
	}

	return user, nil
}
