package handler

import (
	"github.com/avito/internal/delivery/http/middleware"
	"github.com/avito/internal/domain"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	userService        domain.UserService
	merchService       domain.MerchService
	transactionService domain.TransactionService
}

func NewHandler(services *domain.Services) *Handler {
	return &Handler{
		userService:        services.User,
		merchService:       services.Merch,
		transactionService: services.Transaction,
	}
}

func (h *Handler) Init(router *gin.Engine, tokenSecret string) {
	authMiddleware := middleware.AuthMiddleware(tokenSecret)

	v1 := router.Group("/api")
	{
		userHandler := NewUserHandler(h.userService)
		v1.POST("/auth", userHandler.Auth)

		userGroup := v1.Group("/user")
		userGroup.Use(authMiddleware)
		{
			userHandler := NewUserHandler(h.userService)
			userGroup.GET("/balance", userHandler.GetBalance)
		}

		merchGroup := v1.Group("/merch")
		{
			merchHandler := NewMerchHandler(h.merchService)
			merchGroup.GET("", merchHandler.GetList)
			merchGroup.GET("/:id", merchHandler.GetByID)

			protected := merchGroup.Group("/")
			protected.Use(authMiddleware)
			{
				protected.POST("/buy", merchHandler.Buy)
				protected.GET("/my", merchHandler.GetUserPurchases)
			}
		}

		transactionGroup := v1.Group("/transactions")
		transactionGroup.Use(authMiddleware)
		{
			transactionHandler := NewTransactionHandler(h.transactionService)
			transactionGroup.POST("/transfer", transactionHandler.Transfer)
			transactionGroup.GET("/history", transactionHandler.GetHistory)
		}
	}
}
