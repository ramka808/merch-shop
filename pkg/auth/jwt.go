package auth

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// TokenManager предоставляет методы для работы с JWT токенами
type TokenManager struct {
	signingKey string
}

// NewTokenManager создает новый менеджер токенов
func NewTokenManager(signingKey string) (*TokenManager, error) {
	if signingKey == "" {
		return nil, errors.New("empty signing key")
	}

	return &TokenManager{signingKey: signingKey}, nil
}

// NewJWT создает новый JWT токен
func (m *TokenManager) NewJWT(userID int64, username string, ttl time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"exp":      time.Now().Add(ttl).Unix(),
	})

	return token.SignedString([]byte(m.signingKey))
}

// Parse проверяет токен на валидность и возвращает ID пользователя
func (m *TokenManager) Parse(accessToken string) (int64, error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(m.signingKey), nil
	})

	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("invalid token claims")
	}

	userID, ok := claims["user_id"].(float64)
	if !ok {
		return 0, errors.New("invalid user id")
	}

	return int64(userID), nil
}
