package http_test

import (
	"encoding/json"
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/nrmadi02/mini-project/domain"
	"github.com/nrmadi02/mini-project/domain/mocks"
	http2 "github.com/nrmadi02/mini-project/internal/favorite/delivery/http"
	"github.com/nrmadi02/mini-project/internal/user/delivery/http/helper"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

var base_path = "/api/v1"

var password, _ = bcrypt.GenerateFromPassword([]byte("12345678"), bcrypt.DefaultCost)
var dummyUser = domain.Users{
	domain.User{
		ID:       uuid.FromStringOrNil("35d6a9a1-aa5e-41f1-9991-08878dfdf891"),
		Fullname: "user1",
		Email:    "satu@email.com",
		Username: "usr1",
		Password: string(password),
		Roles: []domain.Role{
			domain.Role{
				Name: "ROLE_ADMIN", ID: 1,
			},
		},
		Enterprises:      nil,
		RatingEnterprise: nil,
		Reviews:          nil,
		Favorite:         domain.Favorite{},
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	},
	domain.User{
		ID:       uuid.FromStringOrNil("0cf712fc-e631-40c7-8572-54772e698edf"),
		Fullname: "user3",
		Email:    "dua@email.com",
		Username: "usr2",
		Password: string(password),
		Roles: []domain.Role{
			domain.Role{
				Name: "ROLE_CLIENT", ID: 1,
			},
		},
		Enterprises:      nil,
		RatingEnterprise: nil,
		Reviews:          nil,
		Favorite:         domain.Favorite{},
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	},
}

var dummyEnterprise = domain.Enterprises{
	domain.Enterprise{
		ID:          uuid.FromStringOrNil("35d6a9a1-aa5e-41f1-9991-08878dfdf89a"),
		UserID:      uuid.FromStringOrNil("35d6a9a1-aa5e-41f1-9991-08878dfdf891"),
		Name:        "enterprise satu",
		NumberPhone: "0012798232",
		Address:     "bjb",
		Postcode:    707722,
		Latitude:    "4235,23",
		Longitude:   "4225,233231",
		Description: "testing1",
		Status:      0,
		Tags: domain.Tags{
			domain.Tag{
				ID:   uuid.FromStringOrNil("35d6a9a1-aa5e-41f1-9991-08878dfdf89b"),
				Name: "Tag satu",
			},
		},
		RatingEnterprise: nil,
		Reviews:          nil,
		CreatedAt:        time.Time{},
		UpdatedAt:        time.Time{},
	},
	domain.Enterprise{
		ID:               uuid.FromStringOrNil("35d6a9a1-aa5e-41f1-9991-08878dfdf89b"),
		UserID:           uuid.FromStringOrNil("35d6a9a1-aa5e-41f1-9991-08878dfdf892"),
		Name:             "enterprise dua",
		NumberPhone:      "0012798232",
		Address:          "bjm",
		Postcode:         707722,
		Latitude:         "4235,23",
		Longitude:        "4225,233231",
		Description:      "testing1",
		Status:           0,
		Tags:             nil,
		RatingEnterprise: nil,
		Reviews:          nil,
		CreatedAt:        time.Time{},
		UpdatedAt:        time.Time{},
	},
}

var dummyFavorite = domain.Favorites{
	domain.Favorite{
		ID:     uuid.FromStringOrNil("0cf712fc-e631-40c7-8572-54772e698e10"),
		UserID: uuid.FromStringOrNil("0cf712fc-e631-40c7-8572-54772e698edf"),
		Enterprises: []domain.Enterprise{
			dummyEnterprise[0],
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	},
	domain.Favorite{
		ID:     uuid.FromStringOrNil("0cf712fc-e631-40c7-8572-54772e698e11"),
		UserID: uuid.FromStringOrNil("0cf712fc-e631-40c7-8572-54772e698edf"),
		Enterprises: []domain.Enterprise{
			dummyEnterprise[0],
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	},
}

var dummyRating = domain.RatingEnterprises{
	domain.RatingEnterprise{
		ID:           uuid.FromStringOrNil("35d6a9a1-aa5e-41f1-9991-08878dfdf777"),
		Rating:       3,
		EnterpriseID: uuid.FromStringOrNil("35d6a9a1-aa5e-41f1-9991-08878dfdf89a"),
		UserID:       uuid.FromStringOrNil("0cf712fc-e631-40c7-8572-54772e698edf"),
	},
	domain.RatingEnterprise{
		ID:           uuid.FromStringOrNil("35d6a9a1-aa5e-41f1-9991-08878dfdf878"),
		Rating:       0,
		EnterpriseID: uuid.FromStringOrNil("35d6a9a1-aa5e-41f1-9991-08878dfdf89a"),
		UserID:       uuid.FromStringOrNil("0cf712fc-e631-40c7-8572-54772e698787"),
	},
}

func createToken() string {
	jwtSetToken := helper.NewGoJWT()
	token := jwtSetToken.CreateTokenJWT(&dummyUser[0])
	return token
}

func middlewareToken(handlerFunc echo.HandlerFunc, c echo.Context) error {
	err := middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte("220220"),
	})(handlerFunc)(c)
	return err
}

func parseResponse(rec *httptest.ResponseRecorder) map[string]interface{} {
	var responseBody map[string]interface{}
	resBody := rec.Body.String()
	_ = json.Unmarshal([]byte(resBody), &responseBody)
	return responseBody
}

func makeRequestHttp(request string, method string, path string, isToken bool, isBind bool) (req *http.Request, rec *httptest.ResponseRecorder) {
	req, _ = http.NewRequest(method, base_path+path, strings.NewReader(request))
	if isBind {
		req.Header.Add("Content-Type", "application/json")
	}
	if isToken {
		req.Header.Add(echo.HeaderAuthorization, middleware.DefaultJWTConfig.AuthScheme+" "+createToken())
	}
	rec = httptest.NewRecorder()
	return req, rec
}

func TestFavoriteController_AddFavoriteEnterprise(t *testing.T) {
	favoriteUsecase := new(mocks.FavoriteUsecase)
	authUsecase := new(mocks.AuthUsecase)
	ratingUsecase := new(mocks.RatingUsecase)
	requestBody := []string{dummyEnterprise[0].ID.String()}
	requestFavorite, _ := json.Marshal(requestBody)
	t.Run("success", func(t *testing.T) {
		e := echo.New()
		req, rec := makeRequestHttp(string(requestFavorite), echo.POST, "/favorite", true, true)
		c := e.NewContext(req, rec)
		favoriteUsecase.On("AddFavorite", mock.Anything, mock.Anything).Return(domain.Favorite{}, nil).Once()
		favoriteUsecase.On("GetDetailByUserID", mock.Anything).Return(dummyFavorite[0], nil).Once()
		favoriteController := http2.NewFavoriteController(favoriteUsecase, authUsecase, ratingUsecase)
		err := middlewareToken(favoriteController.AddFavoriteEnterprise, c)
		responseBody := parseResponse(rec)
		assert.NoError(t, err)
		assert.Equal(t, 200, int(responseBody["code"].(float64)))
		assert.NotNil(t, responseBody["data"])
		favoriteUsecase.AssertExpectations(t)
	})
	t.Run("error echo bind", func(t *testing.T) {
		e := echo.New()
		req, rec := makeRequestHttp(string(requestFavorite), echo.POST, "/favorite", true, false)
		c := e.NewContext(req, rec)
		favoriteController := http2.NewFavoriteController(favoriteUsecase, authUsecase, ratingUsecase)
		err := middlewareToken(favoriteController.AddFavoriteEnterprise, c)
		responseBody := parseResponse(rec)
		assert.NoError(t, err)
		assert.Equal(t, 400, int(responseBody["code"].(float64)))
		favoriteUsecase.AssertExpectations(t)
	})
	t.Run("error add favorite", func(t *testing.T) {
		e := echo.New()
		req, rec := makeRequestHttp(string(requestFavorite), echo.POST, "/favorite", true, true)
		c := e.NewContext(req, rec)
		favoriteUsecase.On("AddFavorite", mock.Anything, mock.Anything).Return(domain.Favorite{}, errors.New("error something")).Once()
		favoriteController := http2.NewFavoriteController(favoriteUsecase, authUsecase, ratingUsecase)
		err := middlewareToken(favoriteController.AddFavoriteEnterprise, c)
		responseBody := parseResponse(rec)
		assert.NoError(t, err)
		assert.Equal(t, 400, int(responseBody["code"].(float64)))
		favoriteUsecase.AssertExpectations(t)
	})
	t.Run("error get detail", func(t *testing.T) {
		e := echo.New()
		req, rec := makeRequestHttp(string(requestFavorite), echo.POST, "/favorite", true, true)
		c := e.NewContext(req, rec)
		favoriteUsecase.On("AddFavorite", mock.Anything, mock.Anything).Return(domain.Favorite{}, nil).Once()
		favoriteUsecase.On("GetDetailByUserID", mock.Anything).Return(domain.Favorite{}, errors.New("error someting")).Once()
		favoriteController := http2.NewFavoriteController(favoriteUsecase, authUsecase, ratingUsecase)
		err := middlewareToken(favoriteController.AddFavoriteEnterprise, c)
		responseBody := parseResponse(rec)
		assert.NoError(t, err)
		assert.Equal(t, 404, int(responseBody["code"].(float64)))
		favoriteUsecase.AssertExpectations(t)
	})
}

func TestFavoriteController_RemoveFavoriteEnterprise(t *testing.T) {
	favoriteUsecase := new(mocks.FavoriteUsecase)
	authUsecase := new(mocks.AuthUsecase)
	ratingUsecase := new(mocks.RatingUsecase)
	requestBody := []string{dummyEnterprise[0].ID.String()}
	requestFavorite, _ := json.Marshal(requestBody)
	t.Run("success", func(t *testing.T) {
		e := echo.New()
		req, rec := makeRequestHttp(string(requestFavorite), echo.DELETE, "/favorite", true, true)
		c := e.NewContext(req, rec)
		favoriteUsecase.On("RemoveFavorite", mock.Anything, mock.Anything).Return(domain.Favorite{}, nil).Once()
		favoriteUsecase.On("GetDetailByUserID", mock.Anything).Return(dummyFavorite[0], nil).Once()
		favoriteController := http2.NewFavoriteController(favoriteUsecase, authUsecase, ratingUsecase)
		err := middlewareToken(favoriteController.RemoveFavoriteEnterprise, c)
		responseBody := parseResponse(rec)
		assert.NoError(t, err)
		assert.Equal(t, 200, int(responseBody["code"].(float64)))
		favoriteUsecase.AssertExpectations(t)
	})
	t.Run("error echo bind", func(t *testing.T) {
		e := echo.New()
		req, rec := makeRequestHttp(string(requestFavorite), echo.DELETE, "/favorite", true, false)
		c := e.NewContext(req, rec)
		favoriteController := http2.NewFavoriteController(favoriteUsecase, authUsecase, ratingUsecase)
		err := middlewareToken(favoriteController.RemoveFavoriteEnterprise, c)
		responseBody := parseResponse(rec)
		assert.NoError(t, err)
		assert.Equal(t, 400, int(responseBody["code"].(float64)))
		favoriteUsecase.AssertExpectations(t)
	})
	t.Run("error remove favorite", func(t *testing.T) {
		e := echo.New()
		req, rec := makeRequestHttp(string(requestFavorite), echo.DELETE, "/favorite", true, true)
		c := e.NewContext(req, rec)
		favoriteUsecase.On("RemoveFavorite", mock.Anything, mock.Anything).Return(domain.Favorite{}, errors.New("error something")).Once()
		favoriteController := http2.NewFavoriteController(favoriteUsecase, authUsecase, ratingUsecase)
		err := middlewareToken(favoriteController.RemoveFavoriteEnterprise, c)
		responseBody := parseResponse(rec)
		assert.NoError(t, err)
		assert.Equal(t, 400, int(responseBody["code"].(float64)))
		favoriteUsecase.AssertExpectations(t)
	})
	t.Run("error get detail", func(t *testing.T) {
		e := echo.New()
		req, rec := makeRequestHttp(string(requestFavorite), echo.DELETE, "/favorite", true, true)
		c := e.NewContext(req, rec)
		favoriteUsecase.On("RemoveFavorite", mock.Anything, mock.Anything).Return(domain.Favorite{}, nil).Once()
		favoriteUsecase.On("GetDetailByUserID", mock.Anything).Return(domain.Favorite{}, errors.New("error someting")).Once()
		favoriteController := http2.NewFavoriteController(favoriteUsecase, authUsecase, ratingUsecase)
		err := middlewareToken(favoriteController.RemoveFavoriteEnterprise, c)
		responseBody := parseResponse(rec)
		assert.NoError(t, err)
		assert.Equal(t, 404, int(responseBody["code"].(float64)))
		favoriteUsecase.AssertExpectations(t)
	})
}

func TestFavoriteController_GetDetailFavoriteEnterprise(t *testing.T) {
	favoriteUsecase := new(mocks.FavoriteUsecase)
	authUsecase := new(mocks.AuthUsecase)
	ratingUsecase := new(mocks.RatingUsecase)
	t.Run("success", func(t *testing.T) {
		e := echo.New()
		req, rec := makeRequestHttp("", echo.GET, "/favorite", true, true)
		c := e.NewContext(req, rec)
		favoriteUsecase.On("GetDetailByUserID", mock.Anything).Return(dummyFavorite[0], nil).Once()
		ratingUsecase.On("GetAverageRatingEnterprise", mock.Anything).Return(float64(3), nil).Once()
		favoriteController := http2.NewFavoriteController(favoriteUsecase, authUsecase, ratingUsecase)
		err := middlewareToken(favoriteController.GetDetailFavoriteEnterprise, c)
		responseBody := parseResponse(rec)
		assert.NoError(t, err)
		assert.Equal(t, 200, int(responseBody["code"].(float64)))
		favoriteUsecase.AssertExpectations(t)
	})
	t.Run("error get detail", func(t *testing.T) {
		e := echo.New()
		req, rec := makeRequestHttp("", echo.GET, "/favorite", true, true)
		c := e.NewContext(req, rec)
		favoriteUsecase.On("GetDetailByUserID", mock.Anything).Return(domain.Favorite{}, errors.New("error something")).Once()
		favoriteController := http2.NewFavoriteController(favoriteUsecase, authUsecase, ratingUsecase)
		err := middlewareToken(favoriteController.GetDetailFavoriteEnterprise, c)
		responseBody := parseResponse(rec)
		assert.NoError(t, err)
		assert.Equal(t, 404, int(responseBody["code"].(float64)))
		favoriteUsecase.AssertExpectations(t)
	})
}
