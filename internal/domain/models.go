package domain

import (
	"time"
)

// User представляет пользователя системы
type User struct {
	ID        int64     `json:"id" db:"id"`
	Username  string    `json:"username" db:"username"`
	Password  string    `json:"-" db:"password_hash"`
	Balance   int64     `json:"balance" db:"balance"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// Merch представляет товар в магазине
type Merch struct {
	ID          int64     `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	Price       int64     `json:"price" db:"price"`
	Quantity    int       `json:"quantity" db:"quantity"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// Purchase представляет покупку мерча пользователем
type Purchase struct {
	ID        int64     `json:"id" db:"id"`
	UserID    int64     `json:"user_id" db:"user_id"`
	MerchID   int64     `json:"merch_id" db:"merch_id"`
	Quantity  int       `json:"quantity" db:"quantity"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// PurchaseResponse представляет покупку мерча с дополнительной информацией
type PurchaseResponse struct {
	ID         int64     `json:"id" db:"id"`
	UserID     int64     `json:"user_id" db:"user_id"`
	MerchID    int64     `json:"merch_id" db:"merch_id"`
	MerchName  string    `json:"merch_name" db:"merch_name"`
	MerchPrice int64     `json:"merch_price" db:"merch_price"`
	Quantity   int       `json:"quantity" db:"quantity"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}

// Transaction представляет операцию с монетами
type Transaction struct {
	ID          int64     `json:"id" db:"id"`
	FromUserID  int64     `json:"from_user_id" db:"from_user_id"`
	ToUserID    int64     `json:"to_user_id" db:"to_user_id"`
	Amount      int64     `json:"amount" db:"amount"`
	Description *string   `json:"description,omitempty" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

// UserInfoResponse представляет полную информацию о пользователе
type UserInfoResponse struct {
	Balance      int64               `json:"balance"`
	Transactions []*Transaction      `json:"transactions"`
	Purchases    []*PurchaseResponse `json:"inventory"`
}
