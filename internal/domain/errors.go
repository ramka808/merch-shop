package domain
// package main

import "errors"

var (
	// ErrUserNotFound возвращается, когда пользователь не найден
	ErrUserNotFound = errors.New("user not found")
	

	// ErrUserAlreadyExists возвращается при попытке создать существующего пользователя
	ErrUserAlreadyExists = errors.New("user already exists")

	// ErrInvalidCredentials возвращается при неверных учетных данных
	ErrInvalidCredentials = errors.New("invalid credentials")

	// ErrInsufficientFunds возвращается при недостаточном балансе
	ErrInsufficientFunds = errors.New("insufficient funds")

	// ErrMerchNotFound возвращается, когда товар не найден
	ErrMerchNotFound = errors.New("merch not found")

	// ErrInvalidAmount возвращается при некорректной сумме операции
	ErrInvalidAmount = errors.New("invalid amount")

	// ErrInvalidQuantity возвращается при некорректном количестве товара
	ErrInvalidQuantity = errors.New("invalid quantity")

	// ErrTransactionFailed возвращается при ошибке проведения транзакции
	ErrTransactionFailed = errors.New("transaction failed")
)

