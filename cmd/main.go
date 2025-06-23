package main

import (
	"context"
	"github.com/volkowlad/gRPC/internal/grpc"
	"github.com/volkowlad/gRPC/internal/repos"
	service "github.com/volkowlad/gRPC/internal/service/auth"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"

	"github.com/volkowlad/gRPC/internal/config"
	customLog "github.com/volkowlad/gRPC/internal/logger"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal(errors.Wrap(err, "Error loading .env file"))
	}

	// Загружаем конфигурацию из переменных окружения
	var cfg config.AppConfig
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatal(errors.Wrap(err, "failed to load configuration"))
	}

	lg, err := customLog.NewLogger(cfg.LogLevel)
	if err != nil {
		log.Fatal(errors.Wrap(err, "error initializing logger"))
	}

	ctx := context.Background()

	repository, err := repos.NewPostgres(ctx, cfg.Postgres)
	if err != nil {
		lg.Fatal(errors.Wrap(err, "error initializing postgres"))
	}

	services := service.NewService(cfg.Token, repository, lg)

	app := grpc.NewGRPCServer(lg, services, cfg.GRPC.ListenAddress)

	go app.GRPC.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	app.GRPC.Stop()

	lg.Info("shutting down")
}
