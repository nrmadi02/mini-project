package domain

import (
	request2 "github.com/nrmadi02/mini-project/web/request"
	"github.com/nrmadi02/mini-project/web/response"
	uuid "github.com/satori/go.uuid"
	"time"
)

type User struct {
	ID               uuid.UUID          `json:"id" gorm:"PrimaryKey"`
	Fullname         string             `json:"fullname" gorm:"notnull"`
	Email            string             `json:"email" gorm:"notnull"`
	Username         string             `json:"username" gorm:"unique;notnull"`
	Password         string             `json:"password" gorm:"notnull"`
	Roles            []Role             `json:"roles" gorm:"many2many:user_roles;"`
	Enterprises      []Enterprise       `json:"enterprises" gorm:"foreignKey:UserID;references:ID"`
	RatingEnterprise []RatingEnterprise `json:"rating_enterprise" gorm:"foreignKey:UserID;references:ID"`
	Favorite         Favorite           `json:"favorite" gorm:"foreignKey:UserID;references:ID"`
	CreatedAt        time.Time          `json:"created_at"`
	UpdatedAt        time.Time          `json:"updated_at"`
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
	Login(request request2.LoginRequest) (response.SuccessLogin, error)
	Register(request request2.UserCreateRequest) (User, error)
	GetUserDetails(id string) (User, Favorite, Enterprises, error)
	CheckIfUserIsAdmin(id string) (bool, error)
}
