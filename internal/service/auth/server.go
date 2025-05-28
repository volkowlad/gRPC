package auth

import (
	"context"
	"github.com/volkowlad/gRPC/protos/gen"
	"google.golang.org/grpc"
)

type Server struct {
	gen.UnimplementedAuthServiceServer
}

func RegisterServer(g *grpc.Server) {
	gen.RegisterAuthServiceServer(g, &Server{})
}

func (s *Server) Login(ctx context.Context, req *gen.LoginRequest) (*gen.LoginResponse, error) {
	return &gen.LoginResponse{
		Token: "123",
	}, nil
}

func (s *Server) Register(ctx context.Context, req *gen.RegisterRequest) (*gen.RegisterResponse, error) {
	return &gen.RegisterResponse{
		Message: "success",
	}, nil
}
