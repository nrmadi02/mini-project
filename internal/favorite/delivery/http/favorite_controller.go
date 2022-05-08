package http

import (
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/nrmadi02/mini-project/domain"
	"github.com/nrmadi02/mini-project/web/response"
	uuid "github.com/satori/go.uuid"
	"math"
	"net/http"
	"time"
)

type FavoriteController interface {
	AddFavoriteEnterprise(c echo.Context) error
	RemoveFavoriteEnterprise(c echo.Context) error
	GetDetailFavoriteEnterprise(c echo.Context) error
}

type favoriteController struct {
	favoriteUsecase domain.FavoriteUsecase
	authUsecase     domain.AuthUsecase
	ratingUsecase   domain.RatingUsecase
}

func NewFavoriteController(fu domain.FavoriteUsecase, au domain.AuthUsecase, ru domain.RatingUsecase) FavoriteController {
	return favoriteController{
		favoriteUsecase: fu,
		authUsecase:     au,
		ratingUsecase:   ru,
	}
}

// AddFavoriteEnterprise godoc
// @Summary Add favorite
// @Description add favorite enterprise
// @Tags favorite
// @accept json
// @Produce json
// @Router /favorite [post]
// @param data body []string true "enterprise id"
// @Success 201 {object} response.JSONSuccessResult{data=interface{}}
// @Failure 400 {object} response.JSONBadRequestResult{}
// @Failure 404 {object} response.JSONBadRequestResult{}
// @Security JWT
func (f favoriteController) AddFavoriteEnterprise(c echo.Context) error {
	jwtBearer := c.Get("user").(*jwt.Token)
	claims := jwtBearer.Claims.(jwt.MapClaims)
	userid := claims["UserID"].(string)
	var req []string
	if err := c.Bind(&req); err != nil {
		return response.FailResponse(c, http.StatusBadRequest, false, err.Error())
	}
	_, err := f.favoriteUsecase.AddFavorite(req, userid)
	if err != nil {
		return response.FailResponse(c, http.StatusBadRequest, false, err.Error())
	}
	favorite, err := f.favoriteUsecase.GetDetailByUserID(userid)
	if err != nil {
		return response.FailResponse(c, http.StatusNotFound, false, err.Error())
	}

	return response.SuccessResponse(c, http.StatusOK, true, "success add favorite", favorite)
}

// RemoveFavoriteEnterprise godoc
// @Summary Remove favorite
// @Description remove favorite enterprise
// @Tags favorite
// @accept json
// @Produce json
// @Router /favorite [delete]
// @param data body []string true "enterprise id"
// @Success 201 {object} response.JSONSuccessResult{data=interface{}}
// @Failure 400 {object} response.JSONBadRequestResult{}
// @Failure 404 {object} response.JSONBadRequestResult{}
// @Security JWT
func (f favoriteController) RemoveFavoriteEnterprise(c echo.Context) error {
	jwtBearer := c.Get("user").(*jwt.Token)
	claims := jwtBearer.Claims.(jwt.MapClaims)
	userid := claims["UserID"].(string)
	var req []string
	if err := c.Bind(&req); err != nil {
		return response.FailResponse(c, http.StatusBadRequest, false, err.Error())
	}
	_, err := f.favoriteUsecase.RemoveFavorite(req, userid)
	if err != nil {
		return response.FailResponse(c, http.StatusBadRequest, false, err.Error())
	}
	favorite, err := f.favoriteUsecase.GetDetailByUserID(userid)
	if err != nil {
		return response.FailResponse(c, http.StatusNotFound, false, err.Error())
	}

	return response.SuccessResponse(c, http.StatusOK, true, "success remove favorite", favorite)
}

// GetDetailFavoriteEnterprise godoc
// @Summary Get favorite
// @Description get favorite enterprise
// @Tags favorite
// @accept json
// @Produce json
// @Router /favorite [get]
// @Success 201 {object} response.JSONSuccessResult{data=interface{}}
// @Failure 400 {object} response.JSONBadRequestResult{}
// @Failure 404 {object} response.JSONBadRequestResult{}
// @Security JWT
func (f favoriteController) GetDetailFavoriteEnterprise(c echo.Context) error {
	jwtBearer := c.Get("user").(*jwt.Token)
	claims := jwtBearer.Claims.(jwt.MapClaims)
	userid := claims["UserID"].(string)

	favorite, err := f.favoriteUsecase.GetDetailByUserID(userid)
	if err != nil {
		return response.FailResponse(c, http.StatusNotFound, false, err.Error())
	}

	var res []response.GetListByStatusResponse

	for _, enterprise := range favorite.Enterprises {
		details, _, _, _ := f.authUsecase.GetUserDetails(enterprise.UserID.String())
		rantings, err := f.ratingUsecase.GetAllRatingByEnterpriseID(enterprise.ID.String())
		if err != nil {
			return response.FailResponse(c, http.StatusBadRequest, false, err.Error())
		}
		var currRat int
		var finalRating float64
		if len(rantings) != 0 {
			for _, arr := range rantings {
				currRat += arr.Rating
			}
			var rateAvr float64
			rateAvr = float64(currRat) / float64(len(rantings))
			finalRating = math.Round(rateAvr*100) / 100
		} else {
			finalRating = 0
		}
		res = append(res, response.GetListByStatusResponse{
			ID:          enterprise.ID,
			Name:        enterprise.Name,
			NumberPhone: enterprise.NumberPhone,
			UserID:      enterprise.UserID,
			Address:     enterprise.Address,
			Postcode:    enterprise.Postcode,
			Description: enterprise.Description,
			Status:      enterprise.Status,
			UpdatedAt:   enterprise.UpdatedAt,
			CreatedAt:   enterprise.CreatedAt,
			Latitude:    enterprise.Latitude,
			Longitude:   enterprise.Longitude,
			Rating:      finalRating,
			Owner: response.UserDetailResponse{
				ID: details.ID, Email: details.Email, Fullname: details.Fullname, Username: details.Username, CreatedAt: details.CreatedAt, UpdatedAt: details.UpdatedAt,
			},
		})
	}
	resFinal := struct {
		ID          uuid.UUID                          `json:"id"`
		UserID      uuid.UUID                          `json:"user_id"`
		Enterprises []response.GetListByStatusResponse `json:"enterprises"`
		CreatedAt   time.Time                          `json:"created_at"`
		UpdatedAt   time.Time                          `json:"updated_at"`
	}{
		ID:          favorite.ID,
		UserID:      favorite.UserID,
		UpdatedAt:   favorite.UpdatedAt,
		CreatedAt:   favorite.CreatedAt,
		Enterprises: res,
	}

	return response.SuccessResponse(c, http.StatusOK, true, "success get detail favorite", resFinal)
}
