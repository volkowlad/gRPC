package domain

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"time"
)

type Users struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	PassHash []byte    `json:"password"`
}

type RefreshToken struct {
	ID        uuid.UUID `json:"id"`
	Hash      uuid.UUID `json:"hash"`
	ExpireAt  int64     `json:"expire_at"`
	CreatedAt time.Time `json:"created_at"`
	jwt.RegisteredClaims
}
