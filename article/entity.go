package article

import (
	"database/sql"
)

type Article struct {
	Id        int            `json:"id"`
	Title     string         `json:"title"`
	Content   sql.NullString `json:"content"`
	Author    int            `json:"author"`
	CreatedAt int64          `json:"created_at"`
}

type ArticleRequest struct {
	Id      string `param:"id"`
	Title   string `json:"title" validate:"required"`
	Content string `json:"content"`
}

type CreateArticleRequest struct {
	Title     string `json:"title" validate:"required"`
	Content   string `json:"content"`
	Slug      string
	CreatedAt int64
}
