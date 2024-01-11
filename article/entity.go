package article

import (
	"database/sql"
	"time"
)

type Article struct {
	Id        int            `json:"id"`
	Title     string         `json:"title"`
	Content   sql.NullString `json:"content"`
	Author    int            `json:"author"`
	CreatedAt time.Time      `json:"created_at"`
}

type ArticleRequest struct {
	Id      string `param:"id"`
	Title   string `json:"title" validate:"required"`
	Content string `json:"content"`
}

type CreateArticleRequest struct {
	Title   string `json:"title" validate:"required"`
	Content string `json:"content"`
	Slug    string
}

type UpdateArticleRequest struct {
	Id      string `param:"id"`
	Title   string `json:"title" validate:"required"`
	Content string `json:"content"`
	Slug    string
	Author  int
}
