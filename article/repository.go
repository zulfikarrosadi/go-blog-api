package article

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/zulfikarrosadi/go-blog-api/auth"
)

type ArticleRepository interface {
	GetArticles(context.Context) ([]Article, error)
	FindArticleById(int, context.Context) *Article
	CreateArticle(*CreateArticleRequest, context.Context) (int64, error)
	DeleteArticleById(int, context.Context) error
}

type ArticleRepositoryImpl struct {
	*sql.DB
}

func NewArticleRepository(connection *sql.DB) *ArticleRepositoryImpl {
	return &ArticleRepositoryImpl{
		DB: connection,
	}
}

func (as *ArticleRepositoryImpl) GetArticles(ctx context.Context) ([]Article, error) {
	q := "SELECT id, title, content, created_at FROM articles"
	articles := []Article{}

	r, err := as.QueryContext(ctx, q)
	if err != nil {
		fmt.Println("error in getArticles repo", err)
		return []Article{}, err
	}
	for r.Next() {
		article := Article{}
		r.Scan(&article.Id, &article.Title, &article.Content, &article.CreatedAt)
		articles = append(articles, article)
	}
	return articles, nil
}

func (as *ArticleRepositoryImpl) FindArticleById(id int, ctx context.Context) *Article {
	q := "SELECT id, title, content, created_at FROM articles WHERE id = ?"
	article := Article{}
	r := as.DB.QueryRowContext(ctx, q, id)
	err := r.Scan(&article.Id, &article.Title, &article.Content, &article.CreatedAt)
	if err != nil {
		return nil
	}
	return &article
}

func (as *ArticleRepositoryImpl) CreateArticle(data *ArticleRequest, ctx context.Context) (int64, error) {
	q := "INSERT INTO articles (title, content) VALUES (?,?)"
	r, err := as.DB.ExecContext(ctx, q, data.Title, data.Content)
	if err != nil {
		return 0, err
	}

	id, _ := r.LastInsertId()
	return id, nil
}

func (as *ArticleRepositoryImpl) DeleteArticleById(id int, ctx context.Context) error {
	q := "DELETE FROM articles WHERE id = ?"
	_, err := as.DB.ExecContext(ctx, q, id)
	if err != nil {
		return err
	}
	return nil
}
