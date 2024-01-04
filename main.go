package main

import (
	"database/sql"
	"log"

	"github.com/go-playground/validator/v10"
	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
	"github.com/zulfikarrosadi/go-blog-api/article"
	"github.com/zulfikarrosadi/go-blog-api/auth"
	"github.com/zulfikarrosadi/go-blog-api/lib"
)

func main() {
	e := echo.New()
	validator := validator.New()
	db := GetDBConnection()
	articleRepository := article.NewArticleRepository(GetDBConnection())
	articleService := article.NewArticleService(articleRepository, validator)
	articleHandler := article.NewArticleApi(articleService)

	authRepository := auth.NewAuthRepository(db)
	authService := auth.NewAuthService(authRepository, validator)
	authHandler := auth.NewAuthHandler(authService)
	authMiddleware := auth.NewAuthMiddleware()

	e.POST("/api/signin", authHandler.SignInHandler)
	e.POST("/api/signup", authHandler.SignUpHandler)
	e.POST("/api/refresh", authHandler.RefreshTokenHandler)

	protectedRouteGroup := e.Group("/api/auth")
	protectedRouteGroup.Use(authMiddleware.DeserializeUser)
	protectedRouteGroup.Use(authMiddleware.AuthenticationRequired)

	e.GET("/api/articles", articleHandler.GetArticles)
	e.GET("/api/articles/:id", articleHandler.GetArticleById)
	protectedRouteGroup.POST("/articles", articleHandler.CreateArticle)
	protectedRouteGroup.DELETE("/articles/:id", articleHandler.DeleteArticle)
	protectedRouteGroup.POST("/files", lib.FileUploadHandler)

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
