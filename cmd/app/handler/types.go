package handler

import (
	"github.com/anaabdi/warga-app-go/api/v1"
	"github.com/anaabdi/warga-app-go/cmd/app/config"
	"github.com/anaabdi/warga-app-go/internal/api/parser"
	"github.com/anaabdi/warga-app-go/pkg/postgres"
)

type Params struct {
	Cfg        *config.Config
	Responder  parser.JSONResponder
	DB         *postgres.DB
	ServerImpl api.ServerInterface
}
