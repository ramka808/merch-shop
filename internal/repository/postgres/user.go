package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/avito/internal/domain"
	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	*Repository
}

func NewUserRepository(repo *Repository) *UserRepository {
	return &UserRepository{Repository: repo}
}

func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO users (username, password_hash, balance)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at`

	err := r.db.QueryRowxContext(ctx, query,
		user.Username,
		user.Password,
		user.Balance,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (r *UserRepository) GetByID(ctx context.Context, id int64) (*domain.User, error) {
	user := &domain.User{}

	query := `
		SELECT id, username, password_hash, balance, created_at, updated_at
		FROM users
		WHERE id = $1`

	err := r.db.GetContext(ctx, user, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	user := &domain.User{}

	query := `
		SELECT id, username, password_hash, balance, created_at, updated_at
		FROM users
		WHERE username = $1`

	err := r.db.GetContext(ctx, user, query, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

func (r *UserRepository) UpdateBalance(ctx context.Context, userID int64, amount int64) error {
	return r.Transaction(ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		var newBalance int64
		err := tx.QueryRowContext(ctx, `
			UPDATE users
			SET balance = balance + $1, updated_at = CURRENT_TIMESTAMP
			WHERE id = $2
			RETURNING balance`, amount, userID).Scan(&newBalance)
		if err != nil {
			return fmt.Errorf("failed to update balance: %w", err)
		}

		if newBalance < 0 {
			return fmt.Errorf("insufficient funds")
		}

		return nil
	})
}
