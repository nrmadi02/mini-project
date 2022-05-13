package http_test

import (
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/nrmadi02/mini-project/domain"
	"github.com/nrmadi02/mini-project/domain/mocks"
	"github.com/nrmadi02/mini-project/internal/user/delivery/http"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

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

func TestUserController_User(t *testing.T) {
	mockAuthusecase := new(mocks.AuthUsecase)
	mockRatingUsecase := new(mocks.RatingUsecase)

	t.Run("success", func(t *testing.T) {
		e := echo.New()
		req, rec := makeRequestHttp("", echo.GET, "/user", true, true)
		c := e.NewContext(req, rec)
		userController := http.NewUserController(mockAuthusecase, mockRatingUsecase)
		mockAuthusecase.On("GetUserDetails", mock.Anything).Return(dummyUser[0], dummyFavorite[1], domain.Enterprises{dummyEnterprise[0]}, nil).Once()
		mockRatingUsecase.On("GetAllRatingByEnterpriseID", mock.Anything).Return(domain.RatingEnterprises{dummyRating[0]}, nil).Once()
		err := middlewareToken(userController.User, c)
		responseBody := parseResponse(rec)
		assert.NoError(t, err)
		assert.Equal(t, float64(200), responseBody["code"])
		mockAuthusecase.AssertExpectations(t)
		mockRatingUsecase.AssertExpectations(t)
	})
	t.Run("error get detail user", func(t *testing.T) {
		e := echo.New()
		req, rec := makeRequestHttp("", echo.GET, "/user", true, true)
		c := e.NewContext(req, rec)
		userController := http.NewUserController(mockAuthusecase, mockRatingUsecase)
		mockAuthusecase.On("GetUserDetails", mock.Anything).Return(domain.User{}, domain.Favorite{}, domain.Enterprises{}, errors.New("error something")).Once()
		err := middlewareToken(userController.User, c)
		responseBody := parseResponse(rec)
		assert.NoError(t, err)
		assert.Equal(t, float64(401), responseBody["code"])
		mockAuthusecase.AssertExpectations(t)
		mockRatingUsecase.AssertExpectations(t)
	})
	t.Run("get without favorite", func(t *testing.T) {
		e := echo.New()
		req, rec := makeRequestHttp("", echo.GET, "/user", true, true)
		c := e.NewContext(req, rec)
		userController := http.NewUserController(mockAuthusecase, mockRatingUsecase)
		mockAuthusecase.On("GetUserDetails", mock.Anything).Return(dummyUser[0], domain.Favorite{}, domain.Enterprises{dummyEnterprise[0]}, nil).Once()
		mockRatingUsecase.On("GetAllRatingByEnterpriseID", mock.Anything).Return(domain.RatingEnterprises{dummyRating[0]}, nil).Once()
		err := middlewareToken(userController.User, c)
		responseBody := parseResponse(rec)
		assert.NoError(t, err)
		assert.Equal(t, float64(200), responseBody["code"])
		mockAuthusecase.AssertExpectations(t)
		mockRatingUsecase.AssertExpectations(t)
	})
}
