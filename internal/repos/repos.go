package repos

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/volkowlad/gRPC/internal/config"
	"github.com/volkowlad/gRPC/internal/domain"
	"github.com/volkowlad/gRPC/internal/myerr"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

const (
	insertUserQuery      = `INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id;`
	selectUserByUsername = `SELECT username FROM users WHERE username = $1;`
	selectLogin          = `SELECT id, username, password FROM users WHERE username = $1;`
	insertRefresh        = `INSERT INTO refresh_tokens (user_id, token_hash, expires_at, created_at) VALUES ($1, $2, $3, $4)`
	selectRefresh        = `SELECT EXISTS(
							SELECT 1 FROM refresh_tokens 
							WHERE token_hash = $1 
		)`
	selectUserIDRefresh = `SELECT u.id, u.username
						   FROM refresh_tokens rt
						   JOIN users u ON rt.user_id = u.id
						   WHERE rt.token_hash = $1`
	updateRefresh = `UPDATE refresh_tokens SET token_hash = $1, expires_at = $2, created_at = $3 WHERE user_id = $4`
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

func (r *Repository) UserSaver(ctx context.Context, username string, passHash []byte) (uuid.UUID, error) {
	var id uuid.UUID
	err := r.pool.QueryRow(ctx, insertUserQuery, username, passHash).Scan(&id)
	if err != nil {
		return uuid.Nil, errors.Wrap(err, "failed to insert user")
	}

	return id, nil
}

func (r *Repository) UserByUsername(ctx context.Context, username string) (bool, error) {
	var newName string

	err := r.pool.QueryRow(ctx, selectUserByUsername, username).Scan(&newName)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}

		return false, errors.Wrap(err, "failed to query user by username")
	}

	return true, myerr.ErrAlreadyExists
}

func (r *Repository) Login(ctx context.Context, username string) (domain.Users, error) {
	var users domain.Users

	err := r.pool.QueryRow(ctx, selectLogin, username).Scan(&users.ID, &users.Username, &users.PassHash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return users, myerr.ErrNotFound
		}

		return users, errors.Wrap(err, "failed to query user")
	}

	return users, nil
}

func (r *Repository) RefreshTokenSaver(ctx context.Context, refreshToken domain.RefreshToken) error {
	_, err := r.pool.Exec(ctx, insertRefresh, refreshToken.ID, refreshToken.Hash, refreshToken.ExpireAt, refreshToken.CreatedAt)
	if err != nil {
		return errors.Wrap(err, "failed to insert user")
	}

	return nil
}

func (r *Repository) RefreshTokenCheck(ctx context.Context, tokenID uuid.UUID) (bool, error) {
	var exist bool

	err := r.pool.QueryRow(ctx, selectRefresh, tokenID).Scan(&exist)
	if err != nil {
		return false, errors.Wrap(err, "failed to query refresh token")
	}

	return exist, nil
}

func (r *Repository) UserByID(ctx context.Context, id uuid.UUID) (domain.Users, error) {
	var users domain.Users
	err := r.pool.QueryRow(ctx, selectUserIDRefresh, id).Scan(&users.ID, &users.Username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return users, myerr.ErrNotFound
		}

		return domain.Users{}, errors.Wrap(err, "failed to query user")
	}

	return users, nil
}

func (r *Repository) RefreshUpdate(ctx context.Context, token domain.RefreshToken) error {
	_, err := r.pool.Exec(ctx, updateRefresh, token.Hash, token.ExpireAt, token.CreatedAt, token.ID)
	if err != nil {
		return errors.Wrap(err, "failed to insert user")
	}

	return nil
}
