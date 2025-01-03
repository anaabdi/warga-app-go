package auth

import (
	"context"

	"github.com/anaabdi/warga-app-go/internal/repository"
	"github.com/anaabdi/warga-app-go/internal/repository/user"
)

type AuthService struct {
	transactionalDB repository.Database
	userrepo        user.User
}

type AuthServiceInterface interface {
	Login(ctx context.Context, username, password string) (string, error)
}

func NewAuthService(userrepo user.User, transactionalDB repository.Database) AuthService {
	return AuthService{
		userrepo: userrepo,
	}
}

func (a AuthService) Login(ctx context.Context, username, password string) (string, error) {
	_, err := a.userrepo.FindAccountByUsername(ctx, username)
	if err != nil {
		return "", err
	}

	return "token", nil
}
