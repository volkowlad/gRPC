package validate

import (
	"github.com/volkowlad/gRPC/protos/gen"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	UsernameError = status.Error(codes.InvalidArgument, "username is required")
	PasswordError = status.Error(codes.InvalidArgument, "password is required")
)

func ValidateLogin(req *gen.LoginRequest) error {
	if req.GetUsername() == "" {
		return UsernameError
	}

	if req.GetPassword() == "" {
		return PasswordError
	}

	return nil
}

func ValidateRegister(req *gen.RegisterRequest) error {
	if req.GetUsername() == "" {
		return UsernameError
	}

	if req.GetPassword() == "" {
		return PasswordError
	}

	return nil
}
