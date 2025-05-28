package grpc

import (
	"github.com/volkowlad/gRPC/internal/grpc/server"

	"go.uber.org/zap"
)

type GRPCServer struct {
	GRPC *server.Server
}

func NewGRPCServer(log *zap.SugaredLogger, port int) *GRPCServer {

	gRPC := server.NewGrpc(log, port)
	return &GRPCServer{
		GRPC: gRPC,
	}
}
