package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/anaabdi/warga-app-go/cmd/app"
	"github.com/anaabdi/warga-app-go/cmd/app/config"
)

func main() {
	// This is the main function for the server.
	ctx, terminateFn := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill, syscall.SIGTERM)
	defer terminateFn()

	// Prepare configuration
	appCfg := config.NewConfig()

	// Prepare logger

	// Prepare app
	app := app.NewApp(appCfg)

	// Start app
	if err := app.Start(ctx); err != nil {
		panic(err)
	}

	slog.Info("Server is shutting down")
}
