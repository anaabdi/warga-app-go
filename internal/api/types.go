package api

import (
	"github.com/anaabdi/warga-app-go/cmd/app/config"
	"github.com/anaabdi/warga-app-go/internal/api/parser"
	"github.com/anaabdi/warga-app-go/internal/repository"
	"github.com/anaabdi/warga-app-go/internal/repository/user"
	"github.com/anaabdi/warga-app-go/internal/service/auth"
	"github.com/anaabdi/warga-app-go/pkg/postgres"
)

type ServerImpl struct {
	Config    *config.Config
	Responder parser.JSONResponder
	Requestor parser.RequestParser
	DB        *postgres.DB

	// add more service here
	authService auth.AuthService

	// add more repository here
	userRepository user.User
}

func NewServerImpl(config *config.Config, responder parser.JSONResponder,
	requestor parser.RequestParser,
	db *postgres.DB) *ServerImpl {

	transactionalDB := repository.NewDatabase(db)

	userRepository := user.NewUser(db)
	authService := auth.NewAuthService(userRepository, transactionalDB)

	return &ServerImpl{
		Config:         config,
		Responder:      responder,
		DB:             db,
		authService:    authService,
		userRepository: userRepository,
	}
}
