package repos

import (
	"context"
	"fmt"
	"github.com/volkowlad/gRPC/internal/domain"
	"github.com/volkowlad/gRPC/internal/myerr"

	"github.com/volkowlad/gRPC/internal/config"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

const (
	insertUserQuery      = `INSERT INTO users (username, password) VALUES ($1, $2);`
	selectUserByUsername = `SELECT username FROM users WHERE username = $1;`
	selectLogin          = `SELECT id, username, password FROM users WHERE username = $1;`
)

type Repository struct {
	pool *pgxpool.Pool
}

// NewRepository - создание нового экземпляра репозитория с подключением к PostgreSQL
func NewPostgres(ctx context.Context, cfg config.PostgreSQL) (*Repository, error) {
	// Формируем строку подключения
	connString := fmt.Sprintf(
		`user=%s password=%s host=%s port=%d dbname=%s sslmode=%s 
        pool_max_conns=%d pool_max_conn_lifetime=%s pool_max_conn_idle_time=%s`,
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
		cfg.SSLMode,
		cfg.PoolMaxConns,
		cfg.PoolMaxConnLifetime.String(),
		cfg.PoolMaxConnIdleTime.String(),
	)

	// Парсим конфигурацию подключения
	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse PostgreSQL config")
	}

	// Оптимизация выполнения запросов (кеширование запросов)
	config.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeCacheDescribe

	// Создаём пул соединений с базой данных
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create PostgreSQL connection pool")
	}

	return &Repository{pool}, nil
}

func (r *Repository) UserSaver(ctx context.Context, username string, passHash []byte) error {
	_, err := r.pool.Exec(ctx, insertUserQuery, username, passHash)
	if err != nil {
		return errors.Wrap(err, "failed to insert user")
	}

	return nil
}

func (r *Repository) UserByUsername(ctx context.Context, username string) (bool, error) {
	_, err := r.pool.Query(ctx, selectUserByUsername, username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, err
		}
	}

	return true, myerr.ErrAlreadyExists
}

func (r *Repository) Login(ctx context.Context, username string) (domain.Users, error) {
	var users domain.Users

	err := r.pool.QueryRow(ctx, selectLogin, username).Scan(&users.ID, &users.Username, &users.PassHash)
	if err != nil {
		if err == pgx.ErrNoRows {
			return users, myerr.ErrNotFound
		}

		return users, err
	}

	return users, nil
}
