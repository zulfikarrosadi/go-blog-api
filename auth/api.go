package auth

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type AuthHandler interface {
	SignInHandler(echo.Context) error
	SignUpHandler(echo.Context) error
	RefreshTokenHandler(echo.Context) error
}

type AuthHandlerImpl struct {
	AuthService
}

func NewAuthHandler(authService AuthService) *AuthHandlerImpl {
	return &AuthHandlerImpl{
		AuthService: authService,
	}
}

func (ahi *AuthHandlerImpl) SignInHandler(c echo.Context) error {
	data := &UserSignInRequest{}
	c.Bind(data)
	accessTokenClaims, refreshTokenClaims, errorResponse :=
		ahi.AuthService.SignIn(data, c.Request().Context())

	if errorResponse != nil {
		return c.JSON(errorResponse.Code, errorResponse)
	}

	tokens := CreateToken(true, accessTokenClaims, refreshTokenClaims)

	accessTokenCookie := &http.Cookie{
		Name:     "accessToken",
		Value:    tokens[0],
		Secure:   false,
		HttpOnly: true,
	}
	refreshTokenCookie := &http.Cookie{
		Name:     "refreshToken",
		Value:    tokens[1],
		Secure:   false,
		HttpOnly: true,
		Path:     "/api/refresh",
	}

	c.SetCookie(accessTokenCookie)
	c.SetCookie(refreshTokenCookie)
	return c.String(http.StatusOK, http.StatusText(http.StatusOK))
}

func (ahi *AuthHandlerImpl) SignUpHandler(c echo.Context) error {
	data := &UserSignUpRequest{}
	c.Bind(data)
	fmt.Println(data)
	accessTokenClaims, refreshTokenClaims, errorResponse :=
		ahi.AuthService.SignUp(data, c.Request().Context())

	if errorResponse != nil {
		return c.JSON(errorResponse.Code, errorResponse)
	}

	tokens := CreateToken(true, accessTokenClaims, refreshTokenClaims)
	accessTokenCookie := &http.Cookie{
		Name:     "accessToken",
		Value:    tokens[0],
		Secure:   false,
		SameSite: http.SameSiteNoneMode,
		HttpOnly: true,
	}
	refreshTokenCookie := &http.Cookie{
		Name:     "refreshToken",
		Value:    tokens[1],
		Secure:   false,
		HttpOnly: true,
		Path:     "/api/refresh",
	}

	c.SetCookie(accessTokenCookie)
	c.SetCookie(refreshTokenCookie)
	return c.String(http.StatusOK, http.StatusText(http.StatusOK))
}

func (ahi *AuthHandlerImpl) RefreshTokenHandler(c echo.Context) error {
	refreshTokenCookie, err := c.Request().Cookie("refreshToken")
	if err != nil {
		fmt.Println(err)
		return c.NoContent(http.StatusUnauthorized)
	}

	_, validatedRefreshToken, err := ValidateToken(refreshTokenCookie.Value, true)
	if err != nil {
		fmt.Println(err)
		return c.NoContent(http.StatusUnauthorized)
	}

	accessTokenId := make([]byte, 10)
	rand.Read(accessTokenId)
	newAccessTokenClaims := AccessToken{
		AccessTokenId: base64.URLEncoding.EncodeToString(accessTokenId)[:10],
		UserId:        validatedRefreshToken.Id,
		Username:      validatedRefreshToken.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 15)),
		},
	}
	tokens := CreateToken(false, newAccessTokenClaims)

	newAccessTokenCookie := http.Cookie{
		Name:     "accessToken",
		Value:    tokens[0],
		Secure:   false,
		SameSite: http.SameSiteNoneMode,
		HttpOnly: true,
	}
	c.SetCookie(&newAccessTokenCookie)
	return c.String(http.StatusOK, http.StatusText(http.StatusOK))
}
