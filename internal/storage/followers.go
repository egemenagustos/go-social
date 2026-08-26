package store

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
)

type Follower struct {
	UserId     string `json:"user_id"`
	FollowerId string `json:"follower_id"`
	CreatedAt  string `json:"created_at"`
}

type FollowerStore struct {
	db *sql.DB
}

func (f *FollowerStore) Follow(ctx context.Context, followerUserId, userId string) error {
	query := `INSERT INTO followers(user_id, follower_id) VALUES($1, $2)`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := f.db.ExecContext(
		ctx,
		query,
		userId,
		followerUserId,
	)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return ErrConflict
		}

		return err
	}

	return nil
}

func (f *FollowerStore) Unfollow(ctx context.Context, followerUserId, userId string) error {
	query := `DELETE FROM followers WHERE user_id = $1 AND follower_id = $2`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := f.db.ExecContext(
		ctx,
		query,
		userId,
		followerUserId,
	)

	if err != nil {
		return err
	}

	return nil
}
