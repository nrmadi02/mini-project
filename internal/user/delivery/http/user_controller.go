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

type UserController interface {
	User(c echo.Context) error
}

type userController struct {
	AuthUsecase   domain.AuthUsecase
	RatingUsecase domain.RatingUsecase
}

func NewUserController(au domain.AuthUsecase, ru domain.RatingUsecase) UserController {
	return userController{
		AuthUsecase:   au,
		RatingUsecase: ru,
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

	user, favorite, enterprises, err := u.AuthUsecase.GetUserDetails(claims["UserID"].(string))
	if err != nil {
		return response.FailResponse(c, http.StatusUnauthorized, false, err.Error())
	}

	var resEnterprises []interface{}
	var resFavorite interface{}

	if favorite.ID == uuid.FromStringOrNil("") {
		resFavorite = nil
	} else {
		resFavorite = favorite
	}

	for _, enterprise := range enterprises {
		rating := u.RatingUsecase.GetAverageRatingEnterprise(enterprise.ID.String())
		finalRating := math.Round(rating*100) / 100
		resEnterprises = append(resEnterprises, struct {
			ID          uuid.UUID   `json:"id"`
			UserID      uuid.UUID   `json:"user_id"`
			Name        string      `json:"name"`
			NumberPhone string      `json:"number_phone"`
			Address     string      `json:"address"`
			Postcode    int         `json:"postcode"`
			Description string      `json:"description"`
			Latitude    string      `json:"latitude"`
			Longitude   string      `json:"longitude"`
			Status      int         `json:"status"`
			Tags        interface{} `json:"tags"`
			CreatedAt   time.Time   `json:"created_at"`
			UpdatedAt   time.Time   `json:"updated_at"`
			Rating      float64     `json:"rating"`
		}{
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
			Rating:      finalRating,
		})
	}

	res := response.UserDetailResponse{
		Fullname:    user.Fullname,
		Username:    user.Username,
		Email:       user.Email,
		ID:          user.ID,
		Enterprises: resEnterprises,
		Favorite:    resFavorite,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}
	return response.SuccessResponse(c, http.StatusOK, true, "success get detail user", res)

}
