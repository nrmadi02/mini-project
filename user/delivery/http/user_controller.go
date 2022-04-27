package http

import (
	"github.com/labstack/echo/v4"
	"github.com/nrmadi02/mini-project/domain"
)

type UserController interface {
	User(c echo.Context) error
}

type userController struct {
	AuthUsecase domain.AuthUsecase
}

func NewUserController(au domain.AuthUsecase) UserController {
	return userController{
		AuthUsecase: au,
	}
}

func (u userController) User(c echo.Context) error {
	//TODO implement me
	panic("implement me")
}
