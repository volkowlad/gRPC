package jwt

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/volkowlad/gRPC/internal/myerr"
	"time"

	"github.com/volkowlad/gRPC/internal/config"
	"github.com/volkowlad/gRPC/internal/domain"

	"github.com/golang-jwt/jwt/v5"
)

func NewAccessToken(cfg config.Token, user domain.Users) (string, error) {
	if cfg.JWTSecret == "" {
		return "", errors.New("jwt secret is required")
	}

	if cfg.AccessTTL > time.Hour {
		return "", errors.New("jwt token ttl is more than 3600")
	}

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = user.Username
	claims["uuid"] = user.ID
	claims["exp"] = time.Now().Add(cfg.AccessTTL).Unix()

	tokenString, err := token.SignedString([]byte(cfg.JWTSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func NewRefreshToken(cfg config.Token, id uuid.UUID) (string, domain.RefreshToken, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	tokenID := uuid.New()
	ttl := time.Now().Add(cfg.RefreshTTL).Unix()
	createdAt := time.Now()

	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = id
	claims["hash"] = tokenID
	claims["expire_at"] = ttl
	claims["created_at"] = createdAt

	var refreshToken domain.RefreshToken

	tokenString, err := token.SignedString([]byte(cfg.JWTSecret))
	if err != nil {
		return "", refreshToken, err
	}

	refreshToken.ID = id
	refreshToken.Hash = tokenID
	refreshToken.ExpireAt = ttl
	refreshToken.CreatedAt = createdAt

	return tokenString, refreshToken, nil
}

func ParseRefreshToken(tokenString, secret string) (uuid.UUID, error) {
	token, err := jwt.NewParser().ParseWithClaims(tokenString, &domain.RefreshToken{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})
	if err != nil {
		return uuid.Nil, err
	}

	if claims, ok := token.Claims.(*domain.RefreshToken); ok && token.Valid {
		if time.Now().Unix() > claims.ExpireAt {
			return uuid.Nil, errors.New("token expired")
		}
		return claims.Hash, nil
	}

	return uuid.Nil, myerr.ErrInvalidToken
}
