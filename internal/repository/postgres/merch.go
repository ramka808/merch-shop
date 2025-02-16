package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/avito/internal/domain"
)

type MerchRepository struct {
	*Repository
}

func NewMerchRepository(repo *Repository) *MerchRepository {
	return &MerchRepository{Repository: repo}
}

func (r *MerchRepository) Create(ctx context.Context, merch *domain.Merch) error {
	query := `
		INSERT INTO merch (name, description, price)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at`

	err := r.db.QueryRowxContext(ctx, query,
		merch.Name,
		merch.Description,
		merch.Price,
	).Scan(&merch.ID, &merch.CreatedAt, &merch.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create merch: %w", err)
	}

	return nil
}

func (r *MerchRepository) GetByID(ctx context.Context, id int64) (*domain.Merch, error) {
	merch := &domain.Merch{}

	query := `
		SELECT id, name, description, price, created_at, updated_at
		FROM merch
		WHERE id = $1`

	err := r.db.GetContext(ctx, merch, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("merch not found")
		}
		return nil, fmt.Errorf("failed to get merch: %w", err)
	}

	return merch, nil
}

func (r *MerchRepository) List(ctx context.Context, limit, offset int) ([]*domain.Merch, error) {
	var items []*domain.Merch

	query := `
		SELECT id, name, description, price, created_at, updated_at
		FROM merch
		ORDER BY id
		LIMIT $1 OFFSET $2`

	err := r.db.SelectContext(ctx, &items, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list merch: %w", err)
	}

	return items, nil
}

func (r *MerchRepository) UpdateQuantity(ctx context.Context, merchID int64, quantity int) error {
	// В данной реализации количество не отслеживается, так как по условию
	// "Предполагается, что в магазине бесконечный запас каждого вида мерча"
	return nil
}
