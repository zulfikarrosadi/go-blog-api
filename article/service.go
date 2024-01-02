package article

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/zulfikarrosadi/go-blog-api/lib"
	"github.com/zulfikarrosadi/go-blog-api/web"
)

type Error struct {
	Message string `json:"message"`
	Detail  any    `json:"details"`
}

type ErrorDetail struct {
	Path    string `json:"path"`
	Value   string `json:"value"`
	Message string `json:"message"`
}

type ArticleService interface {
	GetArticles(context.Context) web.Response
	FindArticleById(int, context.Context) web.Response
	CreateArticle(*CreateArticleRequest, context.Context) web.Response
	DeleteArticleById(int, context.Context) web.Response
}

type ArticleServiceImpl struct {
	ArticleRepository
	v *validator.Validate
}

func NewArticleService(
	articleRepository ArticleRepository, v *validator.Validate,
) *ArticleServiceImpl {
	return &ArticleServiceImpl{
		ArticleRepository: articleRepository,
		v:                 v,
	}
}

func (as *ArticleServiceImpl) GetArticles(ctx context.Context) web.Response {
	articlesChannel := make(chan []Article)
	errorChannel := make(chan error)
	defer close(articlesChannel)
	defer close(errorChannel)

	go func() {
		articles, err := as.ArticleRepository.GetArticles(ctx)
		if err != nil {
			errorChannel <- err
			return
		}
		articlesChannel <- articles
	}()

	select {
	case result := <-articlesChannel:
		fmt.Println(result)
		response := &web.Response{
			Status: "success",
			Code:   200,
			Data:   result,
		}
		return *response
	case result := <-errorChannel:
		fmt.Println(result)
		response := &web.Response{
			Status: "fail",
			Code:   400,
			Data:   nil,
			Error: web.Error{
				Message: result.Error(),
			},
		}
		return *response
	}
}

func (as *ArticleServiceImpl) FindArticleById(id int, ctx context.Context) web.Response {
	articleChannel := make(chan *Article)
	defer close(articleChannel)

	go func() {
		articleChannel <- as.ArticleRepository.FindArticleById(id, ctx)
	}()

	article := <-articleChannel
	if article == nil {
		return web.Response{
			Status: "fail",
			Code:   http.StatusNotFound,
			Data:   nil,
		}
	}
	return web.Response{
		Status: "success",
		Code:   http.StatusOK,
		Data:   article,
	}
}

func (as *ArticleServiceImpl) CreateArticle(data *CreateArticleRequest, ctx context.Context) web.Response {
	err := as.v.Struct(data)
	if err != nil {
		validatedError := lib.ValidateError(err.(validator.ValidationErrors))
		return web.Response{
			Status: "fail",
			Code:   http.StatusBadRequest,
			Error: web.Error{
				Message: "validation error",
				Detail:  validatedError,
			},
		}
	}

	errorChannel := make(chan error)
	articleIdChannel := make(chan int64)
	defer close(errorChannel)
	defer close(articleIdChannel)

	go func() {
		id, err := as.ArticleRepository.CreateArticle(data, ctx)
		errorChannel <- err
		articleIdChannel <- id
	}()
	err = <-errorChannel
	if err != nil {
		fmt.Println(err)
		return web.Response{
			Status: "fail",
			Code:   http.StatusBadRequest,
			Error: web.Error{
				Message: "cannot create article, please try again",
				Detail: []ErrorDetail{{
					Path:  "title",
					Value: data.Title,
				}, {
					Path:  "content",
					Value: data.Content,
				}},
			},
		}
	}
	return web.Response{
		Status: "success",
		Code:   http.StatusCreated,
		Data: struct {
			Id int64 `json:"id"`
		}{Id: <-articleIdChannel},
	}
}

func (as *ArticleServiceImpl) DeleteArticleById(id int, ctx context.Context) web.Response {
	errorChannel := make(chan error)
	defer close(errorChannel)

	go func() {
		errorChannel <- as.ArticleRepository.DeleteArticleById(id, ctx)
	}()
	err := <-errorChannel
	if err != nil {
		return web.Response{
			Status: "fail",
			Code:   http.StatusNotFound,
			Error: web.Error{
				Message: "cannot delete article, please try again",
			},
		}
	}
	return web.Response{
		Status: "success",
		Code:   http.StatusNoContent,
	}
}
