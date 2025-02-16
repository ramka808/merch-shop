package app

import (
	"github.com/avito/internal/config"
	"github.com/avito/internal/delivery/http/handler"
	"github.com/avito/internal/domain"
	"github.com/avito/internal/repository/postgres"
	"github.com/avito/internal/service"
	"github.com/gin-gonic/gin"
)

type App struct {
	router *gin.Engine
	cfg    *config.Config
}

func NewApp() (*App, error) {
	// Загружаем конфигурацию
	cfg, err := config.New()
	if err != nil {
		return nil, err
	}

	// Инициализируем репозитории
	repos, err := postgres.NewRepositories(cfg.Postgres.DSN())
	if err != nil {
		return nil, err
	}

	// Инициализируем сервисы
	deps := domain.Deps{
		Repos: &domain.Repositories{
			User:        repos.User,
			Merch:       repos.Merch,
			Purchase:    repos.Purchase,
			Transaction: repos.Transaction,
		},
		TokenSecret: cfg.JWT.SecretKey,
	}
	services := service.NewServices(deps)

	// Инициализируем handler
	h := handler.NewHandler(services)

	// Создаем новый роутер
	router := gin.Default()

	// Инициализируем маршруты через handler
	h.Init(router, cfg.JWT.SecretKey)

	return &App{
		router: router,
		cfg:    cfg,
	}, nil
}

func (a *App) Run(addr string) error {
	return a.router.Run(addr)
}
