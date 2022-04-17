package http

import (
	"github.com/labstack/echo/v4"
	"go_cicd/domain"
	"go_cicd/domain/web/request"
	"go_cicd/domain/web/response"
	mid "go_cicd/user/delivery/http/middleware"
	"net/http"
	"strconv"
)

type UserController struct {
	UserUsecase domain.UserUsecase
}

func NewUserController(e *echo.Echo, Usecase domain.UserUsecase) {
	UserController := &UserController{
		UserUsecase: Usecase,
	}

	e.POST("/login", UserController.Login)
	authMiddleware := mid.NewGoMiddleware().AuthMiddleware()
	e.GET("/users/:id", UserController.GetUserByID, authMiddleware)
	e.GET("/users", UserController.GetUsers, authMiddleware)
	e.POST("/users", UserController.CreateUser)
}

func (u *UserController) Login(c echo.Context) error {
	var req request.LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	res, err := u.UserUsecase.Login(req)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"status":  false,
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"code":    200,
		"status":  true,
		"id_user": res.ID,
		"email":   res.Email,
		"token":   res.Token,
	})
}

func (u *UserController) CreateUser(c echo.Context) error {
	var req request.UserCreateRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	createdUser, err := u.UserUsecase.Create(req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"status":  false,
			"message": err.Error(),
		})
	}

	res := response.UserCreateResponse{
		ID:    int(createdUser.ID),
		Email: createdUser.Email,
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"code":   201,
		"status": true,
		"data":   res,
	})
}

func (u *UserController) GetUserByID(c echo.Context) error {
	id, err := strconv.Atoi((c.Param("id")))
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	foundUser, err := u.UserUsecase.ReadByID(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"code":    404,
			"status":  false,
			"message": err.Error(),
		})
	}

	res := response.UserResponse{
		ID:       int(foundUser.ID),
		Email:    foundUser.Email,
		Password: foundUser.Password,
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"code":   200,
		"status": true,
		"data":   res,
	})
}

func (u *UserController) GetUsers(c echo.Context) error {
	foundUsers, err := u.UserUsecase.ReadAll()
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	var res []response.UsersResponse
	for _, foundUser := range *foundUsers {
		res = append(res, response.UsersResponse{
			ID:    int(foundUser.ID),
			Email: foundUser.Email,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"code":   200,
		"status": true,
		"data":   res,
	})
}
