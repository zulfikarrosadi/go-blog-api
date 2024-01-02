package article

import (
	"context"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/zulfikarrosadi/go-blog-api/auth"
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
	articleRequest := new(ArticleRequest)
	c.Bind(articleRequest)
	id, _ := strconv.Atoi(articleRequest.Id)

	r := aa.ArticleServiceImpl.FindArticleById(id, c.Request().Context())

	if r.Code == http.StatusNotFound {
		return c.JSON(r.Code, r)
	}
	return c.JSON(r.Code, r)
}

func (aa *ArticleApiImpl) CreateArticle(c echo.Context) error {
	articleRequest := &CreateArticleRequest{}
	c.Bind(&articleRequest)

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
