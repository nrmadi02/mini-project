package http

import (
	"github.com/labstack/echo/v4"
	"github.com/nrmadi02/mini-project/domain"
	"github.com/nrmadi02/mini-project/web/request"
	"github.com/nrmadi02/mini-project/web/response"
	"net/http"
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

// Register godoc
// @Summary Register new user
// @Description Register for create new user
// @Tags Auth
// @param data body request.UserCreateRequest true "required"
// @accept json
// @Produce json
// @Router /register [post]
// @Success 201 {object} response.JSONSuccessResult{data=response.UserCreateResponse}
// @Failure 400 {object} response.JSONBadRequestResult{}
func (a authController) Register(c echo.Context) error {
	var req request.UserCreateRequest

	if err := c.Bind(&req); err != nil {
		return response.FailResponse(c, http.StatusBadRequest, false, err.Error())
	}
	if val, err := request.ValidateCreation(req); val == false {
		return response.FailResponse(c, http.StatusBadRequest, false, err.Error())
	}

	createdUser, err := a.AuthUsecase.Register(req)

	if err != nil {
		return response.FailResponse(c, http.StatusBadRequest, false, err.Error())
	}

	res := response.UserCreateResponse{
		ID:        createdUser.ID,
		Email:     createdUser.Email,
		Fullname:  createdUser.Fullname,
		Username:  createdUser.Username,
		CreatedAt: createdUser.CreatedAt,
		UpdatedAt: createdUser.UpdatedAt,
	}

	return response.SuccessResponse(c, http.StatusCreated, true, "success create new user", res)

}

// Login godoc
// @Summary Login user
// @Description Login for get JWT token
// @Tags Auth
// @param data body request.LoginRequest true "required"
// @accept json
// @Produce json
// @Router /login [post]
// @Success 201 {object} response.JSONSuccessResult{data=response.SuccessLogin}
// @Failure 400 {object} response.JSONBadRequestResult{}
// @Failure 401 {object} response.JSONUnauthorizedResult{}
func (a authController) Login(c echo.Context) error {
	var req request.LoginRequest

	if err := c.Bind(&req); err != nil {
		return response.FailResponse(c, http.StatusBadRequest, false, err.Error())
	}

	res, err := a.AuthUsecase.Login(req)
	if err != nil {
		return response.FailResponse(c, http.StatusUnauthorized, false, err.Error())
	}

	return response.SuccessResponse(c, http.StatusOK, true, "login success", res)

}
