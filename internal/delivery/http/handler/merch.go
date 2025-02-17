package handler

import (
	"net/http"
	"strconv"

	httpDelivery "github.com/avito/internal/delivery/http"
	"github.com/avito/internal/delivery/http/middleware"
	"github.com/avito/internal/domain"
	"github.com/gin-gonic/gin"
)

type merchHandler struct {
	merchService domain.MerchService
}

func NewMerchHandler(merchService domain.MerchService) *merchHandler {
	return &merchHandler{
		merchService: merchService,
	}
}

type buyMerchInput struct {
	MerchID  int64 `json:"merch_id" binding:"required"`
	Quantity int   `json:"quantity" binding:"required,min=1"`
}

func (h *merchHandler) GetList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	items, err := h.merchService.List(c.Request.Context(), page, pageSize)
	if err != nil {
		httpDelivery.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "internal_error")
		return
	}

	httpDelivery.OK(c, "Успешный ответ", items)
}

func (h *merchHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		httpDelivery.NewErrorResponse(c, http.StatusBadRequest, "invalid id", "invalid_input")
		return
	}

	merch, err := h.merchService.GetByID(c.Request.Context(), id)
	if err != nil {
		if err == domain.ErrMerchNotFound {
			httpDelivery.NewErrorResponse(c, http.StatusNotFound, err.Error(), "merch_not_found")
			return
		}
		httpDelivery.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "internal_error")
		return
	}

	httpDelivery.OK(c, "Успешный ответ", merch)
}

func (h *merchHandler) Buy(c *gin.Context) {
	var input buyMerchInput
	if err := c.ShouldBindJSON(&input); err != nil {
		httpDelivery.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid_input")
		return
	}

	userID, err := middleware.GetUserID(c)
	if err != nil {
		httpDelivery.NewErrorResponse(c, http.StatusUnauthorized, err.Error(), "unauthorized")
		return
	}

	err = h.merchService.Buy(c.Request.Context(), userID, input.MerchID, input.Quantity)
	if err != nil {
		switch err {
		case domain.ErrMerchNotFound:
			httpDelivery.NewErrorResponse(c, http.StatusNotFound, err.Error(), "merch_not_found")
		case domain.ErrInsufficientFunds:
			httpDelivery.NewErrorResponse(c, http.StatusPaymentRequired, err.Error(), "insufficient_funds")
		case domain.ErrInvalidQuantity:
			httpDelivery.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid_quantity")
		default:
			httpDelivery.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "internal_error")
		}
		return
	}

	httpDelivery.OK(c, "purchase successful", nil)
}

func (h *merchHandler) GetUserPurchases(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		httpDelivery.NewErrorResponse(c, http.StatusUnauthorized, err.Error(), "unauthorized")
		return
	}

	purchases, err := h.merchService.GetUserPurchases(c.Request.Context(), userID)
	if err != nil {
		if err == domain.ErrUserNotFound {
			httpDelivery.NewErrorResponse(c, http.StatusNotFound, err.Error(), "user_not_found")
			return
		}
		httpDelivery.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "internal_error")
		return
	}

	httpDelivery.OK(c, "Успешный ответ", purchases)
}
