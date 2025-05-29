package jwt

import (
	"github.com/volkowlad/gRPC/internal/config"
	"time"

	"github.com/volkowlad/gRPC/internal/domain"

	"github.com/golang-jwt/jwt/v5"
)

func NewToken(cfg config.Token, user domain.Users) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = user.Username
	claims["uuid"] = user.ID
	claims["exp"] = time.Now().Add(cfg.TTL).Unix()

	tokenString, err := token.SignedString([]byte(cfg.JWTSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
