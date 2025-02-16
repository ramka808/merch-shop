package service

import "github.com/avito/internal/domain"

// NewServices создает новый экземпляр всех сервисов
func NewServices(deps domain.Deps) *domain.Services {
	userService := NewUserService(deps.Repos.User, deps.TokenSecret)
	merchService := NewMerchService(deps.Repos.Merch, deps.Repos.Purchase, deps.Repos.User)
	transactionService := NewTransactionService(deps.Repos.Transaction, deps.Repos.User)

	return &domain.Services{
		User:        userService,
		Merch:       merchService,
		Transaction: transactionService,
	}
}
