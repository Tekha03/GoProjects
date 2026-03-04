package repository

import (
	"context"
	"database/sql"

	"marketplace/internal/domain"
)

type ProductRepository interface {
	GetByID(ctx context.Context, id string) (*domain.Product, error)
	UpdateStock(ctx context.Context, id string, stock int) error
	List(ctx context.Context, limit, offset int) ([]domain.Product, error)
	Create(ctx context.Context, p *domain.Product) error
}

type productRepo struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) ProductRepository {
	return &productRepo{db: db}
}

func (r *productRepo) executor(ctx context.Context) execer {
	if tx := GetTx(ctx); tx != nil {
		return tx
	}
	return r.db
}

type execer interface {
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

func (r *productRepo) GetByID(ctx context.Context, id string) (*domain.Product, error) {
	query := `
		SELECT id, name, description, price, stock, created_at, updated_at
		FROM products WHERE id = $1
	`
	row := r.executor(ctx).QueryRowContext(ctx, query, id)

	var p domain.Product
	err := row.Scan(
		&p.ID,
		&p.Name,
		&p.Description,
		&p.Price,
		&p.Stock,
		&p.CreatedAt,
		&p.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *productRepo) UpdateStock(ctx context.Context, id string, stock int) error {
	query := `UPDATE products SET stock = $1, updated_at = now() WHERE id = $2`
	_, err := r.executor(ctx).ExecContext(ctx, query, stock, id)
	return err
}

func (r *productRepo) List(ctx context.Context, limit, offset int) ([]domain.Product, error) {
	query := `
		SELECT id, name, description, price, stock, created_at, updated_at
		FROM products
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.executor(ctx).QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []domain.Product
	for rows.Next() {
		var p domain.Product
		if err := rows.Scan(
			&p.ID,
			&p.Name,
			&p.Description,
			&p.Price,
			&p.Stock,
			&p.CreatedAt,
			&p.UpdatedAt,
		); err != nil {
			return nil, err
		}
		result = append(result, p)
	}

	return result, nil
}

func (r *productRepo) Create(ctx context.Context, p *domain.Product) error {
	query := `
		INSERT INTO products (id, name, description, price, stock)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.executor(ctx).ExecContext(
		ctx,
		query,
		p.ID,
		p.Name,
		p.Description,
		p.Price,
		p.Stock,
	)
	return err
}
