package handler

import (
	"context"
	"net/http"

	"github.com/anaabdi/warga-app-go/api/v1"
	"github.com/go-chi/chi/v5"
)

func NewHandler(ctx context.Context, params Params) (http.Handler, error) {
	router := chi.NewRouter()

	router.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		params.ServerImpl.GetPing(w, r)
	})

	router.Mount("/api/v1", api.HandlerWithOptions(params.ServerImpl, api.ChiServerOptions{
		BaseRouter:       router,
		BaseURL:          "/api/v1",
		Middlewares:      nil,
		ErrorHandlerFunc: nil,
	}))

	return router, nil
}
