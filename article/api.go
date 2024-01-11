package article

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

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
	id, err := strconv.Atoi(c.QueryParam("id"))
	if err != nil {
		return c.JSON(http.StatusNotFound, web.Response{
			Status: "fail",
			Code:   http.StatusNotFound,
			Data:   nil,
		})
	}
	s := c.Param("slug")
	if trimed := strings.Trim(s, " "); len(trimed) < 1 || s == "" {
		fmt.Println("slug is nil")
		return c.JSON(http.StatusNotFound, web.Response{
			Status: "fail",
			Code:   http.StatusNotFound,
			Data:   nil,
		})
	}

	r := aa.ArticleServiceImpl.FindArticleById(id, c.Request().Context())

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

func (aa *ArticleApiImpl) GetUserLoginInfo(c echo.Context) context.Context {
	accessToken := c.Get("accessToken").(auth.AccessToken)
	ctx := context.WithValue(c.Request().Context(), "accessToken", accessToken)

	return ctx
}
