package http

import (
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/nrmadi02/mini-project/domain"
	"github.com/nrmadi02/mini-project/web/response"
	"net/http"
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

// User godoc
// @Summary Get detail user by JWT Token
// @Description User id get default by claims JWT Token
// @Tags User
// @accept json
// @Produce json
// @Router /user [get]
// @Success 200 {object} response.JSONSuccessResult{data=response.UserDetailResponse}
// @Failure 401 {object} response.JSONUnauthorizedResult{}
// @Security JWT
func (u userController) User(c echo.Context) error {
	jwtBearer := c.Get("user").(*jwt.Token)
	claims := jwtBearer.Claims.(jwt.MapClaims)

	details, err := u.AuthUsecase.GetUserDetails(claims["UserID"].(string))
	if err != nil {
		return response.FailResponse(c, http.StatusUnauthorized, false, err.Error())
	}

	res := response.UserDetailResponse{
		Fullname:  details.Fullname,
		Username:  details.Username,
		Email:     details.Email,
		ID:        details.ID,
		CreatedAt: details.CreatedAt,
		UpdatedAt: details.UpdatedAt,
	}

	return response.SuccessResponse(c, http.StatusOK, true, "success get detail user", res)

}
