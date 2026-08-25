package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrNotFound          = errors.New("Record not found!")
	QueryTimeoutDuration = time.Second * 5
)

type Storage struct {
	Posts interface {
		Create(context.Context, *Post) error
		GetById(context.Context, string) (*Post, error)
		Delete(context.Context, string) error
		Update(context.Context, *Post) error
	}

	Users interface {
		Create(context.Context, *User) error
		GetById(context.Context, string) (*User, error)
	}

	Comments interface {
		Create(context.Context, *Comment) error
		GetPostById(ctx context.Context, postId string) ([]Comment, error)
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Posts:    &PostStore{db: db},
		Users:    &UserStore{db: db},
		Comments: &CommentStore{db: db},
	}
}
