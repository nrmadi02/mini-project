package http

import (
	"github.com/labstack/echo/v4"
	"github.com/nrmadi02/mini-project/domain"
)

type AdminController interface {
	GetUserList(c echo.Context) error
}

type adminController struct {
	AuthUsecase domain.AuthUsecase
	UserUsecase domain.UserUsecase
}

func NewAdminController(au domain.AuthUsecase, uu domain.UserUsecase) AdminController {
	return adminController{
		AuthUsecase: au,
		UserUsecase: uu,
	}
}

func (a adminController) GetUserList(c echo.Context) error {
	//TODO implement me
	panic("implement me")
}
