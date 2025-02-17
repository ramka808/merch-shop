package service

import (
	"context"

	"github.com/avito/internal/domain"
)

type TransactionService struct {
	transactionRepo domain.TransactionRepository
	userRepo        domain.UserRepository
}

func NewTransactionService(
	transactionRepo domain.TransactionRepository,
	userRepo domain.UserRepository,
) *TransactionService {
	return &TransactionService{
		transactionRepo: transactionRepo,
		userRepo:        userRepo,
	}
}

func (s *TransactionService) Transfer(ctx context.Context, fromUserID, toUserID int64, amount int64, description string) error {
	if amount <= 0 {
		return domain.ErrInvalidAmount
	}

	// Проверяем существование пользователей
	fromUser, err := s.userRepo.GetByID(ctx, fromUserID)
	if err != nil {
		return domain.ErrUserNotFound
	}

	_, err = s.userRepo.GetByID(ctx, toUserID)
	if err != nil {
		return domain.ErrUserNotFound
	}

	// Проверяем достаточность средств
	if fromUser.Balance < amount {
		return domain.ErrInsufficientFunds
	}

	// Создаем транзакцию
	transaction := &domain.Transaction{
		FromUserID:  fromUserID,
		ToUserID:    toUserID,
		Amount:      amount,
		Description: &description,
	}

	// Выполняем перевод денег и создаем запись о транзакции
	err = s.transactionRepo.TransferMoney(ctx, fromUserID, toUserID, amount)
	if err != nil {
		return domain.ErrTransactionFailed
	}

	err = s.transactionRepo.Create(ctx, transaction)
	if err != nil {
		// В случае ошибки пытаемся откатить перевод
		_ = s.transactionRepo.TransferMoney(ctx, toUserID, fromUserID, amount)
		return domain.ErrTransactionFailed
	}

	return nil
}

func (s *TransactionService) GetUserTransactions(ctx context.Context, userID int64) ([]*domain.Transaction, error) {
	// Проверяем существование пользователя
	_, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, domain.ErrUserNotFound
	}

	return s.transactionRepo.GetByUserID(ctx, userID)
}
