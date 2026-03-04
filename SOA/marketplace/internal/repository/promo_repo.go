package repository

import (
	"context"
	"database/sql"

	"marketplace/internal/domain"
)

type PromoRepository interface {
	GetByCode(ctx context.Context, code string) (*domain.PromoCode, error)
	IncrementUsage(ctx context.Context, id string) error
}

type promoRepo struct {
	db *sql.DB
}

func NewPromoRepository(db *sql.DB) PromoRepository {
	return &promoRepo{db: db}
}

func (r *promoRepo) GetByCode(ctx context.Context, code string) (*domain.PromoCode, error) {
	query := `
		SELECT id, code, discount_percent, expires_at,
		       usage_limit, used_count, active, created_at
		FROM promo_codes WHERE code = $1
	`
	row := r.db.QueryRowContext(ctx, query, code)

	var p domain.PromoCode
	err := row.Scan(
		&p.ID,
		&p.Code,
		&p.DiscountPercent,
		&p.ExpiresAt,
		&p.UsageLimit,
		&p.UsedCount,
		&p.Active,
		&p.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *promoRepo) IncrementUsage(ctx context.Context, id string) error {
	query := `
		UPDATE promo_codes
		SET used_count = used_count + 1
		WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
