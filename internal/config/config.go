package config

import "time"

type AppConfig struct {
	LogLevel string `envconfig:"LOG_LEVEL" default:"info"`
	GRPC     GRPC
	Postgres PostgreSQL
	Token    Token
}

type GRPC struct {
	ListenAddress int           `envconfig:"PORT" required:"true"`
	WriteTimeout  time.Duration `envconfig:"WRITE_TIMEOUT" required:"true"`
}

type PostgreSQL struct {
	Host                string        `envconfig:"DB_HOST" required:"true"`
	Port                int           `envconfig:"DB_PORT" required:"true"`
	Name                string        `envconfig:"DB_NAME" required:"true"`
	User                string        `envconfig:"DB_USER" required:"true"`
	Password            string        `envconfig:"DB_PASSWORD" required:"true"`
	SSLMode             string        `envconfig:"DB_SSL_MODE" default:"disable"`
	PoolMaxConns        int           `envconfig:"DB_POOL_MAX_CONNS" default:"5"`
	PoolMaxConnLifetime time.Duration `envconfig:"DB_POOL_MAX_CONN_LIFETIME" default:"180s"`
	PoolMaxConnIdleTime time.Duration `envconfig:"DB_POOL_MAX_CONN_IDLE_TIME" default:"100s"`
}

type Token struct {
	JWTSecret  string        `envconfig:"JWT_SECRET" required:"true"`
	AccessTTL  time.Duration `envconfig:"ACCESS_TTL" default:"1h"`
	RefreshTTL time.Duration `envconfig:"REFRESH_TTL" default:"744h"`
}
