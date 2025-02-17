package http

import (
	"github.com/gin-gonic/gin"
)

// Response представляет структуру успешного ответа
type Response struct {
	Message string      `json:"description,omitempty"`
	Data    interface{} `json:"schema,omitempty"`
}

// NewResponse отправляет успешный ответ
func NewResponse(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, Response{
		Message: message,
		Data:    data,
	})
}

// NoContent отправляет ответ без содержимого
func NoContent(c *gin.Context) {
	c.Status(204)
}

// Created отправляет ответ о успешном создании ресурса
func Created(c *gin.Context, message string, data interface{}) {
	NewResponse(c, 201, message, data)
}

// OK отправляет успешный ответ
func OK(c *gin.Context, message string, data interface{}) {
	NewResponse(c, 200, message, data)
}

type errorResponse struct {
	Message string `json:"message"`
	Code    string `json:"code"`
}

func newErrorResponse(c *gin.Context, statusCode int, message, code string) {
	c.AbortWithStatusJSON(statusCode, errorResponse{
		Message: message,
		Code:    code,
	})
}
