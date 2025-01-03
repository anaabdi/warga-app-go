package handler

import (
	"context"
	"net/http"

	"github.com/anaabdi/warga-app-go/api/v1"
	"github.com/anaabdi/warga-app-go/internal/api/middlewares"
	"github.com/go-chi/chi/v5"
)

func NewHandler(ctx context.Context, params Params) (http.Handler, error) {
	router := chi.NewRouter()

	router.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		params.ServerImpl.GetPing(w, r)
	})

	authMiddleware := middlewares.NewAuthMiddleware(params.Responder, params.Cfg)

	var handlerMiddlewares []api.MiddlewareFunc

	handlerMiddlewares = append(handlerMiddlewares, authMiddleware.Handler)

	router.Mount("/api/v1", api.HandlerWithOptions(params.ServerImpl, api.ChiServerOptions{
		BaseRouter:       router,
		BaseURL:          "/api/v1",
		Middlewares:      handlerMiddlewares,
		ErrorHandlerFunc: nil,
	}))

	return router, nil
}
