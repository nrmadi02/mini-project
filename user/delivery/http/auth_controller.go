package http

import (
	"github.com/labstack/echo/v4"
	"github.com/nrmadi02/mini-project/domain"
)

type AuthController interface {
	Register(c echo.Context) error
	Login(c echo.Context) error
}

type authController struct {
	AuthUsecase domain.AuthUsecase
}

func NewAuthController(au domain.AuthUsecase) AuthController {
	return authController{
		AuthUsecase: au,
	}
}

func (a authController) Register(c echo.Context) error {
	//TODO implement me
	panic("implement me")
}

func (a authController) Login(c echo.Context) error {
	//TODO implement me
	panic("implement me")
}
