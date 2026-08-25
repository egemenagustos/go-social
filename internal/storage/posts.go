package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type Post struct {
	Id        string    `json:"id"`
	Content   string    `json:"content"`
	Title     string    `json:"title"`
	UserId    string    `json:"user_id"`
	Tags      []string  `json:"tags"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
	Comments  []Comment `json:"comments"`
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

func (s *PostStore) GetById(ctx context.Context, id string) (*Post, error) {
	var post Post
	query :=
		`
		SELECT id, title, content, user_id, created_at, tags, updated_at
		FROM posts 
		where id=$1
		`

	err := s.db.QueryRowContext(
		ctx,
		query,
		id,
	).Scan(
		&post.Id,
		&post.Title,
		&post.Content,
		&post.UserId,
		&post.CreatedAt,
		pq.Array(&post.Tags),
		&post.UpdatedAt,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return &post, nil
}

func (s *PostStore) Delete(ctx context.Context, id string) error {

	query := `DELETE FROM posts where id = $1`

	result, err := s.db.ExecContext(
		ctx,
		query,
		id,
	)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

func (s *PostStore) Update(ctx context.Context, post *Post) error {
	query :=
		`
		UPDATE posts SET title=$1, content=$2
		WHERE id = $3
		`

	_, err := s.db.ExecContext(
		ctx,
		query,
		post.Title,
		post.Content,
		post.Id,
	)

	if err != nil {
		return err
	}

	return nil
}
