package myerr

import "github.com/pkg/errors"

var (
	ErrNotFound      = errors.New("not found")
	ErrAlreadyExists = errors.New("already exists")
	ErrInvalidToken  = errors.New("invalid token claims")
)
