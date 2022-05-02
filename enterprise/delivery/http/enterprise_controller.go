package http

import (
	"encoding/json"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/nrmadi02/mini-project/domain"
	"github.com/nrmadi02/mini-project/web/request"
	"github.com/nrmadi02/mini-project/web/response"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	"net/url"
	"strconv"
)

type EnterpriseController interface {
	CreateNewEnterprise(c echo.Context) error
	UpdateStatusEnterprise(c echo.Context) error
	UpdateEnterpriseByID(c echo.Context) error
	GetEnterpriseByStatus(c echo.Context) error
	GetDetailEnterpriseByID(c echo.Context) error
	GetAllEnterprises(c echo.Context) error
	GetDistance(c echo.Context) error
	DeleteEnterpriseByID(c echo.Context) error

	//rating enterprise
	AddNewRanting(c echo.Context) error
}

type enterpriseController struct {
	authUsecase       domain.AuthUsecase
	enterpriseUsecase domain.EnterpriseUsecase
	ratingUsecase     domain.RatingUsecase
}

func NewEnterpriseController(au domain.AuthUsecase, eu domain.EnterpriseUsecase, ru domain.RatingUsecase) EnterpriseController {
	return enterpriseController{
		authUsecase:       au,
		enterpriseUsecase: eu,
		ratingUsecase:     ru,
	}
}

// CreateNewEnterprise godoc
// @Summary Create new enterprise
// @Description create new enterprise
// @Tags Enterprise
// @accept json
// @Produce json
// @Router /enterprise [post]
// @param data body request.CreateEnterpriseRequest true "required"
// @Success 201 {object} response.JSONSuccessResult{data=interface{}}
// @Failure 400 {object} response.JSONBadRequestResult{}
// @Security JWT
func (e enterpriseController) CreateNewEnterprise(c echo.Context) error {
	var req request.CreateEnterpriseRequest
	if err := c.Bind(&req); err != nil {
		return response.FailResponse(c, http.StatusBadRequest, false, err.Error())
	}

	jwtBearer := c.Get("user").(*jwt.Token)
	claims := jwtBearer.Claims.(jwt.MapClaims)

	id := claims["UserID"].(string)

	res, err := e.enterpriseUsecase.CreateNewEnterprise(req, id)
	if err != nil {
		return response.FailResponse(c, http.StatusBadRequest, false, err.Error())
	}

	enterprise, err := e.enterpriseUsecase.GetDetailEnterpriseByID(res.ID.String())
	if err != nil {
		return response.FailResponse(c, http.StatusBadRequest, false, err.Error())
	}

	return response.SuccessResponse(c, http.StatusCreated, true, "success create new enterprise", enterprise)
}

// UpdateStatusEnterprise godoc
// @Summary Update status enterprise
// @Description 0 = draft, 1 = publish
// @Tags Enterprise
// @accept json
// @Produce json
// @Router /enterprise/{id}/status [put]
// @Param id path string true "id"
// @Param status query int true "status" Enums(0, 1)
// @Success 200 {object} response.JSONSuccessResult{data=interface{}}
// @Failure 400 {object} response.JSONBadRequestResult{}
// @Failure 401 {object} response.JSONUnauthorizedResult{}
// @Security JWT
func (e enterpriseController) UpdateStatusEnterprise(c echo.Context) error {
	id := c.Param("id")
	status, err := strconv.Atoi(c.QueryParam("status"))
	if err != nil {
		return response.FailResponse(c, http.StatusBadRequest, false, err.Error())
	}

	jwtBearer := c.Get("user").(*jwt.Token)
	claims := jwtBearer.Claims.(jwt.MapClaims)

	isAdmin, err := e.authUsecase.CheckIfUserIsAdmin(claims["UserID"].(string))
	if err != nil {
		return response.FailResponse(c, http.StatusBadRequest, false, err.Error())
	}
	if !isAdmin {
		return response.FailResponse(c, http.StatusUnauthorized, false, "only access admin")
	}

	_, err = e.enterpriseUsecase.UpdateStatusEnterprise(id, status)
	if err != nil {
		return response.FailResponse(c, http.StatusBadRequest, false, err.Error())
	}

	enterprise, err := e.enterpriseUsecase.GetDetailEnterpriseByID(id)
	if err != nil {
		return response.FailResponse(c, http.StatusBadRequest, false, err.Error())
	}

	return response.SuccessResponse(c, http.StatusOK, true, "success update status enterprise", enterprise)
}

// GetEnterpriseByStatus godoc
// @Summary Get list enterprise by status
// @Description status draft and publish
// @Tags Enterprise
// @accept json
// @Produce json
// @Router /enterprises/{status} [get]
// @Param status path string true "status"
// @Success 200 {object} response.JSONSuccessResult{data=interface{}}
// @Failure 400 {object} response.JSONBadRequestResult{}
// @Security JWT
func (e enterpriseController) GetEnterpriseByStatus(c echo.Context) error {
	statusParam := c.Param("status")
	var status int
	if statusParam == "draft" {
		status = 0
	} else if statusParam == "publish" {
		status = 1
	} else {
		return response.FailResponse(c, http.StatusBadRequest, false, "setting status not found")
	}

	resEnterprises, err := e.enterpriseUsecase.GetListEnterpriseByStatus(status)
	if err != nil {
		return err
	}

	var res []response.GetListByStatusResponse

	for _, enterprise := range resEnterprises {
		details, _ := e.authUsecase.GetUserDetails(enterprise.UserID.String())
		res = append(res, response.GetListByStatusResponse{
			ID:          enterprise.ID,
			Name:        enterprise.Name,
			NumberPhone: enterprise.NumberPhone,
			UserID:      enterprise.UserID,
			Address:     enterprise.Address,
			Postcode:    enterprise.Postcode,
			Description: enterprise.Description,
			Status:      enterprise.Status,
			Tags:        enterprise.Tags,
			UpdatedAt:   enterprise.UpdatedAt,
			CreatedAt:   enterprise.CreatedAt,
			Latitude:    enterprise.Latitude,
			Longitude:   enterprise.Longitude,
			Owner: response.UserDetailResponse{
				ID: details.ID, Email: details.Email, Fullname: details.Fullname, Username: details.Username, CreatedAt: details.CreatedAt, UpdatedAt: details.UpdatedAt,
			},
		})
	}

	return response.SuccessResponse(c, http.StatusOK, true, "success get list enterprise by status "+statusParam, res)
}

// UpdateEnterpriseByID godoc
// @Summary Update enterprise by id
// @Description Update enterprise
// @Tags Enterprise
// @accept json
// @Produce json
// @Router /enterprise/{id} [put]
// @param id path string true "id"
// @param data body request.CreateEnterpriseRequest true "required"
// @Success 201 {object} response.JSONSuccessResult{data=interface{}}
// @Failure 400 {object} response.JSONBadRequestResult{}
// @Security JWT
func (e enterpriseController) UpdateEnterpriseByID(c echo.Context) error {
	var req request.CreateEnterpriseRequest
	id := c.Param("id")
	if err := c.Bind(&req); err != nil {
		return response.FailResponse(c, http.StatusBadRequest, false, err.Error())
	}

	jwtBearer := c.Get("user").(*jwt.Token)
	claims := jwtBearer.Claims.(jwt.MapClaims)

	userid := claims["UserID"].(string)

	_, err := e.enterpriseUsecase.UpdateEnterpriseByID(id, userid, req)
	if err != nil {
		return response.FailResponse(c, http.StatusBadRequest, false, err.Error())
	}

	enterprise, err := e.enterpriseUsecase.GetDetailEnterpriseByID(id)
	return response.SuccessResponse(c, http.StatusOK, true, "success update enterprise", enterprise)
}

// GetAllEnterprises godoc
// @Summary Get list enterprises
// @Description get all list enterprises
// @Tags Enterprise
// @accept json
// @Produce json
// @Router /enterprises [get]
// @Success 200 {object} response.JSONSuccessResult{data=interface{}}
// @Failure 400 {object} response.JSONBadRequestResult{}
// @Security JWT
func (e enterpriseController) GetAllEnterprises(c echo.Context) error {
	enterprises, err := e.enterpriseUsecase.GetListAllEnterprise()
	if err != nil {
		return response.FailResponse(c, http.StatusBadRequest, false, err.Error())
	}
	return response.SuccessResponse(c, http.StatusOK, true, "success get list enterprises", enterprises)
}

// DeleteEnterpriseByID godoc
// @Summary Delete enterprise by id
// @Description delete enterprise
// @Tags Enterprise
// @accept json
// @Produce json
// @Router /enterprise/{id} [delete]
// @Param id path string true "id"
// @Success 200 {object} response.JSONSuccessDeleteResult{}
// @Failure 400 {object} response.JSONBadRequestResult{}
// @Failure 401 {object} response.JSONUnauthorizedResult{}
// @Security JWT
func (e enterpriseController) DeleteEnterpriseByID(c echo.Context) error {
	id := c.Param("id")
	jwtBearer := c.Get("user").(*jwt.Token)
	claims := jwtBearer.Claims.(jwt.MapClaims)
	userID := claims["UserID"].(string)

	isAdmin, err := e.authUsecase.CheckIfUserIsAdmin(userID)
	if err != nil {
		return response.FailResponse(c, http.StatusBadRequest, false, err.Error())
	}

	enterprise, err := e.enterpriseUsecase.GetDetailEnterpriseByID(id)
	if err != nil {
		return response.FailResponse(c, http.StatusBadRequest, false, err.Error())
	}

	if isAdmin || enterprise.UserID.String() != userID {
		err := e.enterpriseUsecase.DeleteEnterpriseByID(id)
		if err != nil {
			return response.FailResponse(c, http.StatusBadRequest, false, err.Error())
		}

		return response.SuccessDeleteResponse(c, http.StatusOK, true, "success delete enterprise")
	}

	return response.FailResponse(c, http.StatusUnauthorized, false, "not current user or admin")
}

// GetDetailEnterpriseByID godoc
// @Summary Get detail by id
// @Description get detail enterprise
// @Tags Enterprise
// @accept json
// @Produce json
// @Router /enterprise/{id} [get]
// @Param id path string true "id"
// @Success 200 {object} response.JSONSuccessDeleteResult{}
// @Failure 400 {object} response.JSONBadRequestResult{}
// @Security JWT
func (e enterpriseController) GetDetailEnterpriseByID(c echo.Context) error {
	id := c.Param("id")
	enterprise, _ := e.enterpriseUsecase.GetDetailEnterpriseByID(id)
	isNotFound := enterprise.ID.String() != id
	if isNotFound {
		return response.FailResponse(c, http.StatusNotFound, false, "enterprise not found")
	}

	rantings, err := e.ratingUsecase.GetAllRatingByEnterpriseID(enterprise.ID.String())
	if err != nil {
		return response.FailResponse(c, http.StatusBadRequest, false, err.Error())
	}
	var currRat int
	for _, arr := range rantings {
		currRat += arr.Rating
	}
	var rateAvr float64
	rateAvr = float64(currRat) / float64(len(rantings))

	details, _ := e.authUsecase.GetUserDetails(enterprise.UserID.String())
	res := response.GetListByStatusResponse{
		ID:          enterprise.ID,
		Name:        enterprise.Name,
		NumberPhone: enterprise.NumberPhone,
		UserID:      enterprise.UserID,
		Address:     enterprise.Address,
		Postcode:    enterprise.Postcode,
		Description: enterprise.Description,
		Status:      enterprise.Status,
		Tags:        enterprise.Tags,
		UpdatedAt:   enterprise.UpdatedAt,
		CreatedAt:   enterprise.CreatedAt,
		Latitude:    enterprise.Latitude,
		Longitude:   enterprise.Longitude,
		Rating:      math.Round(rateAvr*100) / 100,
		Owner: response.UserDetailResponse{
			ID: details.ID, Email: details.Email, Fullname: details.Fullname, Username: details.Username, CreatedAt: details.CreatedAt, UpdatedAt: details.UpdatedAt,
		},
	}

	return response.SuccessResponse(c, http.StatusOK, true, "success get detail enterprise", res)
}

// GetDistance godoc
// @Summary Get distance
// @Description get distance from you to enterprise
// @Tags Enterprise
// @accept json
// @Produce json
// @Router /enterprise/{id}/distance [get]
// @Param id path string true "id"
// @Param longitude query string true "longitude"
// @Param latitude query string true "latitude"
// @Success 200 {object} response.JSONSuccessResult{}
// @Failure 400 {object} response.JSONBadRequestResult{}
// @Security JWT
func (e enterpriseController) GetDistance(c echo.Context) error {
	longitude := c.QueryParam("longitude")
	latitude := c.QueryParam("latitude")
	id := c.Param("id")

	enterprise, err := e.enterpriseUsecase.GetDetailEnterpriseByID(id)
	if err != nil {
		return response.FailResponse(c, http.StatusBadRequest, false, err.Error())
	}

	urlPath := "https://geo-services-by-mvpc-com.p.rapidapi.com/distance?locationB=" + url.QueryEscape(latitude+","+longitude) + "&locationA=" + url.QueryEscape(enterprise.Latitude+","+enterprise.Longitude) + "&unit=kms"
	req, _ := http.NewRequest("GET", urlPath, nil)
	req.Header.Add("X-RapidAPI-Host", "geo-services-by-mvpc-com.p.rapidapi.com")
	req.Header.Add("X-RapidAPI-Key", "f0e1ed1db2mshab6b204e17c54edp12a883jsne3b08398fedc")
	res, _ := http.DefaultClient.Do(req)

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			//panic(err)
		}
	}(res.Body)
	body, _ := ioutil.ReadAll(res.Body)
	var resBody interface{}
	err = json.Unmarshal(body, &resBody)
	if err != nil {
		return response.FailResponse(c, http.StatusBadRequest, false, err.Error())
	}
	return response.SuccessResponse(c, http.StatusOK, true, "success get distance enterprise", map[string]interface{}{
		"distance":   resBody.(map[string]interface{})["data"],
		"enterprise": enterprise,
	})
}

// AddNewRanting godoc
// @Summary Add rating enterprise
// @Description add rating enterprise rate 1-5
// @Tags Enterprise
// @accept json
// @Produce json
// @Router /enterprise/{id}/rating [post]
// @param id path string true "id"
// @Param value query int true "value rate"
// @Success 201 {object} response.JSONSuccessResult{data=interface{}}
// @Failure 400 {object} response.JSONBadRequestResult{}
// @Security JWT
func (e enterpriseController) AddNewRanting(c echo.Context) error {
	enterpriseid := c.Param("id")
	jwtBearer := c.Get("user").(*jwt.Token)
	claims := jwtBearer.Claims.(jwt.MapClaims)
	userid := claims["UserID"].(string)
	value, _ := strconv.Atoi(c.QueryParam("value"))
	ranting, err := e.ratingUsecase.AddNewRanting(enterpriseid, userid, value)
	if err != nil {
		return response.FailResponse(c, http.StatusBadRequest, false, err.Error())
	}
	rantings, err := e.ratingUsecase.GetAllRatingByEnterpriseID(enterpriseid)
	if err != nil {
		return response.FailResponse(c, http.StatusBadRequest, false, err.Error())
	}
	var currRat int
	for _, arr := range rantings {
		currRat += arr.Rating
	}
	var rateAvr float64
	rateAvr = float64(currRat) / float64(len(rantings))
	return response.SuccessResponse(c, http.StatusCreated, true, "success add rating", map[string]interface{}{
		"rating":         ranting,
		"rating_average": math.Round(rateAvr*100) / 100,
	})
}
