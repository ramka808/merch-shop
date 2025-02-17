package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/avito/internal/domain"
)

type PurchaseRepository struct {
	*Repository
}

func NewPurchaseRepository(repo *Repository) *PurchaseRepository {
	return &PurchaseRepository{Repository: repo}
}

func (r *PurchaseRepository) Create(ctx context.Context, purchase *domain.Purchase) error {
	query := `
		INSERT INTO purchases (user_id, merch_id, quantity)
		VALUES ($1, $2, $3)
		RETURNING id, created_at`

	err := r.db.QueryRowxContext(ctx, query,
		purchase.UserID,
		purchase.MerchID,
		purchase.Quantity,
	).Scan(&purchase.ID, &purchase.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to create purchase: %w", err)
	}

	return nil
}

func (r *PurchaseRepository) GetByID(ctx context.Context, id int64) (*domain.Purchase, error) {
	purchase := &domain.Purchase{}

	query := `
		SELECT id, user_id, merch_id, quantity, created_at
		FROM purchases
		WHERE id = $1`

	err := r.db.GetContext(ctx, purchase, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("purchase not found")
		}
		return nil, fmt.Errorf("failed to get purchase: %w", err)
	}

	return purchase, nil
}

func (r *PurchaseRepository) GetByUserID(ctx context.Context, userID int64) ([]*domain.PurchaseResponse, error) {
	var purchases []*domain.PurchaseResponse

	query := `
		SELECT p.id, p.user_id, p.merch_id, p.quantity, p.created_at,
			   m.name as merch_name, m.price as merch_price
		FROM purchases p
		JOIN merch m ON p.merch_id = m.id
		WHERE p.user_id = $1
		ORDER BY p.created_at DESC`

	err := r.db.SelectContext(ctx, &purchases, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user purchases: %w", err)
	}

	return purchases, nil
}
