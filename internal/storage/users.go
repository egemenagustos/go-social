package store

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

type User struct {
	Id        string `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"-"`
	CreatedAt string `json:"created_at"`
}

type UserStore struct {
	db *sql.DB
}

func (u *UserStore) Create(ctx context.Context, user *User) error {

	id, err := uuid.NewV7()
	if err != nil {
		return err
	}

	user.Id = id.String()

	query := `
	INSERT INTO users (id,username,email,password)
	VALUES($1, $2, $3, $4) RETURNING id, created_at
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err = u.db.QueryRowContext(
		ctx,
		query,
		user.Id,
		user.Username,
		user.Email,
		user.Password,
	).Scan(
		&user.Id,
		&user.CreatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func (u *UserStore) GetById(ctx context.Context, id string) (*User, error) {

	user := &User{}

	query := `SELECT id, email, username, password, created_at FROM users where id = $1`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := u.db.QueryRowContext(
		ctx,
		query,
		id,
	).Scan(
		&user.Id,
		&user.Email,
		&user.Username,
		&user.Password,
		&user.CreatedAt,
	)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return user, nil
}
