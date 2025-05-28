package validate

import (
	"github.com/volkowlad/gRPC/protos/gen"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func UsernameLogin(req *gen.LoginRequest) error {
	if req.GetUsername() == "" {
		return status.Error(codes.InvalidArgument, "username is required")
	}

	return nil
}

func PasswordLogin(req *gen.LoginRequest) error {
	if req.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "password is required")
	}

	return nil
}

func UsernameRegister(req *gen.RegisterRequest) error {
	if req.GetUsername() == "" {
		return status.Error(codes.InvalidArgument, "username is required")
	}

	return nil
}

func PasswordRegister(req *gen.RegisterRequest) error {
	if req.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "password is required")
	}

	return nil
}
