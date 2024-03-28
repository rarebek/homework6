package model

import (
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation/v3"
	"github.com/go-ozzo/ozzo-validation/v3/is"
)

func (r *RegisterUserRequest) Validate() error {
	return validation.ValidateStruct(
		r,
		validation.Field(&r.Name, validation.Required),
		validation.Field(&r.Email, validation.Required, is.Email),
		validation.Field(&r.Password, validation.Required, validation.Length(5, 15), validation.Match(regexp.MustCompile("[a-z]|[A-Z][0-9]"))),
	)
}

type RegisterUserRequest struct {
	Name     string `json:"name"`
	Age      int    `json:"age"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	OTP      string `json:"code"`
}

type UserModel struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Age         int    `json:"age"`
	Username    string `json:"username"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	OTP         string `json:"code"`
	AccessToken string `json:"access_token"`
}

type GetAllUserResponse struct {
	Count int64   `json:"count"`
	Users []*User `json:"users"`
}

type User struct {
	Name     string `json:"name"`
	Age      int    `json:"age"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	OTP      string `json:"code"`
}

type ListUsers struct {
	User []User
}

type RegisterUserResponse struct {
	Message string `json:"message"`
}

type LogInResponse struct {
	AccessToken string `json:"access_token"`
	Message     string `json:"message"`
}

type VerifyUserResponse struct {
	AccessToken string `json:"access_token"`
	Message     string `json:"message"`
}

type Email struct {
	Email string `json:"email"`
}

type ResponseError struct {
	Code  string      `json:"code"`
	Error interface{} `json:"error"`
}

type CreateUserRequest struct {
	Id           string `json:"id"`
	Name         string `json:"name"`
	Age          int    `json:"age"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	RefreshToken string `json:"refresh_token"`
}
