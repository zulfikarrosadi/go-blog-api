package main

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
	"github.com/zulfikarrosadi/go-blog-api/article"
)

func main() {
	e := echo.New()
	validator := validator.New()
	articleRepository := article.NewArticleRepository(GetDBConnection())
	articleService := article.NewArticleService(articleRepository, validator)

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello World")
	})
	e.GET("/api/articles", func(c echo.Context) error {
		r := articleService.GetArticles(c.Request().Context())
		return c.JSON(200, r)
	})
	e.POST("/api/articles", func(c echo.Context) error {
		articleRequest := article.ArticleRequest{}
		c.Bind(&articleRequest)
		r := articleService.CreateArticle(
			&article.Article{
				Title: articleRequest.Title,
				Content: sql.NullString{
					String: articleRequest.Content,
					Valid:  len(articleRequest.Content) != 0,
				},
			},
			c.Request().Context(),
		)
		return c.JSON(r.Code, r)
	})
	e.DELETE("/api/articles/:id", func(c echo.Context) error {
		articleRequest := new(article.ArticleRequest)
		c.Bind(articleRequest)
		id, _ := strconv.Atoi(articleRequest.Id)
		r := articleService.DeleteArticleById(id, c.Request().Context())
		return c.JSON(r.Code, r)
	})
	e.GET("/api/articles/:id", func(c echo.Context) error {
		articleRequest := new(article.ArticleRequest)
		c.Bind(articleRequest)
		id, _ := strconv.Atoi(articleRequest.Id)
		r := articleService.FindArticleById(id, c.Request().Context())
		if r.Code == http.StatusNotFound {
			return c.JSON(r.Code, r)
		}
		return c.JSON(r.Code, r)
	})

	e.Logger.Fatal(e.Start("localhost:3000"))
}

func GetDBConnection() *sql.DB {
	dsn := "root:@tcp(localhost:3306)/golang_article?parseTime=true"
	d, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	d.SetMaxOpenConns(6)
	d.SetMaxIdleConns(2)
	return d
}
