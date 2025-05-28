package server

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/volkowlad/gRPC/internal/service/auth"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
)

type Server struct {
	log  *zap.SugaredLogger
	gRPC *grpc.Server
	port int
}

func NewGrpc(log *zap.SugaredLogger, port int) *Server {
	gRPCServer := grpc.NewServer()

	auth.RegisterServer(gRPCServer)

	return &Server{
		log:  log,
		gRPC: gRPCServer,
		port: port,
	}
}

func (s *Server) Run() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		return errors.Wrap(err, "failed to run gRPC server")
	}

	s.log.Infof("gRPC server listening on port %d", s.port)

	if err := s.gRPC.Serve(lis); err != nil {
		return errors.Wrap(err, "failed to run gRPC server")
	}

	return nil
}

func (s *Server) Stop() {
	s.gRPC.GracefulStop()

	s.log.Infof("gRPC server stopped")
}

func (s *Server) MustRun() {
	if err := s.Run(); err != nil {
		s.log.Fatal(err)
	}
}
