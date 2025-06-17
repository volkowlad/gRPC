package auth

import (
	"context"
	"github.com/google/uuid"
	"github.com/volkowlad/gRPC/internal/domain"

	"github.com/volkowlad/gRPC/internal/config"
	"github.com/volkowlad/gRPC/internal/myerr"
	"github.com/volkowlad/gRPC/pkg/jwt"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

//go:generate mockgen -source=service.go -destination=mock/mock.go

type Repository interface {
	UserSaver(ctx context.Context, email string, passHash []byte) (uuid.UUID, error)
	Login(ctx context.Context, username string) (domain.Users, error)
	UserByUsername(ctx context.Context, username string) (bool, error)
	RefreshTokenSaver(ctx context.Context, refreshToken domain.RefreshToken) error
	RefreshTokenCheck(ctx context.Context, tokenID uuid.UUID) (bool, error)
	UserByID(ctx context.Context, id uuid.UUID) (domain.Users, error)
	RefreshUpdate(ctx context.Context, token domain.RefreshToken) error
}

type Service struct {
	repository Repository
	log        *zap.SugaredLogger
	cfg        config.Token
}

func NewService(cfg config.Token, repos Repository, log *zap.SugaredLogger) *Service {
	return &Service{
		repository: repos,
		log:        log,
		cfg:        cfg,
	}
}

func (s *Service) Login(ctx context.Context, username, password string) (string, string, error) {
	user, err := s.repository.Login(ctx, username)
	if err != nil {
		if errors.Is(err, myerr.ErrNotFound) {
			s.log.Errorf("login failed: %v", err)

			return "", "", errors.New("failed to login")
		}

		s.log.Errorf("login failed: %v", err)

		return "", "", errors.Wrap(err, "failed to login")
	}

	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		s.log.Errorf("login failed: %v", err)

		return "", "", errors.New("failed to login")
	}

	tokenAccess, err := jwt.NewAccessToken(s.cfg, user)
	if err != nil {
		s.log.Errorf("login failed: %v", err)

		return "", "", errors.Wrap(err, "failed to login")
	}

	tokenRefreshString, tokenRefresh, err := jwt.NewRefreshToken(s.cfg, user.ID)
	if err != nil {
		s.log.Errorf("login failed: %v", err)

		return "", "", errors.Wrap(err, "failed to login")
	}

	if err := s.repository.RefreshTokenSaver(ctx, tokenRefresh); err != nil {
		s.log.Errorf("login failed: %v", err)

		return "", "", errors.Wrap(err, "failed to login")
	}

	return tokenAccess, tokenRefreshString, nil
}

func (s *Service) Register(ctx context.Context, username, password string) (string, error) {
	exist, err := s.repository.UserByUsername(ctx, username)
	if err != nil {
		s.log.Errorf("failed to register user: %v", err)

		return "", errors.Wrap(err, "failed to register user")
	}
	if exist {
		s.log.Errorf("user %s already exists", username)

		return "", errors.Wrap(err, "failed to register user")
	}

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		s.log.Errorf("failed to hash password: %v", err)

		return "", errors.Wrap(err, "failed to register user")
	}

	id, err := s.repository.UserSaver(ctx, username, passHash)
	if err != nil {
		s.log.Errorf("failed to register user: %v", err)

		return "", errors.Wrap(err, "failed to register user")
	}

	strId := id.String()
	s.log.Infof("user %s registered", username)

	return strId, nil
}
func (s *Service) CheckToken(ctx context.Context, tokenString string) (string, string, error) {
	tokenID, err := jwt.ParseRefreshToken(tokenString, s.cfg.JWTSecret)
	if err != nil {
		s.log.Errorf("failed to parse refresh token: %v", err)

		return "", "", errors.Wrap(err, "failed to refresh token")
	}

	exist, err := s.repository.RefreshTokenCheck(ctx, tokenID)
	if err != nil {
		s.log.Errorf("failed to check refresh token: %v", err)

		return "", "", errors.Wrap(err, "failed to refresh token")
	}
	if !exist {
		s.log.Errorf("refresh token %s not found", tokenID)

		return "", "", errors.Wrap(myerr.ErrNotFound, "failed to refresh token")
	}

	user, err := s.repository.UserByID(ctx, tokenID)
	if err != nil {
		s.log.Errorf("failed to check user: %v", err)

		return "", "", errors.Wrap(err, "failed to refresh token")
	}

	tokenAccess, err := jwt.NewAccessToken(s.cfg, user)
	if err != nil {
		s.log.Errorf("login failed: %v", err)

		return "", "", errors.Wrap(err, "failed to refresh token")
	}

	tokenRefreshString, tokenRefresh, err := jwt.NewRefreshToken(s.cfg, user.ID)
	if err != nil {
		s.log.Errorf("login failed: %v", err)

		return "", "", errors.Wrap(err, "failed to refresh token")
	}

	err = s.repository.RefreshUpdate(ctx, tokenRefresh)
	if err != nil {
		s.log.Errorf("login failed: %v", err)

		return "", "", errors.Wrap(err, "failed to refresh token")
	}

	return tokenAccess, tokenRefreshString, nil
}
