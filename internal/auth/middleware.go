package auth

import (
	"net/http"
	"painaway_test/internal/config"
	"painaway_test/internal/response"
	"painaway_test/internal/utils"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"go.uber.org/zap"
)

func AuthMiddleware(cfg *config.JWTConfig, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.NewErrorResponse(c, http.StatusUnauthorized, "authorization header missing", logger)
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Token" {
			response.NewErrorResponse(c, http.StatusUnauthorized, "invalid authorization header format", logger)
			c.Abort()
			return
		}

		tokenString := parts[1]

		token, err := jwt.ParseWithClaims(tokenString, &utils.Claims{}, func(t *jwt.Token) (interface{}, error) {
			return []byte(cfg.SecretKey), nil
		})
		if err != nil || !token.Valid {
			response.NewErrorResponse(c, http.StatusUnauthorized, "invalid or expired token", logger)
			c.Abort()
			return
		}

		claims, ok := token.Claims.(*utils.Claims)
		if !ok {
			response.NewErrorResponse(c, http.StatusUnauthorized, "invalid token claims", logger)
			c.Abort()
			return
		}

		c.Set("userID", uint(claims.UserID))
		c.Set("groups", string(claims.Groups))
		c.Next()
	}
}
