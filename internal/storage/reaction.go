package storage

import (
	"database/sql"
	"forum/internal/models"
)

type ReactionIR interface {
	CreateEmotionPost(models.Like) error
	GetEmotionPost(models.Like) (models.Like, error)
	EmotionPostExists(postID, userID int) (bool, error)
	UpdateEmotionPost(models.Like) error
	EmotionPostExistsFull(models.Like) (bool, error)

	CreateEmotionComment(models.Like) error
	GetEmotionComment(models.Like) (models.Like, error)
	EmotionCommentExists(commentID, userID int) (bool, error)
	UpdateEmotionComment(models.Like) error
	EmotionCommentExistsFull(CommentEmo models.Like) (bool, error)
}

type EmotionSQL struct {
	db *sql.DB
}

func NewEmotionSQL(db *sql.DB) ReactionIR {
	return &EmotionSQL{
		db: db,
	}
}

func (e *EmotionSQL) CreateEmotionPost(PostEmo models.Like) error {
	stmt, err := e.db.Prepare("INSERT INTO likesPost(userId, postId,like1) values(?,?,?)")
	if err != nil {
		return err
	}
	if _, err := stmt.Exec(PostEmo.UserID, PostEmo.PostID, PostEmo.Islike); err != nil {
		return err
	}
	asd, err := e.GetEmotionPost(PostEmo)
	if err != nil {
		return err
	}

	query := `UPDATE post SET likes = $1 , dislikes  = $2 WHERE id = $3`
	if _, err := e.db.Exec(query, asd.CountLike, asd.Countdislike, PostEmo.PostID); err != nil {
		return err
	}
	return nil
}

func (e *EmotionSQL) EmotionPostExistsFull(PostEmo models.Like) (bool, error) {
	query := "SELECT COUNT(*) FROM likesPost WHERE postId=$1 AND userId=$2 and like1 =$3 "
	var count int
	err := e.db.QueryRow(query, PostEmo.PostID, PostEmo.UserID, PostEmo.Islike).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (e *EmotionSQL) EmotionCommentExistsFull(CommentEmo models.Like) (bool, error) {
	query := "SELECT COUNT(*) FROM likesComment WHERE commentsId=$1 AND userId=$2 and like1 =$3 "
	var count int
	err := e.db.QueryRow(query, CommentEmo.CommentID, CommentEmo.UserID, CommentEmo.Islike).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (e *EmotionSQL) CreateEmotionComment(CommentEmo models.Like) error {
	stmt, err := e.db.Prepare("INSERT INTO likesComment(userId, commentsId,like1) values(?,?,?)")
	if err != nil {
		return err
	}
	if _, err := stmt.Exec(CommentEmo.UserID, CommentEmo.CommentID, CommentEmo.Islike); err != nil {
		return err
	}

	asd, err := e.GetEmotionComment(CommentEmo)
	if err != nil {
		return err
	}
	query := `UPDATE comment SET likes = $1 , dislikes = $2 WHERE id = $3`
	if _, err := e.db.Exec(query, asd.CountLike, asd.Countdislike, CommentEmo.CommentID); err != nil {
		return err
	}
	return nil
}

func (e *EmotionSQL) GetEmotionPost(PostEmo models.Like) (models.Like, error) {
	query := `SELECT 
    (SELECT COUNT(*) FROM likesPost WHERE postId = $1 AND like1 = 1) AS likes, 
    (SELECT COUNT(*) FROM likesPost WHERE postId = $1 AND like1 = 0) AS dislikes;`
	row := e.db.QueryRow(query, PostEmo.PostID)
	if err := row.Scan(&PostEmo.CountLike, &PostEmo.Countdislike); err != nil {
		return PostEmo, err
	}
	return PostEmo, nil
}

func (e *EmotionSQL) GetEmotionComment(CommentEmo models.Like) (models.Like, error) {
	query := `SELECT 
    (SELECT COUNT(*) FROM likesComment WHERE commentsId = $1 AND like1 = 1) AS likes, 
    (SELECT COUNT(*) FROM likesComment WHERE commentsId = $1 AND like1 = 0) AS dislikes;`
	row := e.db.QueryRow(query, CommentEmo.CommentID)
	if err := row.Scan(&CommentEmo.CountLike, &CommentEmo.Countdislike); err != nil {
		return CommentEmo, err
	}

	return CommentEmo, nil
}

func (e *EmotionSQL) EmotionCommentExists(commentID, userID int) (bool, error) {
	query := "SELECT COUNT(*) FROM likesComment WHERE commentsId=$1 AND userId=$2"
	var count int
	err := e.db.QueryRow(query, commentID, userID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (e *EmotionSQL) EmotionPostExists(postID, userID int) (bool, error) {
	query := "SELECT COUNT(*) FROM likesPost WHERE postId=$1 AND userId=$2"
	var count int
	err := e.db.QueryRow(query, postID, userID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (e *EmotionSQL) UpdateEmotionComment(CommentEmo models.Like) error {
	query := `UPDATE likesComment SET like1 = $1 WHERE userId = $2 AND commentsId = $3`
	_, err := e.db.Exec(query, CommentEmo.Islike, CommentEmo.UserID, CommentEmo.CommentID)
	if err != nil {
		return err
	}

	asd, err := e.GetEmotionComment(CommentEmo)
	if err != nil {
		return err
	}
	query = `UPDATE comment SET likes = $1 , dislikes = $2 WHERE id = $3`
	if _, err := e.db.Exec(query, asd.CountLike, asd.Countdislike, CommentEmo.CommentID); err != nil {
		return err
	}
	return nil
}

func (e *EmotionSQL) UpdateEmotionPost(PostEmo models.Like) error {
	query := `UPDATE likesPost SET like1 = $1 WHERE userId = $2 AND postId = $3`
	_, err := e.db.Exec(query, PostEmo.Islike, PostEmo.UserID, PostEmo.PostID)
	if err != nil {
		return err
	}
	asd, err := e.GetEmotionPost(PostEmo)
	if err != nil {
		return err
	}

	query = `UPDATE post SET likes = $1 , dislikes  = $2 WHERE id = $3`
	if _, err := e.db.Exec(query, asd.CountLike, asd.Countdislike, PostEmo.PostID); err != nil {
		return err
	}
	return nil
}
