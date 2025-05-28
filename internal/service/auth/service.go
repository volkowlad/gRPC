package auth

import (
	"context"
	"github.com/volkowlad/gRPC/internal/repos"
	"go.uber.org/zap"
)

type Service interface {
	Login(ctx context.Context, username, password string) (string, error)
	Register(ctx context.Context, username, password string) (string, error)
}

type service struct {
	repository repos.Repository
	log        *zap.SugaredLogger
}

func NewService(repos repos.Repository, log *zap.SugaredLogger) Service {
	return &service{
		repository: repos,
		log:        log,
	}
}

func (s *service) Login(ctx context.Context, username, password string) (string, error) {
	return "token", nil
}

func (s *service) Register(ctx context.Context, username, password string) (string, error) {
	return "done", nil
}
