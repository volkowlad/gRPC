package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"

	"gRPC/app/internal/config"
	customLog "gRPC/app/internal/logger"
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

	lg.Infof("%v", cfg)

	// TODO:database

	// TODO:app

	// TODO:server
}
