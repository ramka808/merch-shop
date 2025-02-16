package domain

import "context"

// Repositories содержит все репозитории приложения
type Repositories struct {
	User        UserRepository
	Merch       MerchRepository
	Purchase    PurchaseRepository
	Transaction TransactionRepository
}

// UserRepository определяет методы для работы с пользователями
type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id int64) (*User, error)
	GetByUsername(ctx context.Context, username string) (*User, error)
	UpdateBalance(ctx context.Context, userID int64, amount int64) error
}

// MerchRepository определяет методы для работы с мерчем
type MerchRepository interface {
	Create(ctx context.Context, merch *Merch) error
	GetByID(ctx context.Context, id int64) (*Merch, error)
	List(ctx context.Context, limit, offset int) ([]*Merch, error)
	UpdateQuantity(ctx context.Context, merchID int64, quantity int) error
}

// PurchaseRepository определяет методы для работы с покупками
type PurchaseRepository interface {
	Create(ctx context.Context, purchase *Purchase) error
	GetByUserID(ctx context.Context, userID int64) ([]*Purchase, error)
	GetByID(ctx context.Context, id int64) (*Purchase, error)
}

// TransactionRepository определяет методы для работы с транзакциями
type TransactionRepository interface {
	Create(ctx context.Context, transaction *Transaction) error
	GetByUserID(ctx context.Context, userID int64) ([]*Transaction, error)
	GetByID(ctx context.Context, id int64) (*Transaction, error)
	// TransferMoney выполняет перевод денег между пользователями в транзакции
	TransferMoney(ctx context.Context, fromUserID, toUserID int64, amount int64) error
}
