package jwt

import (
	"github.com/google/uuid"
	"time"

	"github.com/volkowlad/gRPC/internal/config"
	"github.com/volkowlad/gRPC/internal/domain"

	"github.com/golang-jwt/jwt/v5"
)

func NewAccessToken(cfg config.Token, user domain.Users) (string, error) {
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

func NewRefreshToken(cfg config.Token, id uuid.UUID) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = id
	claims["uuid"] = uuid.New()
	claims["exp"] = time.Now().Add(cfg.RefreshTTL).Unix()
	claims["created_at"] = time.Now().Unix()

	tokenString, err := token.SignedString([]byte(cfg.JWTSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

//func ParseRefreshToken(tokenString string) (uuid.UUID, error) {
//	token, _, err := jwt.NewParser().ParseUnverified(tokenString, &jwt.MapClaims{})
//	if err != nil {
//		return uuid.Nil, err
//	}
//
//	if claims, ok := token.Claims.(*jwt.MapClaims); ok && token.Valid {
//		return claims.
//	}
//}
