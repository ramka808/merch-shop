package domain

import "context"

// UserService определяет методы для работы с пользователями
type UserService interface {
	Register(ctx context.Context, username, password string) (*User, error)
	Login(ctx context.Context, username, password string) (string, error) // Возвращает JWT токен
	GetByID(ctx context.Context, id int64) (*User, error)
	GetBalance(ctx context.Context, userID int64) (int64, error)
	Auth(ctx context.Context, username, password string) (string, error) // Объединенный метод для регистрации/входа
}

// MerchService определяет методы для работы с мерчем
type MerchService interface {
	List(ctx context.Context, page, pageSize int) ([]*Merch, error)
	GetByID(ctx context.Context, id int64) (*Merch, error)
	Buy(ctx context.Context, userID, merchID int64, quantity int) error
	GetUserPurchases(ctx context.Context, userID int64) ([]*PurchaseResponse, error)
}

// TransactionService определяет методы для работы с транзакциями
type TransactionService interface {
	Transfer(ctx context.Context, fromUserID, toUserID int64, amount int64, description string) error
	GetUserTransactions(ctx context.Context, userID int64) ([]*Transaction, error)
}

// Services объединяет все сервисы приложения
type Services struct {
	User        UserService
	Merch       MerchService
	Transaction TransactionService
}

// Deps содержит зависимости для сервисов
type Deps struct {
	Repos       *Repositories
	TokenSecret string
}
