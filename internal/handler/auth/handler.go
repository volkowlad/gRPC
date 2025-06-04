package auth

import (
	"context"

	service "github.com/volkowlad/gRPC/internal/service/auth"
	"github.com/volkowlad/gRPC/internal/validate"
	"github.com/volkowlad/gRPC/protos/gen"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

type Service interface {
	Login(ctx context.Context, username, password string) (string, string, error)
	Register(ctx context.Context, username, password string) (string, error)
	CheckToken(ctx context.Context, token string) (string, string, error)
}

type Server struct {
	gen.UnimplementedAuthServiceServer
	service *service.Service
}

func NewHandlers(g *grpc.Server, service *service.Service) {
	gen.RegisterAuthServiceServer(g, &Server{
		service: service,
	})
}

func (s *Server) Login(ctx context.Context, req *gen.LoginRequest) (*gen.LoginResponse, error) {
	if err := validate.ValidateLogin(req); err != nil {
		return nil, errors.Wrap(err, "invalid login")
	}

	tokenAccess, tokenRefresh, err := s.service.Login(ctx, req.GetUsername(), req.GetPassword())
	if err != nil {
		return nil, errors.Wrap(err, "failed to login")
	}

	return &gen.LoginResponse{
		Access:  tokenAccess,
		Refresh: tokenRefresh,
	}, nil
}

func (s *Server) Register(ctx context.Context, req *gen.RegisterRequest) (*gen.RegisterResponse, error) {
	if err := validate.ValidateRegister(req); err != nil {
		return nil, errors.Wrap(err, "invalid register")
	}

	message, err := s.service.Register(ctx, req.GetUsername(), req.GetPassword())
	if err != nil {
		return nil, errors.Wrap(err, "failed to register")
	}

	return &gen.RegisterResponse{
		Message: message,
	}, nil
}

func (s *Server) CheckToken(ctx context.Context, req *gen.CheckTokenRequest) (*gen.CheckTokenResponse, error) {
	if err := validate.ValidateToken(req); err != nil {
		return nil, errors.Wrap(err, "invalid token")
	}

	access, refresh, err := s.service.CheckToken(ctx, req.GetToken())
	if err != nil {
		return nil, errors.Wrap(err, "failed to check token")
	}

	return &gen.CheckTokenResponse{
		Access:  access,
		Refresh: refresh,
	}, nil
}
