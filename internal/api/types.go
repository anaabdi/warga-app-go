package api

import "github.com/anaabdi/warga-app-go/cmd/app/config"

type ServerImpl struct {
	Config *config.Config
}

func NewServerImpl(config *config.Config) *ServerImpl {
	return &ServerImpl{
		Config: config,
	}
}
