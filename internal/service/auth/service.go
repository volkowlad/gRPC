package auth

import (
	"context"

	"github.com/volkowlad/gRPC/internal/config"
	"github.com/volkowlad/gRPC/internal/jwt"
	"github.com/volkowlad/gRPC/internal/myerr"
	"github.com/volkowlad/gRPC/internal/repos"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type Repository interface {
	UserSaver(ctx context.Context, email string, passHash []byte) error
	Login(ctx context.Context, email, passHash string) error
	UserByUsername(ctx context.Context, username string) (bool, error)
}

type Service struct {
	repository *repos.Repository
	log        *zap.SugaredLogger
	cfg        config.Token
}

func NewService(cfg config.Token, repos *repos.Repository, log *zap.SugaredLogger) *Service {
	return &Service{
		repository: repos,
		log:        log,
		cfg:        cfg,
	}
}

func (s *Service) Login(ctx context.Context, username, password string) (string, error) {
	user, err := s.repository.Login(ctx, username)
	if err != nil {
		if errors.Is(err, myerr.ErrNotFound) {
			s.log.Errorf("login failed: %v", err)

			return "", errors.New("failed to login")
		}

		s.log.Errorf("login failed: %v", err)

		return "", errors.Wrap(err, "failed to login")
	}

	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		s.log.Errorf("login failed: %v", err)

		return "", errors.New("failed to login")
	}

	token, err := jwt.NewToken(s.cfg, user)
	if err != nil {
		s.log.Errorf("login failed: %v", err)

		return "", errors.Wrap(err, "failed to login")
	}

	return token, nil
}

func (s *Service) Register(ctx context.Context, username, password string) (string, error) {
	exist, err := s.repository.UserByUsername(ctx, username)
	if err != nil {
		s.log.Errorf("failed to register user: %v", err)

		return "", errors.Wrap(err, "failed to register")
	}
	if exist {
		s.log.Errorf("user %s already exists", username)

		return "", err
	}

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		s.log.Errorf("failed to hash password: %v", err)

		return "", errors.Wrap(err, "failed to register user")
	}

	err = s.repository.UserSaver(ctx, username, passHash)
	if err != nil {
		s.log.Errorf("failed to register user: %v", err)

		return "", errors.Wrap(err, "failed to register user")
	}

	s.log.Infof("user %s registered", username)

	return "done", nil
}
