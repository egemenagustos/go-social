package store

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrDuplicateEmail    = errors.New("A user with that email already exist!")
	ErrDuplicateUsername = errors.New("A user with that username already exist!")
)

type User struct {
	Id        string   `json:"id"`
	Username  string   `json:"username"`
	Email     string   `json:"email"`
	Password  password `json:"-"`
	CreatedAt string   `json:"created_at"`
	IsActive  bool     `json:"is_active"`
}

type password struct {
	text *string
	hash []byte
}

func (p *password) Set(text string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(text), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	p.text = &text
	p.hash = hash
	return nil
}

type UserStore struct {
	db *sql.DB
}

func (u *UserStore) Create(ctx context.Context, tx *sql.Tx, user *User) error {

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

	err = tx.QueryRowContext(
		ctx,
		query,
		user.Id,
		user.Username,
		user.Email,
		user.Password.hash,
	).Scan(
		&user.Id,
		&user.CreatedAt,
	)

	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmail
		case err.Error() == `pq: duplicate key value violates unique constraint "users_username_key`:
			return ErrDuplicateUsername
		default:
			return err
		}
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

func (u *UserStore) CreateAndInvite(ctx context.Context, user *User, token string, invitationExp time.Duration) error {
	return withTx(u.db, ctx, func(tx *sql.Tx) error {
		if err := u.Create(ctx, tx, user); err != nil {
			return err
		}

		if err := u.createUserInvitation(ctx, tx, token, invitationExp, user.Id); err != nil {
			return err
		}

		return nil
	})
}

func (u *UserStore) Activate(ctx context.Context, token string) error {

	return withTx(u.db, ctx, func(tx *sql.Tx) error {

		user, err := u.getUserFromInvitation(ctx, tx, token)

		if err != nil {
			return err
		}

		user.IsActive = true
		if err := u.update(ctx, tx, user); err != nil {
			return err
		}

		if err := u.deleteUserInvitations(ctx, tx, user.Id); err != nil {
			return err
		}

		return nil
	})
}

func (u *UserStore) createUserInvitation(ctx context.Context, tx *sql.Tx, token string, invitationExp time.Duration, userId string) error {

	query := `INSERT INTO user_invitations (token, user_id, expiry) VALUES ($1,$2,$3)`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, token, userId, time.Now().Add(invitationExp))
	if err != nil {
		return err
	}

	return nil
}

func (u *UserStore) getUserFromInvitation(ctx context.Context, tx *sql.Tx, token string) (*User, error) {

	query := `
	SELECT u.id, u.username, u.email, u.created_at, u.is_active
	FROM users u
	JOIN user_invitations ui ON u.id = ui.user_id
	WHERE ui.token = $1 AND ui.expiry > $2
	 `

	hash := sha256.Sum256([]byte(token))
	hashToken := hex.EncodeToString(hash[:])

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	user := &User{}

	err := tx.QueryRowContext(
		ctx,
		query,
		hashToken,
		time.Now(),
	).Scan(
		&user.Id,
		&user.Username,
		&user.Email,
		&user.CreatedAt,
		&user.IsActive,
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

func (u *UserStore) update(ctx context.Context, tx *sql.Tx, user *User) error {

	query :=
		`
	UPDATE users 
	set username = $1, email = $2, is_active = $3
	WHERE id = $4
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := tx.ExecContext(
		ctx,
		query,
		user.Username,
		user.Email,
		user.IsActive,
		user.Id,
	)

	if err != nil {
		return err
	}

	return nil
}

func (u *UserStore) deleteUserInvitations(ctx context.Context, tx *sql.Tx, userId string) error {

	query :=
		`
	DELETE FROM user_invitations 
	WHERE user_id = $1
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := tx.ExecContext(
		ctx,
		query,
		userId,
	)

	if err != nil {
		return err
	}

	return nil
}
