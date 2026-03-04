package repository

import (
	"context"
	"database/sql"
)

type TxManager interface {
	WithinTx(ctx context.Context, fn func(ctx context.Context) error) error
}

type txKey struct{}

type txManager struct {
	db *sql.DB
}

func NewTxManager(db *sql.DB) TxManager {
	return &txManager{db: db}
}

func (m *txManager) WithinTx(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	ctxWithTx := context.WithValue(ctx, txKey{}, tx)

	if err := fn(ctxWithTx); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}

func GetTx(ctx context.Context) *sql.Tx {
	tx, _ := ctx.Value(txKey{}).(*sql.Tx)
	return tx
}
