package postgres

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// Repository представляет собой обертку над подключением к базе данных
type Repository struct {
	db *sqlx.DB
}

// New создает новый экземпляр Repository
func New(dsn string) (*Repository, error) {
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Устанавливаем максимальное количество открытых соединений
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)

	return &Repository{db: db}, nil
}

// Close закрывает соединение с базой данных
func (r *Repository) Close() error {
	return r.db.Close()
}

// Transaction выполняет функцию в транзакции
func (r *Repository) Transaction(ctx context.Context, fn func(ctx context.Context, tx *sqlx.Tx) error) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
	}()

	if err := fn(ctx, tx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx rollback error: %v (original error: %w)", rbErr, err)
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
