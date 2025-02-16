package handler

import (
	"net/http"

	httpDelivery "github.com/avito/internal/delivery/http"
	"github.com/avito/internal/delivery/http/middleware"
	"github.com/avito/internal/domain"
	"github.com/gin-gonic/gin"
)

type userHandler struct {
	userService domain.UserService
}

func NewUserHandler(userService domain.UserService) *userHandler {
	return &userHandler{
		userService: userService,
	}
}

type signUpInput struct {
	Username string `json:"username" binding:"required,min=3,max=32"`
	Password string `json:"password" binding:"required,min=6,max=32"`
}

type signInInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type authInput struct {
	Username string `json:"username" binding:"required,min=3,max=32"`
	Password string `json:"password" binding:"required,min=6,max=32"`
}

func (h *userHandler) SignUp(c *gin.Context) {
	var input signUpInput
	if err := c.ShouldBindJSON(&input); err != nil {
		httpDelivery.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid_input")
		return
	}

	user, err := h.userService.Register(c.Request.Context(), input.Username, input.Password)
	if err != nil {
		if err == domain.ErrUserAlreadyExists {
			httpDelivery.NewErrorResponse(c, http.StatusConflict, err.Error(), "user_exists")
			return
		}
		httpDelivery.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "internal_error")
		return
	}

	httpDelivery.OK(c, "user created", gin.H{
		"id":       user.ID,
		"username": user.Username,
		"balance":  user.Balance,
	})
}

func (h *userHandler) SignIn(c *gin.Context) {
	var input signInInput
	if err := c.ShouldBindJSON(&input); err != nil {
		httpDelivery.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid_input")
		return
	}

	token, err := h.userService.Login(c.Request.Context(), input.Username, input.Password)
	if err != nil {
		if err == domain.ErrInvalidCredentials {
			httpDelivery.NewErrorResponse(c, http.StatusUnauthorized, err.Error(), "invalid_credentials")
			return
		}
		httpDelivery.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "internal_error")
		return
	}

	httpDelivery.OK(c, "success", gin.H{
		"token": token,
	})
}

func (h *userHandler) GetBalance(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		httpDelivery.NewErrorResponse(c, http.StatusUnauthorized, err.Error(), "unauthorized")
		return
	}

	balance, err := h.userService.GetBalance(c.Request.Context(), userID)
	if err != nil {
		if err == domain.ErrUserNotFound {
			httpDelivery.NewErrorResponse(c, http.StatusNotFound, err.Error(), "user_not_found")
			return
		}
		httpDelivery.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "internal_error")
		return
	}

	httpDelivery.OK(c, "success", gin.H{
		"balance": balance,
	})
}

func (h *userHandler) Auth(c *gin.Context) {
	var input authInput
	if err := c.ShouldBindJSON(&input); err != nil {
		httpDelivery.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid_input")
		return
	}

	token, err := h.userService.Auth(c.Request.Context(), input.Username, input.Password)
	if err != nil {
		if err == domain.ErrInvalidCredentials {
			httpDelivery.NewErrorResponse(c, http.StatusUnauthorized, err.Error(), "invalid_credentials")
			return
		}
		httpDelivery.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "internal_error")
		return
	}

	httpDelivery.OK(c, "success", gin.H{
		"token": token,
	})
}
