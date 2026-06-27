package middleware

import (
	"fmt"
	"net/http"
	"seconda/internal/service"
	"seconda/pkg/config"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr, err := c.Cookie("access_token")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Access token missing"})
			c.Abort()
			return
		}

		token, err := jwt.ParseWithClaims(tokenStr, &service.Claims{}, func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			return []byte(viper.GetString(config.SecretKey)), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired access token"})
			c.Abort()
			return
		}

		//Устанавливает user_id в контекст, если токен валиден
		if claims, ok := token.Claims.(*service.Claims); ok {
			c.Set("user_id", claims.UserId)
		}

		c.Next()
	}
}
