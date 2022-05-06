package http

import (
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/nrmadi02/mini-project/domain"
	"github.com/nrmadi02/mini-project/web/request"
	"github.com/nrmadi02/mini-project/web/response"
	"net/http"
)

type TagController interface {
	GetTagsList(c echo.Context) error
	DeleteTag(c echo.Context) error
	CreateTag(c echo.Context) error
}

type tagController struct {
	authUsecase domain.AuthUsecase
	tagUsecase  domain.TagUsecase
}

func NewTagController(au domain.AuthUsecase, tu domain.TagUsecase) TagController {
	return tagController{
		authUsecase: au,
		tagUsecase:  tu,
	}
}

// GetTagsList godoc
// @Summary Get list tags
// @Description Get list tags
// @Tags Tag
// @accept json
// @Produce json
// @Router /tags [get]
// @Success 200 {object} response.JSONSuccessResult{data=[]response.TagsListResponse}
// @Failure 400 {object} response.JSONBadRequestResult{}
// @Security JWT
func (t tagController) GetTagsList(c echo.Context) error {

	tags, err := t.tagUsecase.GetAllTags()
	if err != nil {
		return response.FailResponse(c, http.StatusBadRequest, false, err.Error())
	}

	var res []response.TagsListResponse
	for _, tag := range tags {
		res = append(res, response.TagsListResponse{
			ID:   tag.ID,
			Name: tag.Name,
		})
	}

	return response.SuccessResponse(c, http.StatusOK, true, "success get list tags", res)
}

// DeleteTag godoc
// @Summary Delete tag by id
// @Description delete tag can access only admin
// @Tags Tag
// @accept json
// @Produce json
// @Router /tag/{id} [delete]
// @Param id path string true "tag id"
// @Success 200 {object} response.JSONSuccessDeleteResult{}
// @Failure 400 {object} response.JSONBadRequestResult{}
// @Failure 401 {object} response.JSONUnauthorizedResult{}
// @Security JWT
func (t tagController) DeleteTag(c echo.Context) error {
	id := c.Param("id")

	jwtBearer := c.Get("user").(*jwt.Token)
	claims := jwtBearer.Claims.(jwt.MapClaims)

	isAdmin, err := t.authUsecase.CheckIfUserIsAdmin(claims["UserID"].(string))
	if err != nil {
		return response.FailResponse(c, http.StatusBadRequest, false, err.Error())
	}
	if !isAdmin {
		return response.FailResponse(c, http.StatusUnauthorized, false, "only access admin")
	}

	err = t.tagUsecase.DeleteTag(id)
	if err != nil {
		return response.FailResponse(c, http.StatusBadRequest, false, err.Error())
	}

	return response.SuccessDeleteResponse(c, http.StatusOK, true, "success delete tag")
}

// CreateTag godoc
// @Summary Create tag
// @Description create tag can access only admin
// @Tags Tag
// @accept json
// @Produce json
// @param data body request.CreateTagRequest true "required"
// @Router /tag [post]
// @Success 200 {object} response.JSONSuccessResult{data=domain.Tag}
// @Failure 400 {object} response.JSONBadRequestResult{}
// @Failure 401 {object} response.JSONUnauthorizedResult{}
// @Security JWT
func (t tagController) CreateTag(c echo.Context) error {
	var req request.CreateTagRequest
	if err := c.Bind(&req); err != nil {
		return response.FailResponse(c, http.StatusBadRequest, false, err.Error())
	}

	if val, err := request.ValidateCreationTag(req); val == false {
		return response.FailResponse(c, http.StatusBadRequest, false, err.Error())
	}

	jwtBearer := c.Get("user").(*jwt.Token)
	claims := jwtBearer.Claims.(jwt.MapClaims)

	isAdmin, err := t.authUsecase.CheckIfUserIsAdmin(claims["UserID"].(string))
	if err != nil {
		return response.FailResponse(c, http.StatusBadRequest, false, err.Error())
	}
	if !isAdmin {
		return response.FailResponse(c, http.StatusUnauthorized, false, "only access admin")
	}

	res, err := t.tagUsecase.CreateNewTag(req)
	if err != nil {
		return response.FailResponse(c, http.StatusBadRequest, false, err.Error())
	}

	return response.SuccessResponse(c, http.StatusCreated, true, "success create new tag", res)
}
