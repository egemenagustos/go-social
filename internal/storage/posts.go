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
	Version   int       `json:"version"`
	Comments  []Comment `json:"comments"`
	User      User      `json:"user"`
}

type PostWithMetadata struct {
	Post
	CommentsCount int `json:"comments_count"`
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

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

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
		SELECT id, title, content, user_id, created_at, tags, updated_at, version
		FROM posts 
		where id=$1
		`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

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
		&post.Version,
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

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

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
		UPDATE posts SET title=$1, content=$2, version = version + 1
		WHERE id = $3 AND version = $4
		RETURNING version
		`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRowContext(
		ctx,
		query,
		post.Title,
		post.Content,
		post.Id,
		post.Version,
	).Scan(
		&post.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrNotFound
		default:
			return err
		}
	}

	return nil
}

func (s *PostStore) GetUserFeed(ctx context.Context, userId string, pfq PaginatedFeedQuery) ([]PostWithMetadata, error) {
	query := `
			SELECT
				p.id,
				p.user_id,
				p.title,
				p.content,
				p.created_at,
				p.version,
				p.tags,
				u.username,
				COUNT(c.id) AS comments_count
			FROM posts p
			LEFT JOIN comments c ON c.post_id = p.id
			LEFT JOIN users u ON p.user_id = u.id
			JOIN followers f
				ON f.follower_id = p.user_id
				OR p.user_id = $1
			WHERE f.user_id = $1
				AND (
					p.title ILIKE '%' || $4 || '%'
					OR p.content ILIKE '%' || $4 || '%'
				)
				AND (
					p.tags @> $5::varchar[]
					OR $5::varchar[] = '{}'
				)
			GROUP BY p.id, u.username
			ORDER BY p.created_at ` + pfq.Sort + `
			LIMIT $2 OFFSET $3
			`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(
		ctx,
		query,
		userId,
		pfq.Limit,
		pfq.Offset,
		pfq.Search,
		pq.Array(pfq.Tags),
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var feed []PostWithMetadata

	for rows.Next() {
		var p PostWithMetadata
		err := rows.Scan(
			&p.Id,
			&p.UserId,
			&p.Title,
			&p.Content,
			&p.CreatedAt,
			&p.Version,
			pq.Array(&p.Tags),
			&p.User.Username,
			&p.CommentsCount,
		)

		if err != nil {
			return nil, err
		}

		feed = append(feed, p)
	}

	return feed, nil
}
