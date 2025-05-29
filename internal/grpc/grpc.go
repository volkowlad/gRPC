package grpc

import (
	"github.com/volkowlad/gRPC/internal/grpc/server"
	service "github.com/volkowlad/gRPC/internal/service/auth"

	"go.uber.org/zap"
)

type GRPCServer struct {
	GRPC *server.Server
}

func NewGRPCServer(log *zap.SugaredLogger, service *service.Service, port int) *GRPCServer {

	gRPC := server.NewGrpc(log, service, port)
	return &GRPCServer{
		GRPC: gRPC,
	}
}
