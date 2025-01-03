package user

import (
	"context"

	"github.com/anaabdi/warga-app-go/internal/entity"
	"github.com/anaabdi/warga-app-go/pkg/postgres"
)

type UserInterface interface {
	CreateUser(ctx context.Context, arg CreateUserParams) (int32, error)
	DeleteAuthor(ctx context.Context, id int32) error
	FindUserByEmail(ctx context.Context, email string) (*entity.User, error)
	ListUser(ctx context.Context) ([]*entity.User, error)
}

type user struct {
	db *postgres.DB
}

// NewUser ...
func NewUser(db *postgres.DB) UserInterface {
	return &user{
		db: db,
	}
}

const createUser = `-- name: CreateUser :one
INSERT INTO users (
 name, email, password, role, id_card_number, id_family_card_number
) VALUES (
 $1, $2, $3, $4, $5, $6
)
RETURNING id
`

type CreateUserParams struct {
	Name               string
	Email              string
	Password           string
	Role               string
	IDCardNumber       string
	IDFamilyCardNumber string
}

func (q *user) CreateUser(ctx context.Context, arg CreateUserParams) (int32, error) {
	row := q.db.Pool.QueryRowContext(ctx, createUser,
		arg.Name,
		arg.Email,
		arg.Password,
		arg.Role,
		arg.IDCardNumber,
		arg.IDFamilyCardNumber,
	)
	var id int32
	err := row.Scan(&id)
	return id, err
}

const deleteAuthor = `-- name: DeleteAuthor :exec
DELETE FROM users
WHERE id = $1
`

func (q *user) DeleteAuthor(ctx context.Context, id int32) error {
	_, err := q.db.Pool.ExecContext(ctx, deleteAuthor, id)
	return err
}

const findUserByEmail = `-- name: FindUserByEmail :one
SELECT id, name, email, password, role, id_card_number, id_family_card_number, created_at, updated_at FROM users
WHERE email = $1 LIMIT 1
`

func (q *user) FindUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	row := q.db.Pool.QueryRowContext(ctx, findUserByEmail, email)
	var i entity.User
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Email,
		&i.Password,
		&i.Role,
		&i.IDCardNumber,
		&i.IDFamilyCardNumber,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return &i, err
}

const listUser = `-- name: ListUser :many
SELECT id, name, email, password, role, id_card_number, id_family_card_number, created_at, updated_at FROM users
ORDER BY name
`

func (q *user) ListUser(ctx context.Context) ([]*entity.User, error) {
	rows, err := q.db.Pool.QueryContext(ctx, listUser)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*entity.User
	for rows.Next() {
		var i *entity.User
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Email,
			&i.Password,
			&i.Role,
			&i.IDCardNumber,
			&i.IDFamilyCardNumber,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
