package article

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"

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
	FindArticleById(string, context.Context) web.Response
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

func (as *ArticleServiceImpl) FindArticleById(slug string, ctx context.Context) web.Response {
	articleChannel := make(chan *Article)
	errorChannel := make(chan error)
	defer close(articleChannel)
	defer close(errorChannel)

	timestamp, err := extractTimestampFromSlug(slug)
	if err != nil {
		return web.Response{
			Status: "fail",
			Code:   http.StatusNotFound,
			Error: web.Error{
				Message: "article not found",
			},
			Data: nil,
		}
	}

	go func() {
		article, err := as.ArticleRepository.FindArticleById(timestamp, ctx)
		if err != nil {
			errorChannel <- err
			return
		}
		articleChannel <- article
	}()

	select {
	case result := <-errorChannel:
		if newErr := result.(*net.OpError); newErr != nil {
			return web.Response{
				Status: "fail",
				Code:   http.StatusInternalServerError,
				Error: web.Error{
					Message: "something went wrong, please wait and try again",
				},
				Data: nil,
			}
		}
		return web.Response{
			Status: "fail",
			Code:   http.StatusNotFound,
			Error: web.Error{
				Message: "article not found",
			},
			Data: nil,
		}
	case result := <-articleChannel:
		return web.Response{
			Status: "success",
			Code:   http.StatusOK,
			Data:   result,
		}
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
	data.Slug = createSlug(data.Title, data.CreatedAt)

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
			Id   int64  `json:"id"`
			Slug string `json:"slug"`
		}{Id: <-articleIdChannel, Slug: data.Slug},
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

func createSlug(title string, timestamp int64) string {
	splitedTitle := strings.Split(strings.Trim(title, " "), " ")
	slug := strings.ToLower(strings.Join(splitedTitle, "-"))
	stringifyTimestamp := strconv.FormatInt(timestamp, 10)
	slug = slug + "-" + stringifyTimestamp

	return slug
}

func extractTimestampFromSlug(slug string) (int64, error) {
	splitedSlug := strings.Split(slug, "-")
	i, err := strconv.ParseInt(splitedSlug[len(splitedSlug)-1], 10, 64)
	if err != nil {
		return 0, err
	}

	return i, nil
}
