package domain

import (
	"go_cicd/domain/web/request"
	"go_cicd/domain/web/response"
)

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Users []User

type UserRepository interface {
	Create(user *User) (*User, error)
	ReadByID(id int) (*User, error)
	ReadAll() (*Users, error)
	CheckLogin(user *User) (*User, bool, error)
}

type UserUsecase interface {
	Create(request request.UserCreateRequest) (*User, error)
	ReadByID(id int) (*User, error)
	ReadAll() (*Users, error)
	Login(request request.LoginRequest) (*response.SuccessLogin, error)
}
