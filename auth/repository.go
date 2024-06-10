package auth

import (
	"context"
	"database/sql"
	"errors"

	"github.com/zulfikarrosadi/go-blog-api/lib"
)

type AuthRepository interface {
	FindUserByUsername(*UserSignInRequest, context.Context) (*User, error)
	CreateUser(*UserSignUpRequest, context.Context) (*UserAuthResponse, error)
}

type AuthRepositoryImpl struct {
	*sql.DB
}

func NewAuthRepository(connection *sql.DB) *AuthRepositoryImpl {
	return &AuthRepositoryImpl{
		DB: connection,
	}
}

func (as *AuthRepositoryImpl) CreateUser(
	data *UserSignUpRequest, ctx context.Context,
) (*UserAuthResponse, error) {
	q := "INSERT INTO users (username, password) VALUES (?,?)"
	r, err := as.DB.ExecContext(ctx, q, data.Username, data.Password)
	if err != nil {
		lib.ValidateErrorV2("craete_user_repo", err)
		return nil, errors.New("this username is already in use. please use a different username or try logging in")
	}
	i, _ := r.LastInsertId()

	return &UserAuthResponse{
		UserId:   i,
		Username: data.Username,
	}, nil
}

func (as *AuthRepositoryImpl) FindUserByUsername(data *UserSignInRequest, ctx context.Context) (*User, error) {
	q := "SELECT id, username, password FROM users WHERE username = ?"
	r := as.DB.QueryRowContext(ctx, q, data.Username)
	user := &User{}
	err := r.Scan(&user.Id, &user.Username, &user.Password)
	if err != nil {
		lib.ValidateErrorV2("find_user_by_username_repo", err)
		return nil, errors.New("username or password is incorrect")
	}
	return user, nil
}
