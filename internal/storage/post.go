package storage

import (
	"database/sql"
	"fmt"
	"forum/internal/models"
	"strings"
)

type PostIR interface {
	CreatePost(post models.Post) error
	AddCategoriesByPostId(id int, category []string) error
	GetPostByID(id int) (models.Post, error)
	GetCategoriesByPostId(id int) ([]string, error)
	GetAllPosts() ([]models.Post, error)
	Category() ([]string, error)
	GetAllPostsByCategories(category string) ([]models.Post, error)
	GetMyPost(id int) ([]models.Post, error)
	GetMyLikedPost(id int) ([]models.Post, error)
}

type PostStorage struct {
	db *sql.DB
}

func NewPostStorage(db *sql.DB) PostIR {
	return &PostStorage{
		db: db,
	}
}

func (p *PostStorage) GetAllPosts() ([]models.Post, error) {
	posts := []models.Post{}
	query := "SELECT * FROM post"
	row, err := p.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("storage: get all posts: %w", err)
	}
	for row.Next() {
		var post models.Post
		var categoriesStr string
		if err := row.Scan(&post.Id, &post.Author, &post.Title, &post.Description, &post.Likes, &post.Likes, &categoriesStr, &post.CreateAt); err != nil {
			return nil, fmt.Errorf("storage: get all posts: %w", err)
		}
		post.Category = strings.Split(categoriesStr, ", ")

		posts = append(posts, post)
	}

	return posts, nil
}

func (p *PostStorage) CreatePost(post models.Post) error {
	query := `INSERT INTO post(title, description, author, category) VALUES ($1, $2, $3, $4);`
	var categoriesStr string
	if len(post.Category) == 1 {
		categoriesStr = post.Category[0]
	} else {
		post.Category = uniqueStrings(post.Category)
		categoriesStr = strings.Join(post.Category, ", ")
	}

	res, err := p.db.Exec(query, post.Title, post.Description, post.Author, categoriesStr)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	if err := p.AddCategoriesByPostId(int(id), post.Category); err != nil {
		return err
	}
	return nil
}

func uniqueStrings(input []string) []string {
	uniqueMap := make(map[string]struct{})
	for _, str := range input {
		uniqueMap[str] = struct{}{}
	}
	uniqueStrings := make([]string, 0, len(uniqueMap))
	for str := range uniqueMap {
		uniqueStrings = append(uniqueStrings, str)
	}
	return uniqueStrings
}

func (p *PostStorage) AddCategoriesByPostId(id int, category []string) error {
	query := `INSERT INTO category(tag, id_post) VALUES ($1, $2);`
	for _, v := range category {
		_, err := p.db.Exec(query, v, id)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *PostStorage) GetPostByID(id int) (models.Post, error) {
	query := `SELECT id, author, title, description, created_at, likes, dislikes FROM post WHERE id = $1;`
	row := p.db.QueryRow(query, id)
	var post models.Post
	if err := row.Scan(&post.Id, &post.Author, &post.Title, &post.Description, &post.CreateAt, &post.Likes, &post.Dislikes); err != nil {
		return models.Post{}, err
	}
	pCat, err := p.GetCategoriesByPostId(post.Id)
	if err == nil {
		post.Category = pCat
	}
	return post, nil
}

func (p *PostStorage) GetCategoriesByPostId(id int) ([]string, error) {
	var allTags []string
	query := `SELECT tag FROM category WHERE id_post = $1;`
	row, err := p.db.Query(query, id)
	if err != nil {
		return []string{}, err
	}

	for row.Next() {

		tag := ""
		if err := row.Scan(&tag); err != nil {
			return []string{}, err
		}
		allTags = append(allTags, tag)
	}
	return allTags, nil
}

func (p *PostStorage) Category() ([]string, error) {
	var categories []string
	query := `SELECT DISTINCT tag FROM category;`
	rows, err := p.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("storage: get all categories: %w", err)
	}
	for rows.Next() {
		cat := ""
		if err := rows.Scan(&cat); err != nil {
			return nil, fmt.Errorf("storage: get all categories: %w", err)
		}
		categories = append(categories, cat)
	}
	return categories, nil
}

func (p *PostStorage) GetAllPostsByCategories(category string) ([]models.Post, error) {
	posts := []models.Post{}
	query := `SELECT 
		p.id,
		p.author,
		p.title,
		p.description,
		p.likes,
		p.dislikes,
		p.category,
		p.created_at
	FROM 
		post p
	JOIN
		category c
	ON
		p.id = c.id_post
	WHERE
		c.tag = $1	
	`
	row, err := p.db.Query(query, category)
	if err != nil {
		return posts, fmt.Errorf("storage: get all posts: %w", err)
	}
	for row.Next() {
		var post models.Post
		var categoriesStr string
		if err := row.Scan(&post.Id, &post.Author, &post.Title, &post.Description, &post.Likes, &post.Likes, &categoriesStr, &post.CreateAt); err != nil {
			return nil, fmt.Errorf("storage: get all posts: %w", err)
		}
		post.Category = strings.Split(categoriesStr, ", ")

		posts = append(posts, post)
	}

	return posts, nil
}

func (p *PostStorage) GetMyPost(id int) ([]models.Post, error) {
	posts := []models.Post{}
	query := `SELECT 
		p.id,
		p.author,
		p.title,
		p.description,
		p.likes,
		p.dislikes,
		p.category,
		p.created_at
	FROM 
		post p
	LEFT JOIN 
		user u
	ON
		u.username = p.author
	where 
		u.id = $1`
	row, err := p.db.Query(query, id)
	if err != nil {
		return nil, fmt.Errorf("storage: get all posts: %w", err)
	}
	for row.Next() {
		var post models.Post
		var categoriesStr string
		if err := row.Scan(&post.Id, &post.Author, &post.Title, &post.Description, &post.Likes, &post.Likes, &categoriesStr, &post.CreateAt); err != nil {
			return nil, fmt.Errorf("storage: get all posts: %w", err)
		}
		post.Category = strings.Split(categoriesStr, ", ")
		posts = append(posts, post)
	}

	return posts, nil
}

func (p *PostStorage) GetMyLikedPost(id int) ([]models.Post, error) {
	posts := []models.Post{}
	query := `SELECT 
		p.id,
		p.author,
		p.title,
		p.description,
		p.likes,
		p.dislikes,
		p.category,
		p.created_at
	FROM 
		post p
	LEFT JOIN 
		user u
	ON
		u.username = p.author
	JOIN 
		likesPost lp
	ON
		lp.postId = p.id
	WHERE  
		lp.userId = $1 AND lp.like1 = 1`
	row, err := p.db.Query(query, id)
	if err != nil {
		return nil, fmt.Errorf("storage: get all posts: %w", err)
	}
	for row.Next() {
		var post models.Post
		var categoriesStr string
		if err := row.Scan(&post.Id, &post.Author, &post.Title, &post.Description, &post.Likes, &post.Likes, &categoriesStr, &post.CreateAt); err != nil {
			return nil, fmt.Errorf("storage: get all posts: %w", err)
		}
		post.Category = strings.Split(categoriesStr, ", ")

		posts = append(posts, post)
	}

	return posts, nil
}
