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
		Create(context.Context, *User) error
		GetById(context.Context, string) (*User, error)
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
