package store

import (
	"context"
	"database/sql"
	"errors"
)

var (
	ErrNotFound = errors.New("Record not found!")
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
	}

	Comments interface {
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
