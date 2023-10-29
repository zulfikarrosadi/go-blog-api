package article

import (
	"database/sql"
	"time"
)

type Article struct {
	Id        int            `json:"id"`
	Title     string         `json:"title"`
	Content   sql.NullString `json:"content"`
	CreatedAt time.Time      `json:"created_at"`
}

type ArticleRequest struct {
	Id      string `param:"id"`
	Title   string `json:"title" validate:"required"`
	Content string `json:"content"`
}
