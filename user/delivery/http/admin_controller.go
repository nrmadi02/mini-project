package http

import (
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/nrmadi02/mini-project/domain"
	"github.com/nrmadi02/mini-project/web/response"
	"net/http"
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

// GetUserList godoc
// @Summary Get list users
// @Description Get list users can access only admin
// @Tags Admin
// @accept json
// @Produce json
// @Router /admin/users [get]
// @Success 200 {object} response.JSONSuccessResult{data=[]response.UsersListResponse}
// @Failure 400 {object} response.JSONBadRequestResult{}
// @Failure 401 {object} response.JSONUnauthorizedResult{}
// @Security JWT
func (a adminController) GetUserList(c echo.Context) error {

	jwtBearer := c.Get("user").(*jwt.Token)
	claims := jwtBearer.Claims.(jwt.MapClaims)

	isAdmin, err := a.AuthUsecase.CheckIfUserIsAdmin(claims["UserID"].(string))
	if err != nil {
		return response.FailResponse(c, http.StatusBadRequest, false, err.Error())
	}

	if !isAdmin {
		return response.FailResponse(c, http.StatusUnauthorized, false, "only access admin")
	}

	foundUsers, err := a.UserUsecase.GetAllUsers()
	if err != nil {
		return response.FailResponse(c, http.StatusBadRequest, false, err.Error())
	}

	var res []response.UsersListResponse
	for _, foundUser := range foundUsers {
		res = append(res, response.UsersListResponse{
			ID:       foundUser.ID,
			Email:    foundUser.Email,
			Fullname: foundUser.Fullname,
			Username: foundUser.Username,
		})
	}

	return response.SuccessResponse(c, http.StatusOK, true, "success get list users", res)
}
