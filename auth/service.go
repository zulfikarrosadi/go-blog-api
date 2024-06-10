package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/zulfikarrosadi/go-blog-api/lib"
	"github.com/zulfikarrosadi/go-blog-api/web"
	"golang.org/x/crypto/bcrypt"
)

const FIFTEEN_DAY_IN_HOUR = 360
const BCRYPT_COST = 10

type AuthService interface {
	SignIn(*UserSignInRequest, context.Context) (*AccessToken, *RefreshToken, *web.Response)
	SignUp(*UserSignUpRequest, context.Context) (*AccessToken, *RefreshToken, *web.Response)
}

type AuthServiceImpl struct {
	AuthRepository
	v *validator.Validate
}

type AccessToken struct {
	AccessTokenId string `json:"accessTokenId"`
	UserId        int64  `json:"id"`
	Username      string `json:"username"`
	jwt.RegisteredClaims
}

type RefreshToken struct {
	RefreshTokenId string `json:"refreshTokenId"`
	Id             int64  `json:"id"`
	Username       string `json:"username"`
	jwt.RegisteredClaims
}

func NewAuthService(authRepository AuthRepository, v *validator.Validate) *AuthServiceImpl {
	return &AuthServiceImpl{
		AuthRepository: authRepository,
		v:              v,
	}
}

func (asi *AuthServiceImpl) SignIn(data *UserSignInRequest, ctx context.Context) (*AccessToken, *RefreshToken, *web.Response) {
	err := asi.v.Struct(data)
	if err != nil {
		validatedError := lib.ValidateError(err.(validator.ValidationErrors))
		return nil, nil, &web.Response{
			Status: "fail",
			Code:   http.StatusBadRequest,
			Error: web.Error{
				Message: "validation error",
				Detail:  validatedError,
			},
		}
	}

	user, err := asi.AuthRepository.FindUserByUsername(data, ctx)
	if err != nil {
		response := &web.Response{
			Status: web.STATUS_FAIL,
			Code:   http.StatusBadRequest,
			Error: web.Error{
				Message: err.Error(),
			},
		}
		return nil, nil, response
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.Password))
	if err != nil {
		response := &web.Response{
			Status: web.STATUS_FAIL,
			Code:   http.StatusBadRequest,
			Error: web.Error{
				Message: err.Error(),
			},
		}
		return nil, nil, response
	}

	accessTokenId := make([]byte, 10)
	refreshTokenId := make([]byte, 10)
	rand.Read(accessTokenId)
	rand.Read(refreshTokenId)

	accessTokenClaims := &AccessToken{
		AccessTokenId: base64.URLEncoding.EncodeToString(accessTokenId)[:10],
		UserId:        user.Id,
		Username:      user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 15)),
		},
	}
	refreshTokenClaims := &RefreshToken{
		RefreshTokenId: base64.URLEncoding.EncodeToString(refreshTokenId)[:10],
		Id:             user.Id,
		Username:       user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * FIFTEEN_DAY_IN_HOUR)),
		},
	}
	return accessTokenClaims, refreshTokenClaims, nil
}

func (asi *AuthServiceImpl) SignUp(data *UserSignUpRequest, ctx context.Context) (*AccessToken, *RefreshToken, *web.Response) {
	err := asi.v.Struct(data)
	if err != nil {
		validatedError := lib.ValidateError(err.(validator.ValidationErrors))
		return nil, nil, &web.Response{
			Status: "fail",
			Code:   http.StatusBadRequest,
			Error: web.Error{
				Message: "validation error",
				Detail:  validatedError,
			},
		}
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(data.Password), BCRYPT_COST)
	data.Password = string(hashedPassword)
	user, err := asi.AuthRepository.CreateUser(data, ctx)
	if err != nil {
		response := &web.Response{
			Status: web.STATUS_FAIL,
			Code:   http.StatusBadRequest,
			Error: web.Error{
				Message: err.Error(),
			},
		}
		return nil, nil, response
	}
	accessTokenId := make([]byte, 10)
	refreshTokenId := make([]byte, 10)
	rand.Read(accessTokenId)
	rand.Read(refreshTokenId)

	accessTokenClaims := &AccessToken{
		AccessTokenId: base64.URLEncoding.EncodeToString(accessTokenId)[:10],
		UserId:        user.UserId,
		Username:      user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 15)),
		},
	}
	refreshTokenClaims := &RefreshToken{
		RefreshTokenId: base64.URLEncoding.EncodeToString(refreshTokenId)[:10],
		Id:             user.UserId,
		Username:       user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * FIFTEEN_DAY_IN_HOUR)),
		},
	}
	return accessTokenClaims, refreshTokenClaims, nil
}

func CreateToken(newRefreshToken bool, claims ...jwt.Claims) []string {
	accessTokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims[0])
	accessTokenString, _ := accessTokenClaims.SignedString([]byte("temp key"))

	if !newRefreshToken {
		return []string{accessTokenString}
	}
	refreshTokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims[1])
	refreshTokenString, _ := refreshTokenClaims.SignedString([]byte("temp key"))

	return []string{accessTokenString, refreshTokenString}
}

func ValidateToken(token string, isRefreshToken bool) (*AccessToken, *RefreshToken, error) {
	if isRefreshToken {
		refreshToken := &RefreshToken{}
		_, err := jwt.ParseWithClaims(token, refreshToken, func(t *jwt.Token) (interface{}, error) {
			return []byte("temp key"), nil
		})
		if err != nil {
			return nil, nil, errors.New("refresh token invalid")
		}
		return nil, refreshToken, nil
	}

	accessToken := &AccessToken{}
	_, err := jwt.ParseWithClaims(token, accessToken, func(t *jwt.Token) (interface{}, error) {
		return []byte("temp key"), nil
	})
	if err != nil {
		return nil, nil, errors.New("access token invalid")
	}
	return accessToken, nil, nil
}
