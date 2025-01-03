package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"

	"github.com/anaabdi/warga-app-go/cmd/app/config"
	"github.com/anaabdi/warga-app-go/cmd/app/handler"
	intapi "github.com/anaabdi/warga-app-go/internal/api"
	"github.com/anaabdi/warga-app-go/internal/api/parser"
	dbmigration "github.com/anaabdi/warga-app-go/pkg/db/migration"
	"github.com/anaabdi/warga-app-go/pkg/postgres"
)

type App struct {
	Config *config.Config
	DB     *postgres.DB
}

func NewApp(config *config.Config) *App {
	return &App{
		Config: config,
	}
}

func (a *App) Start(ctx context.Context) error {
	// Start migration
	if err := a.InitMigration(ctx); err != nil {
		return err
	}

	// Start app
	httpServer, err := a.InitHTTPServer(ctx)
	if err != nil {
		return err
	}

	go func() {
		<-ctx.Done()
		httpServer.Shutdown(ctx)
	}()

	slog.Info(fmt.Sprintf("Server %s is starting", a.Config.AppName), slog.Any("addr", httpServer.Addr))

	return httpServer.ListenAndServe()
}

func (a *App) InitHTTPServer(ctx context.Context) (*http.Server, error) {
	// Prepare HTTP server
	serverImpl := intapi.NewServerImpl(a.Config, parser.NewJSONResponder(), parser.NewRequestParser(), a.DB)

	// Prepare handlers
	handler, err := handler.NewHandler(ctx, handler.Params{
		ServerImpl: serverImpl,
	})

	if err != nil {
		return nil, err
	}

	return &http.Server{
		Addr:              fmt.Sprintf("%s:%s", a.Config.ServerHost, a.Config.ServerPort),
		ReadHeaderTimeout: a.Config.ReadHeaderTimeout, // if not defined, potential Slowloris Attack (Dos)
		Handler:           handler,
	}, nil
}

func (a *App) InitMigration(ctx context.Context) error {
	disableMigrationValue := os.Getenv("DISABLE_MIGRATION")
	disableMigration, _ := strconv.ParseBool(disableMigrationValue)

	var (
		databaseURL, migrationEnv, migrationType string
		envFound                                 bool
	)

	if !disableMigration {
		databaseURL, envFound = os.LookupEnv("DB_URL")
		if !envFound || len(databaseURL) == 0 {
			return fmt.Errorf("migrate: environment variable not declared: DB_URL")
		}

		databaseURL = "postgres://" + databaseURL

		migrationEnv, envFound = os.LookupEnv("MIGRATION_ENV")
		if !envFound || len(migrationEnv) == 0 {
			return fmt.Errorf("migrate: environment variable not declared: MIGRATION_ENV")
		}

		migrationType, envFound = os.LookupEnv("MIGRATION_TYPE")
		if envFound {
			dbmigration.Exec(migrationType, databaseURL, migrationEnv)
			return nil
		}

		migrationType = "up"
		dbmigration.Exec(migrationType, databaseURL, migrationEnv)
	}

	return nil
}

func (a *App) InitDB() error {
	cfg := a.Config
	dbURL := "postgres://" + cfg.DB.URL
	db, err := postgres.New(dbURL, postgres.MaxOpenConns(cfg.DB.MaxOpenConns),
		postgres.MaxIdleConns(cfg.DB.MaxIdleConns), postgres.ConnMaxLifetime(cfg.DB.ConnMaxLifetime))
	if err != nil {
		return fmt.Errorf("error setup db: %w", err)
	}

	if err := db.Pool.Ping(); err != nil {
		return fmt.Errorf("error pinging to db: %w", err)
	}

	a.DB = db

	return nil
}
