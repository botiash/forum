package storage

import (
	"database/sql"
	"fmt"
	"forum/internal/models"
	"log"
)

type CommentIR interface {
	CreateComment(id int, username string, comment string) error
	GetCommentsByIdPost(id int) ([]models.Comment, error)
}

type CommentStorage struct {
	db *sql.DB
}

func newCommentStorage(db *sql.DB) CommentIR {
	return &CommentStorage{
		db: db,
	}
}

func (c *CommentStorage) CreateComment(id int, user, comment string) error {
	_, err := c.db.Exec(`INSERT INTO comment(id_post, author, comment) VALUES (?, ?, ?)`, id, user, comment)
	if err != nil {
		return fmt.Errorf("repo: create comment: falied %w", err)
	}
	return nil
}

func (c *CommentStorage) GetCommentsByIdPost(id int) ([]models.Comment, error) {
	comments := []models.Comment{}
	query := `SELECT id, author, comment, likes, dislikes, created_at FROM comment WHERE id_post=$1`
	rows, err := c.db.Query(query, id)
	if err != nil {
		return nil, fmt.Errorf("storage: comment by id post: %w", err)
	}
	for rows.Next() {
		var comment models.Comment
		if err := rows.Scan(&comment.Id, &comment.Creator, &comment.Text, &comment.Likes, &comment.Dislikes, &comment.Created_at); err != nil {
			log.Println(err.Error())
			return nil, fmt.Errorf("storage: comment by id post: %w", err)
		}
		comments = append(comments, comment)
	}
	return comments, err
}
