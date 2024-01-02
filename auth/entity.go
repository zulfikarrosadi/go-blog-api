package auth

import "time"

type User struct {
	Id        int64     `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"createdAt"`
}

type UserSignUpRequest struct {
	Username             string `json:"username" validate:"required"`
	Password             string `json:"password" validate:"required"`
	PasswordConfirmation string `json:"passwordConfirmation" validate:"eqfield=Password"`
}

type UserSignInRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type UserAuthResponse struct {
	UserId   int64  `json:"user_id"`
	Username string `json:"username"`
}
