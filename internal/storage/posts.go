package store

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type Post struct {
	Id        string   `json:"id"`
	Content   string   `json:"content"`
	Title     string   `json:"title"`
	UserId    string   `json:"user_id"`
	Tags      []string `json:"tags"`
	CreatedAt string   `json:"created_at"`
	UpdatedAt string   `json:"updated_at"`
}

type PostStore struct {
	db *sql.DB
}

func (s *PostStore) Create(ctx context.Context, post *Post) error {

	id, err := uuid.NewV7()
	if err != nil {
		return err
	}

	post.Id = id.String()

	query := `
	INSERT INTO posts (id,content,title,user_id,tags)
	VALUES($1, $2, $3, $4, $5) RETURNING id, created_at, updated_at
	`

	err = s.db.QueryRowContext(
		ctx,
		query,
		post.Id,
		post.Content,
		post.Title,
		post.UserId,
		pq.Array(post.Tags),
	).Scan(
		&post.Id,
		&post.CreatedAt,
		&post.UpdatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}
