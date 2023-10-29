package article

import (
	"context"
	"fmt"
	"net/http"
)

type Response struct {
	Status string `json:"status"`
	Code   int    `json:"code"`
	Data   any    `json:"data"`
	Error  Error  `json:"errors"`
}

type Error struct {
	Message string        `json:"message"`
	Detail  []ErrorDetail `json:"details"`
}

type ErrorDetail struct {
	Path  string `json:"path"`
	Value string `json:"value"`
}

type ArticleService interface {
	GetArticles(context.Context) Response
	FindArticleById(int, context.Context) Response
	CreateArticle(Article, context.Context) Response
	DeleteArticleById(int, context.Context) Response
}

type ArticleServiceImpl struct {
	ArticleRepository
}

func NewArticleService(articleRepository ArticleRepository) *ArticleServiceImpl {
	return &ArticleServiceImpl{
		ArticleRepository: articleRepository,
	}
}

func (as *ArticleServiceImpl) GetArticles(ctx context.Context) Response {
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
		response := &Response{
			Status: "success",
			Code:   200,
			Data:   result,
			Error:  Error{},
		}
		return *response
	case result := <-errorChannel:
		fmt.Println(result)
		response := &Response{
			Status: "fail",
			Code:   400,
			Data:   nil,
			Error: Error{
				Message: result.Error(),
			},
		}
		return *response
	}
}

func (as *ArticleServiceImpl) FindArticleById(id int, ctx context.Context) Response {
	articleChannel := make(chan *Article)
	defer close(articleChannel)

	go func() {
		articleChannel <- as.ArticleRepository.FindArticleById(id, ctx)
	}()

	article := <-articleChannel
	if article == nil {
		return Response{
			Status: "fail",
			Code:   http.StatusNotFound,
			Data:   nil,
		}
	}
	return Response{
		Status: "success",
		Code:   http.StatusOK,
		Data:   article,
	}
}

func (as *ArticleServiceImpl) CreateArticle(data *Article, ctx context.Context) Response {
	errorChannel := make(chan error)
	articleIdChannel := make(chan int64)
	defer close(errorChannel)
	defer close(articleIdChannel)

	go func() {
		id, err := as.ArticleRepository.CreateArticle(data, ctx)
		errorChannel <- err
		articleIdChannel <- id
	}()
	err := <-errorChannel
	if err != nil {
		return Response{
			Status: "fail",
			Code:   http.StatusBadRequest,
			Error: Error{
				Message: "cannot create article, please try again",
				Detail: []ErrorDetail{{
					Path:  "title",
					Value: data.Title,
				}, {
					Path:  "content",
					Value: data.Content.String,
				}},
			},
		}
	}
	return Response{
		Status: "success",
		Code:   http.StatusCreated,
		Data: struct {
			Id int64 `json:"id"`
		}{Id: <-articleIdChannel},
	}
}

func (as *ArticleServiceImpl) DeleteArticleById(id int, ctx context.Context) Response {
	errorChannel := make(chan error)
	defer close(errorChannel)

	go func() {
		errorChannel <- as.ArticleRepository.DeleteArticleById(id, ctx)
	}()
	err := <-errorChannel
	if err != nil {
		return Response{
			Status: "fail",
			Code:   http.StatusNotFound,
			Error: Error{
				Message: "cannot delete article, please try again",
			},
		}
	}
	return Response{
		Status: "success",
		Code:   http.StatusNoContent,
	}
}
