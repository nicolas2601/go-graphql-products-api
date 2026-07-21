// Command main es el composition root: carga la configuracion, cablea las dependencias y
// levanta el servidor HTTP de la API de productos.
package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/nicolas2601/go-graphql-products-api/internal/config"
	"github.com/nicolas2601/go-graphql-products-api/internal/domain"
	"github.com/nicolas2601/go-graphql-products-api/internal/repository/memory"
	"github.com/nicolas2601/go-graphql-products-api/internal/repository/postgres"
	"github.com/nicolas2601/go-graphql-products-api/internal/server"
	"github.com/nicolas2601/go-graphql-products-api/internal/usecase"
)

const (
	readHeaderTimeout = 5 * time.Second
	readTimeout       = 15 * time.Second
	writeTimeout      = 15 * time.Second
	idleTimeout       = 60 * time.Second
	shutdownTimeout   = 10 * time.Second
	connectTimeout    = 5 * time.Second
)

func main() {
	if err := run(); err != nil {
		slog.Error("server terminated with error", "error", err)
		os.Exit(1)
	}
}

// run carga la configuracion, arma las dependencias y sirve hasta recibir una senal de apagado.
func run() error {
	cfg := config.Load(os.Getenv)
	configureLogger(cfg.LogLevel)

	repo, cleanup, err := buildRepository(cfg)
	if err != nil {
		return err
	}
	defer cleanup()
	products := usecase.NewProductUseCase(repo, uuid.NewString, time.Now)

	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           server.NewHandler(cfg, products),
		ReadHeaderTimeout: readHeaderTimeout,
		ReadTimeout:       readTimeout,
		WriteTimeout:      writeTimeout,
		IdleTimeout:       idleTimeout,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	slog.Info("starting server", "port", cfg.Port, "env", cfg.AppEnv, "repo", cfg.RepoDriver)
	return serve(ctx, srv)
}

// serve arranca el servidor y lo apaga de forma ordenada cuando el contexto se cancela.
// Recibe el contexto ya armado, por lo que es testeable sin enviar senales al proceso.
func serve(ctx context.Context, srv *http.Server) error {
	serverErr := make(chan error, 1)
	go func() {
		serverErr <- srv.ListenAndServe()
	}()

	select {
	case err := <-serverErr:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return fmt.Errorf("listen and serve: %w", err)
	case <-ctx.Done():
		slog.Info("shutdown signal received, draining connections")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()
		if err := srv.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("graceful shutdown: %w", err)
		}
		return nil
	}
}

// buildRepository selecciona la implementacion del repositorio segun REPO_DRIVER y devuelve
// una funcion de limpieza (no-op para memory, cierre del pool para postgres).
func buildRepository(cfg config.Config) (domain.ProductRepository, func(), error) {
	switch cfg.RepoDriver {
	case "memory":
		return memory.New(), func() {}, nil
	case "postgres":
		return buildPostgresRepository(cfg)
	default:
		return nil, nil, fmt.Errorf("unknown repo driver %q", cfg.RepoDriver)
	}
}

// buildPostgresRepository conecta a la base, aplica las migraciones y devuelve el repositorio
// junto con el cierre del pool. Usa un timeout acotado para fallar rapido si la base no responde.
func buildPostgresRepository(cfg config.Config) (domain.ProductRepository, func(), error) {
	if cfg.DatabaseURL == "" {
		return nil, nil, errors.New("DATABASE_URL is required when REPO_DRIVER=postgres")
	}
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout)
	defer cancel()

	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		return nil, nil, fmt.Errorf("connect to postgres: %w", err)
	}
	if err := postgres.Migrate(ctx, pool); err != nil {
		pool.Close()
		return nil, nil, fmt.Errorf("run migrations: %w", err)
	}
	return postgres.New(pool), pool.Close, nil
}

// configureLogger fija el logger estructurado por defecto segun el nivel indicado.
func configureLogger(level string) {
	var lvl slog.Level
	if err := lvl.UnmarshalText([]byte(level)); err != nil {
		lvl = slog.LevelInfo
	}
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: lvl})))
}
