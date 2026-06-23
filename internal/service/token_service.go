package service

import (
	"context"
	"errors"
	"seconda/cmd/service"
	"seconda/internal/dto"
	"seconda/pkg/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/spf13/viper"
)

const ISS = "token_service"

// TODO: удалить, если не будет использоваться
var (
	TokenInvalidErr = errors.New("token is invalid")
	TokenExpiredErr = errors.New("token has expired or revoked")
)

type Claims struct {
	UserId string `json:"user_id"`
	jwt.RegisteredClaims
}

// CreateTokens создает пару Access и Refresh, сохраняя JTI в Redis
func CreateTokens(ctx context.Context, redisDec service.RedisDecorator, userId string) (dto.AccessToken, dto.RefreshToken, error) {
	atStr, err := createAccessToken(userId)
	if err != nil {
		return dto.AccessToken{}, dto.RefreshToken{}, err
	}

	err = redisDec.Client.Set(ctx, userId, atStr, viper.GetDuration(config.AccessTokenLT)).Err()
	if err != nil {
		return dto.AccessToken{}, dto.RefreshToken{}, err
	}

	//TODO: сейчас не используется. Если будет время - написать метод для refresh токена
	rtStr, jti, err := createRefreshToken(userId)
	if err != nil {
		return dto.AccessToken{}, dto.RefreshToken{}, err
	}

	err = redisDec.Client.Set(ctx, jti, userId, viper.GetDuration(config.RefreshTokenLT)).Err()
	if err != nil {
		return dto.AccessToken{}, dto.RefreshToken{}, err
	}

	return dto.AccessToken{
			UserId: userId,
			Token:  atStr,
			Lt:     viper.GetDuration(config.AccessTokenLT),
		}, dto.RefreshToken{
			UserId: userId,
			Token:  rtStr,
			Lt:     viper.GetDuration(config.RefreshTokenLT),
			Jti:    jti,
		}, nil
}

func createAccessToken(userId string) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS512, Claims{
		UserId: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    ISS,
			Subject:   userId,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(viper.GetDuration(config.AccessTokenLT))),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}).SignedString([]byte(viper.GetString(config.SecretKey)))
}

// createRefreshToken возвращает сам токен и его внутренний uuid (JTI)
func createRefreshToken(userId string) (string, string, error) {
	jti := uuid.New().String()
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		UserId: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        jti,
			Issuer:    ISS,
			Subject:   userId,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(viper.GetDuration(config.RefreshTokenLT))),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}).SignedString([]byte(viper.GetString(config.SecretKey)))

	return token, jti, err
}

// Logout удаляет токен из Redis
func Logout(ctx context.Context, redisDec service.RedisDecorator, refreshTokenStr string) error {
	token, _ := jwt.ParseWithClaims(refreshTokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(viper.GetString(config.SecretKey)), nil
	})

	if claims, ok := token.Claims.(*Claims); ok {
		return redisDec.Client.Del(ctx, claims.ID).Err()
	}
	return nil
}
