package user

import (
	"context"

	"github.com/anaabdi/warga-app-go/internal/entity"

	"github.com/anaabdi/warga-app-go/pkg/postgres"
)

type User interface {
	FindAccountByID(ctx context.Context, ID int) (*entity.User, error)
	FindAccountByUsername(ctx context.Context, username string) (*entity.User, error)
}

type user struct {
	db *postgres.DB
}

// NewUser ...
func NewUser(db *postgres.DB) User {
	return &user{
		db: db,
	}
}

func (u *user) FindAccountByID(ctx context.Context, ID int) (*entity.User, error) {
	user := &entity.User{}
	// query := `SELECT id, username, password, role FROM users WHERE id = $1`
	// err := u.db.Pool.QueryRowContext(ctx, query, ID).Scan(&user.ID, &user.Username, &user.Password, &user.Role)
	// if err != nil {
	// 	return nil, err
	// }
	return user, nil
}

func (u *user) FindAccountByUsername(ctx context.Context, username string) (*entity.User, error) {
	user := &entity.User{}
	// query := `SELECT id, username, password, role FROM users WHERE username = $1`
	// err := u.db.Pool.QueryRowContext(ctx, query, username).Scan(&user.ID, &user.Username, &user.Password, &user.Role)
	// if err != nil {
	// 	return nil, err
	// }
	return user, nil
}
