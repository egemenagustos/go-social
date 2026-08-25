package store

import (
	"context"
	"database/sql"
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

func (s *CommentStore) GetPostById(ctx context.Context, postId string) ([]Comment, error) {

	query := `
	SELECT c.id, c.post_id, c.user_id, c.content, c.created_at, u.username, u.id 
	FROM comments c 
    JOIN users u ON c.user_id = u.id 
	WHERE c.post_id = $1
	ORDER BY c.created_at DESC;
	`

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
