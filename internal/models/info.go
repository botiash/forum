package models

type Info struct {
	User
	Post
	Comment []Comment
}

type InfoPosts struct {
	User
	Posts    []Post
	Category []string
}
