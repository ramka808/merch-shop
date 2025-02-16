package http

import (
	"net/http"

	"github.com/avito/internal/domain"
	"github.com/gin-gonic/gin"
)

// ErrorResponse представляет структуру ответа с ошибкой
type ErrorResponse struct {
	Message string `json:"message"`
	Code    string `json:"code"`
}

// NewErrorResponse отправляет ответ с ошибкой
func NewErrorResponse(c *gin.Context, statusCode int, message, code string) {
	c.AbortWithStatusJSON(statusCode, ErrorResponse{
		Message: message,
		Code:    code,
	})
}

// HandleError обрабатывает ошибки и отправляет соответствующий HTTP-ответ
func HandleError(c *gin.Context, err error) {
	switch err {
	case domain.ErrUserNotFound:
		NewErrorResponse(c, http.StatusNotFound, err.Error(), "user_not_found")
	case domain.ErrUserAlreadyExists:
		NewErrorResponse(c, http.StatusConflict, err.Error(), "user_exists")
	case domain.ErrInvalidCredentials:
		NewErrorResponse(c, http.StatusUnauthorized, err.Error(), "invalid_credentials")
	case domain.ErrInsufficientFunds:
		NewErrorResponse(c, http.StatusPaymentRequired, err.Error(), "insufficient_funds")
	case domain.ErrMerchNotFound:
		NewErrorResponse(c, http.StatusNotFound, err.Error(), "merch_not_found")
	case domain.ErrInvalidAmount:
		NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid_amount")
	case domain.ErrInvalidQuantity:
		NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid_quantity")
	case domain.ErrTransactionFailed:
		NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "transaction_failed")
	default:
		NewErrorResponse(c, http.StatusInternalServerError, "internal server error", "internal_error")
	}
}
