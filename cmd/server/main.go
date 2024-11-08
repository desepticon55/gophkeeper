package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/desepticon55/gophkeeper/internal/server"
	"github.com/desepticon55/gophkeeper/internal/server/api/auth"
	"github.com/desepticon55/gophkeeper/internal/server/api/secret"
	customMiddleware "github.com/desepticon55/gophkeeper/internal/server/middleware"
	secretSrv "github.com/desepticon55/gophkeeper/internal/server/service/secret"
	"github.com/desepticon55/gophkeeper/internal/server/service/user"
	"github.com/desepticon55/gophkeeper/internal/server/storage"
	"github.com/desepticon55/gophkeeper/pkg/logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/jackc/pgx/v4/stdlib"
	"github.com/pressly/goose/v3"
	"go.uber.org/zap"
	"net/http"
	"time"
)

func main() {
	log := logger.InitLogger()
	defer log.Sync()

	config := parseConfig()
	log.Debug("Config created",
		zap.String("Server address", config.ServerAddress),
		zap.String("Database connection string", config.DatabaseConnString),
		zap.String("Auth key", config.AuthKey),
		zap.Int("Token expired after minutes", config.ExpirationMinutes))

	router := chi.NewRouter()
	router.Use(middleware.Recoverer)
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Timeout(60 * time.Second))
	router.Use(customMiddleware.CompressingMiddleware())
	router.Use(customMiddleware.DecompressingMiddleware())

	pool, err := createConnectionPool(context.Background(), config.DatabaseConnString)
	if err != nil {
		log.Fatal("Error during initialize DB connection", zap.Error(err))
	}
	runMigrations(config.DatabaseConnString, log)

	userRepository := storage.NewUserRepository(pool)
	userService := user.NewUserService(log, userRepository)

	secretRepository := storage.NewSecretRepository(pool)
	secretService := secretSrv.NewSecretService(log, secretRepository)

	router.Method(http.MethodPost, "/api/user/register", auth.RegisterHandler(log, config, userService))
	router.Method(http.MethodPost, "/api/user/login", auth.LoginHandler(log, config, userService))

	router.Group(func(r chi.Router) {
		r.Use(customMiddleware.CheckAuthMiddleware(log, config))
		r.Method(http.MethodPost, "/api/user/secret", secret.UploadSecretHandler(log, secretService))
		r.Method(http.MethodGet, "/api/user/secret/{name}", secret.ReadOneSecretHandler(log, secretService))
		r.Method(http.MethodDelete, "/api/user/secret/{name}", secret.DeleteSecretHandler(log, secretService))
		r.Method(http.MethodGet, "/api/user/secret", secret.ReadAllSecretsHandler(log, secretService))
	})

	http.ListenAndServe(config.ServerAddress, router)
}

func createConnectionPool(ctx context.Context, connectionString string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(connectionString)
	if err != nil {
		return nil, fmt.Errorf("error parsing database config: %w", err)
	}

	pool, err := pgxpool.ConnectConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("error connecting to the database: %w", err)
	}

	return pool, nil
}

func runMigrations(connectionString string, log *zap.Logger) {
	databaseConfig, err := pgx.ParseConfig(connectionString)
	if err != nil {
		log.Error("Error during parse database URL", zap.Error(err))
		return
	}
	db := stdlib.OpenDB(*databaseConfig)
	defer db.Close()

	goose.SetDialect("postgres")
	if err := goose.Up(db, "migrations"); err != nil {
		log.Error("Error during run database migrations", zap.Error(err))
	}
}

func parseConfig() server.Config {
	config := server.ParseConfig()
	flag.Parse()
	return config
}
