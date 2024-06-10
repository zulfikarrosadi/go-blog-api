package article

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/zulfikarrosadi/go-blog-api/auth"
	"github.com/zulfikarrosadi/go-blog-api/web"
)

type ArticleApi interface {
	GetArticles(echo.Context) error
	GetArticleById(echo.Context) error
	CreateArticle(echo.Context) error
	UpdateArticle(echo.Context) error
	DeleteArticle(echo.Context) error
}

type ArticleApiImpl struct {
	ArticleServiceImpl ArticleService
}

func NewArticleApi(articleService ArticleService) *ArticleApiImpl {
	return &ArticleApiImpl{
		ArticleServiceImpl: articleService,
	}
}

func (aa *ArticleApiImpl) GetArticles(c echo.Context) error {
	r := aa.ArticleServiceImpl.GetArticles(c.Request().Context())
	return c.JSON(r.Code, r)
}

func (aa *ArticleApiImpl) GetArticleById(c echo.Context) error {
	slug := c.Param("slug")

	if trimed := strings.Trim(slug, " "); len(trimed) < 1 || slug == "" {
		fmt.Println("slug is nil")
		return c.JSON(http.StatusNotFound, web.Response{
			Status: "fail",
			Code:   http.StatusNotFound,
			Data:   nil,
		})
	}

	r := aa.ArticleServiceImpl.FindArticleById(slug, c.Request().Context())

	return c.JSON(r.Code, r)
}

func (aa *ArticleApiImpl) CreateArticle(c echo.Context) error {
	articleRequest := &CreateArticleRequest{}
	c.Bind(&articleRequest)
	articleRequest.CreatedAt = time.Now().Unix()

	accessToken := c.Get("accessToken").(auth.AccessToken)
	ctx := context.WithValue(c.Request().Context(), "accessToken", accessToken)

	r := aa.ArticleServiceImpl.CreateArticle(articleRequest, ctx)
	return c.JSON(r.Code, r)
}

func (aa *ArticleApiImpl) DeleteArticle(c echo.Context) error {
	articleRequest := &ArticleRequest{}
	c.Bind(articleRequest)
	id, _ := strconv.Atoi(articleRequest.Id)
	r := aa.ArticleServiceImpl.DeleteArticleById(id, c.Request().Context())
	return c.JSON(r.Code, r)
}

func (aa *ArticleApiImpl) UpdateArticle(c echo.Context) error {
	return c.String(http.StatusNotImplemented, "NOT IMPLEMENTED")
}
