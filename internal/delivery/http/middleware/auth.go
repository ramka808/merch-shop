package middleware

import (
	"net/http"
	"strings"

	"github.com/avito/internal/domain"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "userId"
)

func AuthMiddleware(tokenSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader(authorizationHeader)
		if header == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "empty auth header",
				"code":    "unauthorized",
			})
			return
		}

		headerParts := strings.Split(header, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "invalid auth header",
				"code":    "unauthorized",
			})
			return
		}

		if len(headerParts[1]) == 0 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "token is empty",
				"code":    "unauthorized",
			})
			return
		}

		token, err := jwt.Parse(headerParts[1], func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, domain.ErrInvalidCredentials
			}
			return []byte(tokenSecret), nil
		})

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "invalid token",
				"code":    "unauthorized",
			})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "invalid token claims",
				"code":    "unauthorized",
			})
			return
		}

		userID, ok := claims["user_id"].(float64)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "invalid user id",
				"code":    "unauthorized",
			})
			return
		}

		c.Set(userCtx, int64(userID))
		c.Next()
	}
}

// GetUserID получает ID пользователя из контекста
func GetUserID(c *gin.Context) (int64, error) {
	id, ok := c.Get(userCtx)
	if !ok {
		return 0, domain.ErrUserNotFound
	}

	idInt64, ok := id.(int64)
	if !ok {
		return 0, domain.ErrUserNotFound
	}

	return idInt64, nil
}
