package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"seconda/cmd/base"
	"seconda/internal/model/user"
	"seconda/internal/request"
	"seconda/internal/service"
	"seconda/pkg/config"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

type AuthController struct {
	Controller base.Controller
}

func (ac *AuthController) HandleRequest() {
	e := ac.Controller.E

	group := e.Group("/api/v1")
	group.POST("/register", ac.Register)
	group.POST("/login", ac.Login)
}

func (ac *AuthController) Register(c *gin.Context) {
	var u user.User
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	if err := user.CreateUser(ac.Controller.DI.DBDecorator.GDB(), &u); err != nil &&
		(errors.Is(err, user.WithLoginAlreadyExistsErr) || errors.Is(err, user.WithPhoneAlreadyExistsErr)) {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user: " + err.Error()})
		return
	}

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 30*time.Minute)
	defer cancel()

	//TODO: подумать над передачей по ссылке
	at, rt, err := service.CreateTokens(ctx, *ac.Controller.DI.RedisDecorator, strconv.Itoa(u.Id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user token: " + err.Error()})
		return
	}

	c.SetCookie("access_token", at.Token, int(at.Lt.Seconds()), "/", "", false, true)
	//TODO: оставить что-то одно
	c.SetCookie("refresh_token", rt.Token, int(rt.Lt.Seconds()), "/auth/refresh", "", false, true)
	c.SetCookie("refresh_jti", rt.Jti, int(rt.Lt.Seconds()), "/auth/refresh-jti", "", false, true)

	c.JSON(http.StatusCreated, u)
}

func (ac *AuthController) Login(c *gin.Context) {
	tokenStr, err := c.Cookie("access_token")
	if err != nil && err.Error() != "http: named cookie not present" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "token not found"})
		return
	}

	if tokenStr == "" {
		var ulr request.UserLoginRequest
		if err := c.ShouldBindJSON(&ulr); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
			return
		}

		u, err := user.GetUserByLoginAndPass(ac.Controller.DI.DBDecorator.GDB(), ulr.Login, ulr.Password)
		if err != nil && errors.Is(err, user.NotFoundErr) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		} else if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}

		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, 30*time.Minute)
		defer cancel()

		stringCMD := ac.Controller.DI.RedisDecorator.Client.Get(ctx, strconv.Itoa(u.Id))
		if stringCMD.Err() != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}
		tokenStr = stringCMD.Val()
	}

	claims := &service.Claims{}
	t, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(viper.GetString(config.SecretKey)), nil
	})

	if err != nil || !t.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{
			"valid": false,
			"error": "invalid token",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"payload": gin.H{
			"valid":   true,
			"user_id": claims.UserId,
		},
	})
}
