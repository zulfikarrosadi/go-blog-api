package auth

import (
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type Auth interface {
	AuthenticationRequired(next echo.HandlerFunc) echo.HandlerFunc
	DeserializeUser(next echo.HandlerFunc) echo.HandlerFunc
}

type AuthMiddleware struct{}

func NewAuthMiddleware() AuthMiddleware {
	return AuthMiddleware{}
}

func (am *AuthMiddleware) AuthenticationRequired(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Get("accessToken")
		if token == nil {
			return c.NoContent(http.StatusUnauthorized)
		}
		c.Set("user", token)
		return next(c)
	}
}

func (am *AuthMiddleware) DeserializeUser(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		fmt.Println("deserialize")
		accessTokenCookie, err := c.Cookie("accessToken")
		if err != nil {
			fmt.Println(err)
			return next(c)
		}
		fmt.Println("raw accesstoken cookie string", accessTokenCookie)
		fmt.Println("cookie value", accessTokenCookie.Value)

		accessToken := AccessToken{}
		token, err := jwt.ParseWithClaims(accessTokenCookie.Value, &accessToken, func(t *jwt.Token) (interface{}, error) {
			return []byte("temp key"), nil
		})
		fmt.Println("decoded token aid", accessToken.AccessTokenId)
		fmt.Println("decoded token username", accessToken.Username)

		fmt.Println("decoded token", token)
		if err != nil {
			fmt.Println(err)
			return next(c)
		}
		fmt.Println("final access token", accessToken)
		c.Set("accessToken", accessToken)
		return next(c)
	}
}
