package service

import (
	"context"
	"fmt"

	"github.com/avito/internal/domain"
)

type MerchService struct {
	merchRepo    domain.MerchRepository
	purchaseRepo domain.PurchaseRepository
	userRepo     domain.UserRepository
}

func NewMerchService(
	merchRepo domain.MerchRepository,
	purchaseRepo domain.PurchaseRepository,
	userRepo domain.UserRepository,
) *MerchService {
	return &MerchService{
		merchRepo:    merchRepo,
		purchaseRepo: purchaseRepo,
		userRepo:     userRepo,
	}
}

func (s *MerchService) List(ctx context.Context, page, pageSize int) ([]*domain.Merch, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	return s.merchRepo.List(ctx, pageSize, offset)
}

func (s *MerchService) GetByID(ctx context.Context, id int64) (*domain.Merch, error) {
	merch, err := s.merchRepo.GetByID(ctx, id)
	if err != nil {
		return nil, domain.ErrMerchNotFound
	}
	return merch, nil
}

func (s *MerchService) Buy(ctx context.Context, userID, merchID int64, quantity int) error {
	if quantity <= 0 {
		return domain.ErrInvalidQuantity
	}

	// Получаем информацию о мерче
	merch, err := s.merchRepo.GetByID(ctx, merchID)
	if err != nil {
		return domain.ErrMerchNotFound
	}

	// Проверяем баланс пользователя
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return domain.ErrUserNotFound
	}

	totalCost := merch.Price * int64(quantity)
	if user.Balance < totalCost {
		return domain.ErrInsufficientFunds
	}

	// Создаем запись о покупке
	purchase := &domain.Purchase{
		UserID:   userID,
		MerchID:  merchID,
		Quantity: quantity,
	}

	// Обновляем баланс пользователя и создаем покупку в одной транзакции
	err = s.userRepo.UpdateBalance(ctx, userID, -totalCost)
	if err != nil {
		return fmt.Errorf("failed to update balance: %w", err)
	}

	err = s.purchaseRepo.Create(ctx, purchase)
	if err != nil {
		// В случае ошибки возвращаем деньги
		_ = s.userRepo.UpdateBalance(ctx, userID, totalCost)
		return domain.ErrTransactionFailed
	}

	return nil
}

func (s *MerchService) GetUserPurchases(ctx context.Context, userID int64) ([]*domain.PurchaseResponse, error) {
	return s.purchaseRepo.GetByUserID(ctx, userID)
}
