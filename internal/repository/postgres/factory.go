package postgres

import "github.com/avito/internal/domain"

// Repositories содержит все репозитории
type Repositories struct {
	User        domain.UserRepository
	Merch       domain.MerchRepository
	Purchase    domain.PurchaseRepository
	Transaction domain.TransactionRepository
}

// NewRepositories создает новый экземпляр всех репозиториев
func NewRepositories(dsn string) (*Repositories, error) {
	repo, err := New(dsn)
	if err != nil {
		return nil, err
	}

	return &Repositories{
		User:        NewUserRepository(repo),
		Merch:       NewMerchRepository(repo),
		Purchase:    NewPurchaseRepository(repo),
		Transaction: NewTransactionRepository(repo),
	}, nil
}
