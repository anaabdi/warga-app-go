package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/anaabdi/warga-app-go/cmd/app/config"
	"github.com/anaabdi/warga-app-go/cmd/app/handler"
	intapi "github.com/anaabdi/warga-app-go/internal/api"
)

type App struct {
	Config *config.Config
}

func NewApp(config *config.Config) *App {
	return &App{
		Config: config,
	}
}

func (a *App) Start(ctx context.Context) error {
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
	serverImpl := intapi.NewServerImpl(a.Config)

	// Prepare handlers
	handler, err := handler.NewHandler(ctx, handler.Params{
		ServerImpl: serverImpl,
	})

	if err != nil {
		return nil, err
	}

	return &http.Server{
		Addr:              fmt.Sprintf("%s:%s", a.Config.ServerHost, a.Config.ServerPort),
		ReadHeaderTimeout: a.Config.ReadHeaderTimeout,
		Handler:           handler,
	}, nil
}
