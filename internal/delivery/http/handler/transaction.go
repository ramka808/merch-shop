package handler

import (
	"net/http"

	httpDelivery "github.com/avito/internal/delivery/http"
	"github.com/avito/internal/delivery/http/middleware"
	"github.com/avito/internal/domain"
	"github.com/gin-gonic/gin"
)

type transactionHandler struct {
	transactionService domain.TransactionService
}

func NewTransactionHandler(transactionService domain.TransactionService) *transactionHandler {
	return &transactionHandler{
		transactionService: transactionService,
	}
}

type transferInput struct {
	ToUserID    int64  `json:"to_user_id" binding:"required"`
	Amount      int64  `json:"amount" binding:"required,min=1"`
	Description string `json:"description"`
}

func (h *transactionHandler) Transfer(c *gin.Context) {
	var input transferInput
	if err := c.ShouldBindJSON(&input); err != nil {
		httpDelivery.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid_input")
		return
	}

	fromUserID, err := middleware.GetUserID(c)
	if err != nil {
		httpDelivery.NewErrorResponse(c, http.StatusUnauthorized, err.Error(), "unauthorized")
		return
	}

	if fromUserID == input.ToUserID {
		httpDelivery.NewErrorResponse(c, http.StatusBadRequest, "cannot transfer to yourself", "invalid_recipient")
		return
	}

	err = h.transactionService.Transfer(c.Request.Context(), fromUserID, input.ToUserID, input.Amount, input.Description)
	if err != nil {
		switch err {
		case domain.ErrUserNotFound:
			httpDelivery.NewErrorResponse(c, http.StatusNotFound, err.Error(), "user_not_found")
		case domain.ErrInsufficientFunds:
			httpDelivery.NewErrorResponse(c, http.StatusPaymentRequired, err.Error(), "insufficient_funds")
		case domain.ErrInvalidAmount:
			httpDelivery.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid_amount")
		default:
			httpDelivery.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "internal_error")
		}
		return
	}

	httpDelivery.OK(c, "transfer successful", nil)
}

func (h *transactionHandler) GetHistory(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		httpDelivery.NewErrorResponse(c, http.StatusUnauthorized, err.Error(), "unauthorized")
		return
	}

	transactions, err := h.transactionService.GetUserTransactions(c.Request.Context(), userID)
	if err != nil {
		if err == domain.ErrUserNotFound {
			httpDelivery.NewErrorResponse(c, http.StatusNotFound, err.Error(), "user_not_found")
			return
		}
		httpDelivery.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "internal_error")
		return
	}

	httpDelivery.OK(c, "Успешный ответ", transactions)
}
