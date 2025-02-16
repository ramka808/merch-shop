package service

import (
	"context"
	"fmt"
	"time"

	"github.com/avito/internal/domain"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo        domain.UserRepository
	tokenSecret string
}

func NewUserService(repo domain.UserRepository, tokenSecret string) *UserService {
	return &UserService{
		repo:        repo,
		tokenSecret: tokenSecret,
	}
}

func (s *UserService) Register(ctx context.Context, username, password string) (*domain.User, error) {
	// Проверяем, не существует ли уже пользователь
	existingUser, err := s.repo.GetByUsername(ctx, username)
	if err == nil && existingUser != nil {
		return nil, domain.ErrUserAlreadyExists
	}

	// Хешируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &domain.User{
		Username: username,
		Password: string(hashedPassword),
		Balance:  1000, // Начальный баланс
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

func (s *UserService) Login(ctx context.Context, username, password string) (string, error) {
	user, err := s.repo.GetByUsername(ctx, username)
	if err != nil {
		return "", domain.ErrInvalidCredentials
	}

	// Проверяем пароль
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", domain.ErrInvalidCredentials
	}

	// Создаем JWT токен
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(s.tokenSecret))
	if err != nil {
		return "", fmt.Errorf("failed to create token: %w", err)
	}

	return tokenString, nil
}

func (s *UserService) GetByID(ctx context.Context, id int64) (*domain.User, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, domain.ErrUserNotFound
	}
	return user, nil
}

func (s *UserService) GetBalance(ctx context.Context, userID int64) (int64, error) {
	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return 0, domain.ErrUserNotFound
	}
	return user.Balance, nil
}

func (s *UserService) Auth(ctx context.Context, username, password string) (string, error) {
	// Пробуем найти пользователя
	existingUser, err := s.repo.GetByUsername(ctx, username)
	if err == nil && existingUser != nil {
		// Пользователь существует, проверяем пароль
		if err := bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(password)); err != nil {
			return "", domain.ErrInvalidCredentials
		}
	} else {
		// Пользователь не существует, создаем нового
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return "", fmt.Errorf("failed to hash password: %w", err)
		}

		existingUser = &domain.User{
			Username: username,
			Password: string(hashedPassword),
			Balance:  1000, // Начальный баланс
		}

		if err := s.repo.Create(ctx, existingUser); err != nil {
			return "", fmt.Errorf("failed to create user: %w", err)
		}
	}

	// Создаем JWT токен
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  existingUser.ID,
		"username": existingUser.Username,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(s.tokenSecret))
	if err != nil {
		return "", fmt.Errorf("failed to create token: %w", err)
	}

	return tokenString, nil
}
