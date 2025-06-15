package repos

import (
	"context"
	"fmt"

	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/volkowlad/gRPC/internal/config"
	"github.com/volkowlad/gRPC/internal/myerr"

	"github.com/golang-migrate/migrate/v4"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func setupTestContainer(t *testing.T) (string, int) {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "postgres:15-alpine",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "test",
			"POSTGRES_PASSWORD": "test",
			"POSTGRES_DB":       "test_db",
		},
		WaitingFor: wait.ForAll(
			wait.ForLog("database system is ready to accept connections"),
			wait.ForListeningPort("5432/tcp"),
			wait.ForExec([]string{"pg_isready", "-U", "test"}),
		).WithStartupTimeout(10 * time.Second),
	}

	pgContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	require.NoError(t, err, "Failed to start container")

	t.Cleanup(func() {
		if err := pgContainer.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate container: %s", err)
		}
	})

	port, err := pgContainer.MappedPort(ctx, "5432")
	require.NoError(t, err, "Failed to get container port")

	p, err := strconv.Atoi(port.Port())
	if err != nil {
		t.Fatalf("failed to convert port: %s", err)
	}

	return fmt.Sprintf("postgres://test:test@localhost:%d/test_db?sslmode=disable", p), p
}

func applyMigrations(t *testing.T, connStr string) {
	// Получаем абсолютный путь к папке с миграциями
	wd, err := os.Getwd()
	require.NoError(t, err)
	migrationsPath := filepath.Join(wd, "../../migrations")

	if _, err := os.Stat(migrationsPath); os.IsNotExist(err) {
		t.Fatalf("Migrations directory does not exist: %s", migrationsPath)
	}

	t.Logf("Applying migrations from: %s", migrationsPath)
	t.Logf("Connection string: %s", connStr)

	m, err := migrate.New(
		fmt.Sprintf("file://%s", migrationsPath),
		connStr,
	)
	require.NoError(t, err, "Failed to create migrate instance")

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		require.NoError(t, err, "Failed to apply migrations")
	}

	//t.Cleanup(func() {
	//	if err := m.Down(); err != nil && err != migrate.ErrNoChange {
	//		t.Logf("Warning: failed to rollback migrations: %v", err)
	//	}
	//})
}

func setupTestRepo(t *testing.T) *Repository {
	connStr, port := setupTestContainer(t)
	applyMigrations(t, connStr)

	t.Logf("%d", port)

	cfg := config.PostgreSQL{
		User:                "test",
		Password:            "test",
		Host:                "localhost",
		Name:                "test_db",
		Port:                port, // Будет переопределено в connection string
		SSLMode:             "disable",
		PoolMaxConns:        5,
		PoolMaxConnLifetime: time.Hour,
		PoolMaxConnIdleTime: time.Minute,
	}

	// Подменяем порт из connection string
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	repo, err := NewPostgres(ctx, cfg)
	require.NoError(t, err, "Failed to create repository")

	t.Cleanup(func() {
		repo.pool.Close()
	})

	return repo
}

func TestUserSaver(t *testing.T) {
	repo := setupTestRepo(t)
	ctx := context.Background()

	// Тест создания пользователя
	t.Run("Create and get user", func(t *testing.T) {
		username := "testuser"
		passHash := []byte("hashed_password")

		// Проверяем, что пользователя нет
		exists, err := repo.UserByUsername(ctx, username)
		require.NoError(t, err)
		require.False(t, exists)

		// Создаем пользователя
		err = repo.UserSaver(ctx, username, passHash)
		require.NoError(t, err)

		// Проверяем, что пользователь появился
		exists, err = repo.UserByUsername(ctx, username)
		t.Logf("%t", exists)
		require.ErrorIs(t, err, myerr.ErrAlreadyExists)
		require.True(t, exists)

		// Получаем данные для входа
		user, err := repo.Login(ctx, username)
		require.NoError(t, err)
		require.Equal(t, username, user.Username)
		require.Equal(t, passHash, user.PassHash)
	})
}
