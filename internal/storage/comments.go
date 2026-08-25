package store

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

type Comment struct {
	Id        string `json:"id"`
	PostId    string `json:"post_id"`
	UserId    string `json:"user_id"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
	User      User   `json:"user"`
}

type CommentStore struct {
	db *sql.DB
}

func (s *CommentStore) Create(ctx context.Context, comment *Comment) error {

	id, err := uuid.NewV7()
	if err != nil {
		return err
	}

	comment.Id = id.String()

	query :=
		`
	INSERT INTO comments (id,post_id, user_id,content)
	VALUES($1, $2, $3, $4)
	RETURNING id, created_at
		`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err = s.db.QueryRowContext(
		ctx,
		query,
		&comment.Id,
		&comment.PostId,
		&comment.UserId,
		&comment.Content,
	).Scan(
		&comment.Id,
		&comment.CreatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func (s *CommentStore) GetPostById(ctx context.Context, postId string) ([]Comment, error) {

	query := `
	SELECT c.id, c.post_id, c.user_id, c.content, c.created_at, u.username, u.id 
	FROM comments c 
    JOIN users u ON c.user_id = u.id 
	WHERE c.post_id = $1
	ORDER BY c.created_at DESC;
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(
		ctx,
		query,
		postId,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	comments := []Comment{}

	for rows.Next() {
		var c Comment
		c.User = User{}

		err := rows.Scan(
			&c.Id,
			&c.PostId,
			&c.UserId,
			&c.Content,
			&c.CreatedAt,
			&c.User.Username,
			&c.User.Id,
		)

		if err != nil {
			return nil, err
		}

		comments = append(comments, c)
	}

	return comments, nil
}
