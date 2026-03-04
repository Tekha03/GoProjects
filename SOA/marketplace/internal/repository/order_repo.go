package repository

import (
	"context"
	"database/sql"

	"marketplace/internal/domain"
)

type OrderRepository interface {
	Create(ctx context.Context, order *domain.Order) error
	AddItems(ctx context.Context, orderID string, items []domain.OrderItem) error
	GetByID(ctx context.Context, id string) (*domain.Order, error)
	UpdateStatus(ctx context.Context, id string, status domain.OrderStatus) error
}

type orderRepo struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) OrderRepository {
	return &orderRepo{db: db}
}

func (r *orderRepo) executor(ctx context.Context) execer {
	if tx := GetTx(ctx); tx != nil {
		return tx
	}
	return r.db
}

func (r *orderRepo) Create(ctx context.Context, order *domain.Order) error {
	query := `
		INSERT INTO orders (id, user_id, status, total_amount, promo_id)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.executor(ctx).ExecContext(
		ctx,
		query,
		order.ID,
		order.UserID,
		order.Status,
		order.TotalAmount,
		order.PromoID,
	)
	return err
}

func (r *orderRepo) AddItems(ctx context.Context, orderID string, items []domain.OrderItem) error {
	query := `
		INSERT INTO order_items (id, order_id, product_id, quantity, price)
		VALUES (uuid_generate_v4(), $1, $2, $3, $4)
	`
	for _, item := range items {
		if _, err := r.executor(ctx).ExecContext(
			ctx,
			query,
			orderID,
			item.ProductID,
			item.Quantity,
			item.Price,
		); err != nil {
			return err
		}
	}
	return nil
}

func (r *orderRepo) GetByID(ctx context.Context, id string) (*domain.Order, error) {
	query := `
		SELECT id, user_id, status, total_amount, promo_id, created_at, updated_at
		FROM orders WHERE id = $1
	`
	row := r.db.QueryRowContext(ctx, query, id)

	var o domain.Order
	if err := row.Scan(
		&o.ID,
		&o.UserID,
		&o.Status,
		&o.TotalAmount,
		&o.PromoID,
		&o.CreatedAt,
		&o.UpdatedAt,
	); err != nil {
		return nil, err
	}

	itemsQuery := `
		SELECT product_id, quantity, price
		FROM order_items WHERE order_id = $1
	`
	rows, err := r.db.QueryContext(ctx, itemsQuery, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item domain.OrderItem
		if err := rows.Scan(&item.ProductID, &item.Quantity, &item.Price); err != nil {
			return nil, err
		}
		o.Items = append(o.Items, item)
	}

	return &o, nil
}

func (r *orderRepo) UpdateStatus(ctx context.Context, id string, status domain.OrderStatus) error {
	query := `
		UPDATE orders
		SET status = $1, updated_at = now()
		WHERE id = $2
	`
	_, err := r.executor(ctx).ExecContext(ctx, query, status, id)
	return err
}
