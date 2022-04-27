package domain

import (
	"github.com/nrmadi02/mini-project/domain/web/request"
	"github.com/nrmadi02/mini-project/domain/web/response"
	"time"
)

type User struct {
	ID        string    `json:"id"`
	Fullname  string    `json:"fullname" gorm:"notnull"`
	Email     string    `json:"email" gorm:"notnull"`
	Username  string    `json:"username" gorm:"unique;notnull"`
	Password  []byte    `json:"password" gorm:"notnull"`
	Roles     []Role    `json:"roles" gorm:"many2many:user_roles;"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"UpdatedAt"`
}

type Users []User

type UserRepository interface {
	FindUserByEmail(email string) (User, error)
	FindUserById(id string) (User, error)
	Save(user User) (User, error)
	FindAllUsers() (Users, error)
}

type UserUsecase interface {
	GetAllUsers() (Users, error)
}

type AuthUsecase interface {
	Login(request request.LoginRequest) (response.SuccessLogin, error)
	Register(request request.UserCreateRequest) (User, error)
	GetUserDetails(id string) (User, error)
	CheckIfUserIsAdmin(id string) (bool, error)
}
