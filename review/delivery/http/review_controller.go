package http

import (
	"github.com/labstack/echo/v4"
	"github.com/nrmadi02/mini-project/domain"
	"github.com/nrmadi02/mini-project/web/response"
	uuid "github.com/satori/go.uuid"
	"net/http"
)

type ReviewValue struct {
	Review string `json:"review"`
}

type ReviewController interface {
	AddReviewEnterprise(c echo.Context) error
	UpdateReviewEnterprise(c echo.Context) error
	DeleteReviewEnterprise(c echo.Context) error
	GetListReviewByEnterpriseID(c echo.Context) error
	GetDetailReviewByID(c echo.Context) error
}

type reviewController struct {
	reviewUsecase     domain.ReviewUsecase
	enterpriseUsecase domain.EnterpriseUsecase
	authUsecase       domain.AuthUsecase
}

func NewReviewController(ru domain.ReviewUsecase, eu domain.EnterpriseUsecase, au domain.AuthUsecase) ReviewController {
	return reviewController{
		reviewUsecase:     ru,
		enterpriseUsecase: eu,
		authUsecase:       au,
	}
}

// AddReviewEnterprise godoc
// @Summary Add Review
// @Description add review enterprise
// @Tags Review
// @accept json
// @Produce json
// @Router /review/enterprise/{id} [post]
// @Param id path string true "enterprise id"
// @Param userid query string true "user id"
// @param data body ReviewValue true "value review"
// @Success 201 {object} response.JSONSuccessResult{data=interface{}}
// @Failure 400 {object} response.JSONBadRequestResult{}
// @Failure 404 {object} response.JSONBadRequestResult{}
// @Security JWT
func (r reviewController) AddReviewEnterprise(c echo.Context) error {
	userid := c.QueryParam("userid")
	enterpriseid := c.Param("id")
	var value ReviewValue
	if err := c.Bind(&value); err != nil {
		return response.FailResponse(c, http.StatusBadRequest, false, err.Error())
	}
	if value.Review == "" {
		return response.FailResponse(c, http.StatusBadRequest, false, "value review empty")
	}
	isReview, _ := r.reviewUsecase.GetReviewByUserIDAndEnterpriseID(enterpriseid, userid)
	if isReview.ID != uuid.FromStringOrNil("") {
		return response.FailResponse(c, http.StatusBadRequest, false, "remove old review")
	}
	review, err := r.reviewUsecase.AddReview(enterpriseid, userid, value.Review)
	if err != nil {
		return response.FailResponse(c, http.StatusBadRequest, false, err.Error())
	}
	return response.SuccessResponse(c, http.StatusCreated, true, "success add review", review)
}

// GetListReviewByEnterpriseID godoc
// @Summary Get List Review
// @Description get list review enterprise
// @Tags Review
// @accept json
// @Produce json
// @Router /review/enterprise/{id} [get]
// @Param id path string true "enterprise id"
// @Success 200 {object} response.JSONSuccessResult{data=interface{}}
// @Failure 400 {object} response.JSONBadRequestResult{}
// @Failure 404 {object} response.JSONBadRequestResult{}
// @Security JWT
func (r reviewController) GetListReviewByEnterpriseID(c echo.Context) error {
	enterpriseid := c.Param("id")
	enterprise, _ := r.enterpriseUsecase.GetDetailEnterpriseByID(enterpriseid)
	if enterprise.ID == uuid.FromStringOrNil("") {
		return response.FailResponse(c, http.StatusBadRequest, false, "enterprise null")
	}
	reviews, err := r.reviewUsecase.GetListReviewsByEnterpriseID(enterpriseid)
	if err != nil {
		return response.FailResponse(c, http.StatusBadRequest, false, err.Error())
	}
	var resReview []interface{}

	for _, review := range reviews {
		user, _, _, _ := r.authUsecase.GetUserDetails(review.UserID.String())
		resReview = append(resReview, struct {
			Review   interface{} `json:"review"`
			FromUser interface{} `json:"from_user"`
		}{
			Review: review,
			FromUser: response.UsersListResponse{
				ID: user.ID, Email: user.Email, Fullname: user.Fullname, Username: user.Username, CreatedAt: user.CreatedAt, UpdatedAt: user.UpdatedAt,
			},
		})
	}

	resFinal := struct {
		Enterprise interface{}   `json:"enterprise"`
		Reviews    []interface{} `json:"reviews"`
	}{enterprise, resReview}

	return response.SuccessResponse(c, http.StatusOK, true, "success get list review", resFinal)
}

// UpdateReviewEnterprise godoc
// @Summary Update Review
// @Description update review enterprise
// @Tags Review
// @accept json
// @Produce json
// @Router /review/enterprise/{id} [put]
// @Param id path string true "enterprise id"
// @Param userid query string true "user id"
// @param data body ReviewValue true "value review"
// @Success 200 {object} response.JSONSuccessResult{data=interface{}}
// @Failure 400 {object} response.JSONBadRequestResult{}
// @Failure 404 {object} response.JSONBadRequestResult{}
// @Security JWT
func (r reviewController) UpdateReviewEnterprise(c echo.Context) error {
	userid := c.QueryParam("userid")
	enterpriseid := c.Param("id")
	var value ReviewValue
	if err := c.Bind(&value); err != nil {
		return response.FailResponse(c, http.StatusBadRequest, false, err.Error())
	}
	if value.Review == "" {
		return response.FailResponse(c, http.StatusBadRequest, false, "value review empty")
	}
	isReview, _ := r.reviewUsecase.GetReviewByUserIDAndEnterpriseID(enterpriseid, userid)
	if isReview.ID == uuid.FromStringOrNil("") {
		return response.FailResponse(c, http.StatusBadRequest, false, "review not current")
	}

	_, err := r.reviewUsecase.UpdateReview(enterpriseid, userid, value.Review)
	if err != nil {
		return response.FailResponse(c, http.StatusBadRequest, false, err.Error())
	}
	review, err := r.reviewUsecase.GetReviewByUserIDAndEnterpriseID(enterpriseid, userid)
	if err != nil {
		return response.FailResponse(c, http.StatusBadRequest, false, err.Error())
	}
	return response.SuccessResponse(c, http.StatusOK, true, "success update review enterprise", review)
}

// DeleteReviewEnterprise godoc
// @Summary Delete Review
// @Description delete review enterprise
// @Tags Review
// @accept json
// @Produce json
// @Router /review/enterprise/{id} [delete]
// @Param id path string true "enterprise id"
// @Param userid query string true "user id"
// @Success 200 {object} response.JSONSuccessDeleteResult{}
// @Failure 400 {object} response.JSONBadRequestResult{}
// @Failure 404 {object} response.JSONBadRequestResult{}
// @Security JWT
func (r reviewController) DeleteReviewEnterprise(c echo.Context) error {
	userid := c.QueryParam("userid")
	enterpriseid := c.Param("id")
	isReview, _ := r.reviewUsecase.GetReviewByUserIDAndEnterpriseID(enterpriseid, userid)
	if isReview.ID == uuid.FromStringOrNil("") {
		return response.FailResponse(c, http.StatusBadRequest, false, "review not current")
	}

	err := r.reviewUsecase.DeleteReview(enterpriseid, userid)
	if err == nil {
		return response.SuccessDeleteResponse(c, http.StatusOK, true, "success delete review")
	}
	return response.FailResponse(c, http.StatusBadRequest, false, err.Error())
}

// GetDetailReviewByID godoc
// @Summary Get Detail Review
// @Description get detail review enterprise
// @Tags Review
// @accept json
// @Produce json
// @Router /review/{id} [get]
// @Param id path string true "id"
// @Success 200 {object} response.JSONSuccessResult{data=interface{}}
// @Failure 400 {object} response.JSONBadRequestResult{}
// @Failure 404 {object} response.JSONBadRequestResult{}
// @Security JWT
func (r reviewController) GetDetailReviewByID(c echo.Context) error {
	id := c.Param("id")
	review, _ := r.reviewUsecase.GetDetailReviewByID(id)
	if review.ID == uuid.FromStringOrNil("") {
		return response.FailResponse(c, http.StatusNotFound, false, "review not found")
	}
	enterprise, _ := r.enterpriseUsecase.GetDetailEnterpriseByID(review.EnterpriseID.String())
	if enterprise.ID == uuid.FromStringOrNil("") {
		return response.FailResponse(c, http.StatusNotFound, false, "enterprise not found")
	}

	resFinal := struct {
		Enterprise interface{} `json:"enterprise"`
		Review     interface{} `json:"review"`
	}{
		enterprise, review,
	}

	return response.SuccessResponse(c, http.StatusOK, true, "success get detail review enterprise", resFinal)
}
