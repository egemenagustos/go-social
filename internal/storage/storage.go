package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrNotFound          = errors.New("Record not found!")
	ErrConflict          = errors.New("Resource already exist!")
	QueryTimeoutDuration = time.Second * 5
)

type Storage struct {
	Posts interface {
		Create(context.Context, *Post) error
		GetById(context.Context, string) (*Post, error)
		Delete(context.Context, string) error
		Update(context.Context, *Post) error
		GetUserFeed(context.Context, string, PaginatedFeedQuery) ([]PostWithMetadata, error)
	}

	Users interface {
		Create(context.Context, *sql.Tx, *User) error
		GetById(context.Context, string) (*User, error)
		CreateAndInvite(ctx context.Context, user *User, token string, invitationExp time.Duration) error
		Activate(context.Context, string) error
		Delete(context.Context, string) error
		GetByEmail(context.Context, string) (*User, error)
	}

	Comments interface {
		Create(context.Context, *Comment) error
		GetPostById(ctx context.Context, postId string) ([]Comment, error)
	}

	Followers interface {
		Follow(ctx context.Context, followerUserId, userId string) error
		Unfollow(ctx context.Context, followerUserId, userId string) error
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Posts:     &PostStore{db: db},
		Users:     &UserStore{db: db},
		Comments:  &CommentStore{db: db},
		Followers: &FollowerStore{db: db},
	}
}

func withTx(db *sql.DB, ctx context.Context, fn func(*sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	if err := fn(tx); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}
