package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/avito/internal/domain"
	"github.com/jmoiron/sqlx"
)

type TransactionRepository struct {
	*Repository
}

func NewTransactionRepository(repo *Repository) *TransactionRepository {
	return &TransactionRepository{Repository: repo}
}

func (r *TransactionRepository) Create(ctx context.Context, transaction *domain.Transaction) error {
	query := `
		INSERT INTO transactions (from_user_id, to_user_id, amount, description)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at`

	err := r.db.QueryRowxContext(ctx, query,
		transaction.FromUserID,
		transaction.ToUserID,
		transaction.Amount,
		transaction.Description,
	).Scan(&transaction.ID, &transaction.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to create transaction: %w", err)
	}

	return nil
}

func (r *TransactionRepository) GetByID(ctx context.Context, id int64) (*domain.Transaction, error) {
	transaction := &domain.Transaction{}

	query := `
		SELECT id, from_user_id, to_user_id, amount, description, created_at
		FROM transactions
		WHERE id = $1`

	err := r.db.GetContext(ctx, transaction, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("transaction not found")
		}
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}

	return transaction, nil
}

func (r *TransactionRepository) GetByUserID(ctx context.Context, userID int64) ([]*domain.Transaction, error) {
	var transactions []*domain.Transaction

	query := `
		SELECT id, from_user_id, to_user_id, amount, description, created_at
		FROM transactions
		WHERE from_user_id = $1 OR to_user_id = $1
		ORDER BY created_at DESC`

	err := r.db.SelectContext(ctx, &transactions, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user transactions: %w", err)
	}

	return transactions, nil
}

func (r *TransactionRepository) TransferMoney(ctx context.Context, fromUserID, toUserID int64, amount int64) error {
	return r.Transaction(ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		var senderBalance int64
		err := tx.QueryRowContext(ctx, `
			SELECT balance 
			FROM users 
			WHERE id = $1 
			FOR UPDATE`,
			fromUserID,
		).Scan(&senderBalance)
		if err != nil {
			return fmt.Errorf("failed to get sender balance: %w", err)
		}

		if senderBalance < amount {
			return fmt.Errorf("insufficient funds: available %d, required %d", senderBalance, amount)
		}

		_, err = tx.ExecContext(ctx, `
			UPDATE users 
			SET balance = balance - $1, 
				updated_at = CURRENT_TIMESTAMP
			WHERE id = $2`,
			amount, fromUserID,
		)
		if err != nil {
			return fmt.Errorf("failed to update sender balance: %w", err)
		}

		_, err = tx.ExecContext(ctx, `
			UPDATE users 
			SET balance = balance + $1, 
				updated_at = CURRENT_TIMESTAMP
			WHERE id = $2`,
			amount, toUserID,
		)
		if err != nil {
			return fmt.Errorf("failed to update recipient balance: %w", err)
		}

		return nil
	})
}
