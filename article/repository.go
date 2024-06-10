package article

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/zulfikarrosadi/go-blog-api/auth"
	"github.com/zulfikarrosadi/go-blog-api/lib"
)

type ArticleRepository interface {
	GetArticles(context.Context) ([]Article, error)
	FindArticleById(int64, context.Context) (*Article, error)
	CreateArticle(*CreateArticleRequest, context.Context) (int64, error)
	DeleteArticleById(int, context.Context) error
	UpdateArticleById(int, *UpdateArticleRequest, context.Context) error
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
		lib.ValidateErrorV2("get_articles_repo", err)
		return []Article{}, err
	}
	for r.Next() {
		article := Article{}
		r.Scan(&article.Id, &article.Title, &article.Content, &article.CreatedAt)
		articles = append(articles, article)
	}
	return articles, nil
}

func (as *ArticleRepositoryImpl) FindArticleById(timestamp int64, ctx context.Context) (*Article, error) {
	q := "SELECT id, title, content, created_at, author FROM articles WHERE created_at = ?"
	article := Article{}
	r := as.DB.QueryRowContext(ctx, q, timestamp)
	err := r.Scan(&article.Id, &article.Title, &article.Content, &article.CreatedAt, &article.Author)
	if err != nil {
		lib.ValidateErrorV2("find_article_by_id_repo", err)
		return nil, err
	}
	return &article, err
}

func (as *ArticleRepositoryImpl) CreateArticle(data *CreateArticleRequest, ctx context.Context) (int64, error) {
	accessToken := ctx.Value("accessToken").(auth.AccessToken)
	q := "INSERT INTO articles (title, content, author, slug, created_at) VALUES (?, ?, ?, ?, ?)"
	r, err := as.DB.ExecContext(ctx, q, data.Title, data.Content, accessToken.UserId, data.Slug, data.CreatedAt)

	if err != nil {
		lib.ValidateErrorV2("create_article_repo", err)
		return 0, err
	}

	id, _ := r.LastInsertId()
	return id, nil
}

func (as *ArticleRepositoryImpl) DeleteArticleById(id int, ctx context.Context) error {
	q := "DELETE FROM articles WHERE id = ? AND author = ?"
	user := ctx.Value("accessToken").(auth.AccessToken)
	result, err := as.DB.ExecContext(ctx, q, id, user.UserId)
	deletedArticle, _ := result.RowsAffected()

	fmt.Println("deleted article is:", deletedArticle)
	if err != nil || deletedArticle < 1 {
		fmt.Println("deleted article is:", deletedArticle)
		fmt.Println("error in deleteArticle repo", err)
		return errors.New("article not found")
	}
	return nil
}

func (as *ArticleRepositoryImpl) UpdateArticleById(id int, data *UpdateArticleRequest, ctx context.Context) error {
	user := ctx.Value("accessToken").(auth.AccessToken)
	q := "UPDATE articles SET title = ?, content = ?, slug = ? WHERE id = ? AND author = ?"

	result, err := as.DB.ExecContext(ctx, q, data.Title, data.Content, data.Slug, data.Id, user.UserId)
	updatedArticle, _ := result.RowsAffected()

	fmt.Println("updated article: ", updatedArticle)
	if err != nil || updatedArticle < 1 {
		fmt.Println("error updating article: ", err)
		return errors.New("article not found")
	}

	return nil
}
