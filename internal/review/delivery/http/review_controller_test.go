package http_test

import (
	"encoding/json"
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/nrmadi02/mini-project/domain"
	"github.com/nrmadi02/mini-project/domain/mocks"
	http2 "github.com/nrmadi02/mini-project/internal/review/delivery/http"
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

var dummyReview = domain.Reviews{
	domain.Review{
		ID:           uuid.FromStringOrNil("35d6a9a1-aa5e-41f1-9991-08878dfdf810"),
		Review:       "baguussss",
		EnterpriseID: uuid.FromStringOrNil("35d6a9a1-aa5e-41f1-9991-08878dfdf89a"),
		UserID:       uuid.FromStringOrNil("35d6a9a1-aa5e-41f1-9991-08878dfdf891"),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	},
	domain.Review{
		ID:           uuid.FromStringOrNil("35d6a9a1-aa5e-41f1-9991-08878dfdf811"),
		Review:       "baguussss",
		EnterpriseID: uuid.FromStringOrNil("35d6a9a1-aa5e-41f1-9991-08878dfdf89a"),
		UserID:       uuid.FromStringOrNil("0cf712fc-e631-40c7-8572-54772e698edd"),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
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

func TestReviewController_AddReviewEnterprise(t *testing.T) {
	mockEnterpriseUsecase := new(mocks.EnterpriseUsecase)
	mockReviewUsecase := new(mocks.ReviewUsecase)
	mockAuthUsecase := new(mocks.AuthUsecase)

	reqBody := http2.ReviewValue{
		Review: "bagusss",
	}
	requestReview, _ := json.Marshal(reqBody)
	t.Run("success", func(t *testing.T) {
		e := echo.New()
		req, rec := makeRequestHttp(string(requestReview), echo.POST, "/review/enterprise/"+dummyEnterprise[0].ID.String(), true, true)
		c := e.NewContext(req, rec)
		c.SetPath(base_path + "/review/enterprise/:id")
		c.SetParamNames("id")
		c.SetParamValues(dummyEnterprise[0].ID.String())
		reviewController := http2.NewReviewController(mockReviewUsecase, mockEnterpriseUsecase, mockAuthUsecase)
		mockReviewUsecase.On("GetReviewByUserIDAndEnterpriseID", mock.Anything, mock.Anything).Return(domain.Review{}, nil).Once()
		mockReviewUsecase.On("AddReview", mock.Anything, mock.Anything, mock.Anything).Return(dummyReview[0], nil).Once()
		err := middlewareToken(reviewController.AddReviewEnterprise, c)
		responseBody := parseResponse(rec)
		assert.NoError(t, err)
		assert.Equal(t, float64(201), responseBody["code"])
		mockReviewUsecase.AssertExpectations(t)
	})
	t.Run("error bind echo", func(t *testing.T) {
		e := echo.New()
		req, rec := makeRequestHttp(string(requestReview), echo.POST, "/review/enterprise/"+dummyEnterprise[0].ID.String(), true, false)
		c := e.NewContext(req, rec)
		c.SetPath(base_path + "/review/enterprise/:id")
		c.SetParamNames("id")
		c.SetParamValues(dummyEnterprise[0].ID.String())
		reviewController := http2.NewReviewController(mockReviewUsecase, mockEnterpriseUsecase, mockAuthUsecase)
		err := middlewareToken(reviewController.AddReviewEnterprise, c)
		responseBody := parseResponse(rec)
		assert.NoError(t, err)
		assert.Equal(t, float64(400), responseBody["code"])
		mockReviewUsecase.AssertExpectations(t)
	})
	t.Run("error empty value", func(t *testing.T) {
		reqBody2 := http2.ReviewValue{
			Review: "",
		}
		requestReview2, _ := json.Marshal(reqBody2)
		e := echo.New()
		req, rec := makeRequestHttp(string(requestReview2), echo.POST, "/review/enterprise/"+dummyEnterprise[0].ID.String(), true, true)
		c := e.NewContext(req, rec)
		c.SetPath(base_path + "/review/enterprise/:id")
		c.SetParamNames("id")
		c.SetParamValues(dummyEnterprise[0].ID.String())
		reviewController := http2.NewReviewController(mockReviewUsecase, mockEnterpriseUsecase, mockAuthUsecase)
		err := middlewareToken(reviewController.AddReviewEnterprise, c)
		responseBody := parseResponse(rec)
		assert.NoError(t, err)
		assert.Equal(t, float64(400), responseBody["code"])
		mockReviewUsecase.AssertExpectations(t)
	})
	t.Run("error get review", func(t *testing.T) {
		e := echo.New()
		req, rec := makeRequestHttp(string(requestReview), echo.POST, "/review/enterprise/"+dummyEnterprise[0].ID.String(), true, true)
		c := e.NewContext(req, rec)
		c.SetPath(base_path + "/review/enterprise/:id")
		c.SetParamNames("id")
		c.SetParamValues(dummyEnterprise[0].ID.String())
		reviewController := http2.NewReviewController(mockReviewUsecase, mockEnterpriseUsecase, mockAuthUsecase)
		mockReviewUsecase.On("GetReviewByUserIDAndEnterpriseID", mock.Anything, mock.Anything).Return(dummyReview[0], nil).Once()
		err := middlewareToken(reviewController.AddReviewEnterprise, c)
		responseBody := parseResponse(rec)
		assert.NoError(t, err)
		assert.Equal(t, float64(400), responseBody["code"])
		mockReviewUsecase.AssertExpectations(t)
	})
	t.Run("error add review", func(t *testing.T) {
		e := echo.New()
		req, rec := makeRequestHttp(string(requestReview), echo.POST, "/review/enterprise/"+dummyEnterprise[0].ID.String(), true, true)
		c := e.NewContext(req, rec)
		c.SetPath(base_path + "/review/enterprise/:id")
		c.SetParamNames("id")
		c.SetParamValues(dummyEnterprise[0].ID.String())
		reviewController := http2.NewReviewController(mockReviewUsecase, mockEnterpriseUsecase, mockAuthUsecase)
		mockReviewUsecase.On("GetReviewByUserIDAndEnterpriseID", mock.Anything, mock.Anything).Return(domain.Review{}, nil).Once()
		mockReviewUsecase.On("AddReview", mock.Anything, mock.Anything, mock.Anything).Return(domain.Review{}, errors.New("error something")).Once()
		err := middlewareToken(reviewController.AddReviewEnterprise, c)
		responseBody := parseResponse(rec)
		assert.NoError(t, err)
		assert.Equal(t, float64(400), responseBody["code"])
		mockReviewUsecase.AssertExpectations(t)
	})
}

func TestReviewController_GetListReviewByEnterpriseID(t *testing.T) {
	mockEnterpriseUsecase := new(mocks.EnterpriseUsecase)
	mockReviewUsecase := new(mocks.ReviewUsecase)
	mockAuthUsecase := new(mocks.AuthUsecase)

	t.Run("success", func(t *testing.T) {
		e := echo.New()
		req, rec := makeRequestHttp("", echo.POST, "/review/enterprise/"+dummyEnterprise[0].ID.String(), true, true)
		c := e.NewContext(req, rec)
		c.SetPath(base_path + "/review/enterprise/:id")
		c.SetParamNames("id")
		c.SetParamValues(dummyEnterprise[0].ID.String())
		reviewController := http2.NewReviewController(mockReviewUsecase, mockEnterpriseUsecase, mockAuthUsecase)
		mockEnterpriseUsecase.On("GetDetailEnterpriseByID", mock.Anything).Return(dummyEnterprise[0], nil).Once()
		mockReviewUsecase.On("GetListReviewsByEnterpriseID", mock.Anything).Return(domain.Reviews{dummyReview[0]}, nil).Once()
		mockAuthUsecase.On("GetUserDetails", mock.Anything).Return(dummyUser[0], domain.Favorite{}, domain.Enterprises{}, nil).Once()
		err := middlewareToken(reviewController.GetListReviewByEnterpriseID, c)
		responseBody := parseResponse(rec)
		assert.NoError(t, err)
		assert.Equal(t, float64(200), responseBody["code"])
		mockReviewUsecase.AssertExpectations(t)
	})
	t.Run("error get enterprise", func(t *testing.T) {
		e := echo.New()
		req, rec := makeRequestHttp("", echo.POST, "/review/enterprise/"+dummyEnterprise[0].ID.String(), true, true)
		c := e.NewContext(req, rec)
		c.SetPath(base_path + "/review/enterprise/:id")
		c.SetParamNames("id")
		c.SetParamValues(dummyEnterprise[0].ID.String())
		reviewController := http2.NewReviewController(mockReviewUsecase, mockEnterpriseUsecase, mockAuthUsecase)
		mockEnterpriseUsecase.On("GetDetailEnterpriseByID", mock.Anything).Return(domain.Enterprise{}, errors.New("error something")).Once()
		err := middlewareToken(reviewController.GetListReviewByEnterpriseID, c)
		responseBody := parseResponse(rec)
		assert.NoError(t, err)
		assert.Equal(t, float64(400), responseBody["code"])
		mockReviewUsecase.AssertExpectations(t)
	})
	t.Run("error get list reviews", func(t *testing.T) {
		e := echo.New()
		req, rec := makeRequestHttp("", echo.POST, "/review/enterprise/"+dummyEnterprise[0].ID.String(), true, true)
		c := e.NewContext(req, rec)
		c.SetPath(base_path + "/review/enterprise/:id")
		c.SetParamNames("id")
		c.SetParamValues(dummyEnterprise[0].ID.String())
		reviewController := http2.NewReviewController(mockReviewUsecase, mockEnterpriseUsecase, mockAuthUsecase)
		mockEnterpriseUsecase.On("GetDetailEnterpriseByID", mock.Anything).Return(dummyEnterprise[0], nil).Once()
		mockReviewUsecase.On("GetListReviewsByEnterpriseID", mock.Anything).Return(domain.Reviews{}, errors.New("error something")).Once()
		err := middlewareToken(reviewController.GetListReviewByEnterpriseID, c)
		responseBody := parseResponse(rec)
		assert.NoError(t, err)
		assert.Equal(t, float64(400), responseBody["code"])
		mockReviewUsecase.AssertExpectations(t)
	})
}

func TestReviewController_UpdateReviewEnterprise(t *testing.T) {
	mockEnterpriseUsecase := new(mocks.EnterpriseUsecase)
	mockReviewUsecase := new(mocks.ReviewUsecase)
	mockAuthUsecase := new(mocks.AuthUsecase)

	reqBody := http2.ReviewValue{
		Review: "bagusss",
	}
	requestReview, _ := json.Marshal(reqBody)
	t.Run("success", func(t *testing.T) {
		e := echo.New()
		req, rec := makeRequestHttp(string(requestReview), echo.PUT, "/review/enterprise/"+dummyEnterprise[0].ID.String(), true, true)
		c := e.NewContext(req, rec)
		c.SetPath(base_path + "/review/enterprise/:id")
		c.SetParamNames("id")
		c.SetParamValues(dummyEnterprise[0].ID.String())
		reviewController := http2.NewReviewController(mockReviewUsecase, mockEnterpriseUsecase, mockAuthUsecase)
		mockReviewUsecase.On("GetReviewByUserIDAndEnterpriseID", mock.Anything, mock.Anything).Return(dummyReview[0], nil).Once()
		mockReviewUsecase.On("UpdateReview", mock.Anything, mock.Anything, mock.Anything).Return(dummyReview[0], nil).Once()
		mockReviewUsecase.On("GetReviewByUserIDAndEnterpriseID", mock.Anything, mock.Anything).Return(dummyReview[0], nil).Once()
		err := middlewareToken(reviewController.UpdateReviewEnterprise, c)
		responseBody := parseResponse(rec)
		assert.NoError(t, err)
		assert.Equal(t, float64(200), responseBody["code"])
		mockReviewUsecase.AssertExpectations(t)
	})
	t.Run("error bind echo", func(t *testing.T) {
		e := echo.New()
		req, rec := makeRequestHttp(string(requestReview), echo.POST, "/review/enterprise/"+dummyEnterprise[0].ID.String(), true, false)
		c := e.NewContext(req, rec)
		c.SetPath(base_path + "/review/enterprise/:id")
		c.SetParamNames("id")
		c.SetParamValues(dummyEnterprise[0].ID.String())
		reviewController := http2.NewReviewController(mockReviewUsecase, mockEnterpriseUsecase, mockAuthUsecase)
		err := middlewareToken(reviewController.UpdateReviewEnterprise, c)
		responseBody := parseResponse(rec)
		assert.NoError(t, err)
		assert.Equal(t, float64(400), responseBody["code"])
		mockReviewUsecase.AssertExpectations(t)
	})
	t.Run("error empty value", func(t *testing.T) {
		reqBody2 := http2.ReviewValue{
			Review: "",
		}
		requestReview2, _ := json.Marshal(reqBody2)
		e := echo.New()
		req, rec := makeRequestHttp(string(requestReview2), echo.POST, "/review/enterprise/"+dummyEnterprise[0].ID.String(), true, true)
		c := e.NewContext(req, rec)
		c.SetPath(base_path + "/review/enterprise/:id")
		c.SetParamNames("id")
		c.SetParamValues(dummyEnterprise[0].ID.String())
		reviewController := http2.NewReviewController(mockReviewUsecase, mockEnterpriseUsecase, mockAuthUsecase)
		err := middlewareToken(reviewController.UpdateReviewEnterprise, c)
		responseBody := parseResponse(rec)
		assert.NoError(t, err)
		assert.Equal(t, float64(400), responseBody["code"])
		mockReviewUsecase.AssertExpectations(t)
	})
	t.Run("error get detail review", func(t *testing.T) {
		e := echo.New()
		req, rec := makeRequestHttp(string(requestReview), echo.PUT, "/review/enterprise/"+dummyEnterprise[0].ID.String(), true, true)
		c := e.NewContext(req, rec)
		c.SetPath(base_path + "/review/enterprise/:id")
		c.SetParamNames("id")
		c.SetParamValues(dummyEnterprise[0].ID.String())
		reviewController := http2.NewReviewController(mockReviewUsecase, mockEnterpriseUsecase, mockAuthUsecase)
		mockReviewUsecase.On("GetReviewByUserIDAndEnterpriseID", mock.Anything, mock.Anything).Return(domain.Review{}, errors.New("error something")).Once()
		err := middlewareToken(reviewController.UpdateReviewEnterprise, c)
		responseBody := parseResponse(rec)
		assert.NoError(t, err)
		assert.Equal(t, float64(400), responseBody["code"])
		mockReviewUsecase.AssertExpectations(t)
	})
	t.Run("error/failed update review", func(t *testing.T) {
		e := echo.New()
		req, rec := makeRequestHttp(string(requestReview), echo.PUT, "/review/enterprise/"+dummyEnterprise[0].ID.String(), true, true)
		c := e.NewContext(req, rec)
		c.SetPath(base_path + "/review/enterprise/:id")
		c.SetParamNames("id")
		c.SetParamValues(dummyEnterprise[0].ID.String())
		reviewController := http2.NewReviewController(mockReviewUsecase, mockEnterpriseUsecase, mockAuthUsecase)
		mockReviewUsecase.On("GetReviewByUserIDAndEnterpriseID", mock.Anything, mock.Anything).Return(dummyReview[0], nil).Once()
		mockReviewUsecase.On("UpdateReview", mock.Anything, mock.Anything, mock.Anything).Return(dummyReview[0], errors.New("error something")).Once()
		err := middlewareToken(reviewController.UpdateReviewEnterprise, c)
		responseBody := parseResponse(rec)
		assert.NoError(t, err)
		assert.Equal(t, float64(400), responseBody["code"])
		mockReviewUsecase.AssertExpectations(t)
	})
	t.Run("error get detail review", func(t *testing.T) {
		e := echo.New()
		req, rec := makeRequestHttp(string(requestReview), echo.PUT, "/review/enterprise/"+dummyEnterprise[0].ID.String(), true, true)
		c := e.NewContext(req, rec)
		c.SetPath(base_path + "/review/enterprise/:id")
		c.SetParamNames("id")
		c.SetParamValues(dummyEnterprise[0].ID.String())
		reviewController := http2.NewReviewController(mockReviewUsecase, mockEnterpriseUsecase, mockAuthUsecase)
		mockReviewUsecase.On("GetReviewByUserIDAndEnterpriseID", mock.Anything, mock.Anything).Return(dummyReview[0], nil).Once()
		mockReviewUsecase.On("UpdateReview", mock.Anything, mock.Anything, mock.Anything).Return(dummyReview[0], nil).Once()
		mockReviewUsecase.On("GetReviewByUserIDAndEnterpriseID", mock.Anything, mock.Anything).Return(domain.Review{}, errors.New("error something")).Once()
		err := middlewareToken(reviewController.UpdateReviewEnterprise, c)
		responseBody := parseResponse(rec)
		assert.NoError(t, err)
		assert.Equal(t, float64(400), responseBody["code"])
		mockReviewUsecase.AssertExpectations(t)
	})
}

func TestReviewController_DeleteReviewEnterprise(t *testing.T) {
	mockEnterpriseUsecase := new(mocks.EnterpriseUsecase)
	mockReviewUsecase := new(mocks.ReviewUsecase)
	mockAuthUsecase := new(mocks.AuthUsecase)

	t.Run("success", func(t *testing.T) {
		e := echo.New()
		req, rec := makeRequestHttp("", echo.DELETE, "/review/enterprise/"+dummyReview[0].EnterpriseID.String()+"?userid="+dummyReview[0].UserID.String(), true, true)
		c := e.NewContext(req, rec)
		c.SetPath(base_path + "/review/enterprise/:id")
		c.SetParamNames("id")
		c.SetParamValues(dummyReview[0].EnterpriseID.String())
		reviewController := http2.NewReviewController(mockReviewUsecase, mockEnterpriseUsecase, mockAuthUsecase)
		mockReviewUsecase.On("GetReviewByUserIDAndEnterpriseID", mock.Anything, mock.Anything).Return(dummyReview[0], nil).Once()
		mockReviewUsecase.On("DeleteReview", mock.Anything, mock.Anything).Return(nil).Once()
		err := middlewareToken(reviewController.DeleteReviewEnterprise, c)
		responseBody := parseResponse(rec)
		assert.NoError(t, err)
		assert.Equal(t, float64(200), responseBody["code"])
		mockReviewUsecase.AssertExpectations(t)
	})
	t.Run("failed delete", func(t *testing.T) {
		e := echo.New()
		req, rec := makeRequestHttp("", echo.DELETE, "/review/enterprise/"+dummyReview[0].EnterpriseID.String()+"?userid="+dummyReview[0].UserID.String(), true, true)
		c := e.NewContext(req, rec)
		c.SetPath(base_path + "/review/enterprise/:id")
		c.SetParamNames("id")
		c.SetParamValues(dummyReview[0].EnterpriseID.String())
		reviewController := http2.NewReviewController(mockReviewUsecase, mockEnterpriseUsecase, mockAuthUsecase)
		mockReviewUsecase.On("GetReviewByUserIDAndEnterpriseID", mock.Anything, mock.Anything).Return(dummyReview[0], nil).Once()
		mockReviewUsecase.On("DeleteReview", mock.Anything, mock.Anything).Return(errors.New("error something")).Once()
		err := middlewareToken(reviewController.DeleteReviewEnterprise, c)
		responseBody := parseResponse(rec)
		assert.NoError(t, err)
		assert.Equal(t, float64(400), responseBody["code"])
		mockReviewUsecase.AssertExpectations(t)
	})
	t.Run("error get detail review", func(t *testing.T) {
		e := echo.New()
		req, rec := makeRequestHttp("", echo.DELETE, "/review/enterprise/"+dummyReview[0].EnterpriseID.String()+"?userid="+dummyReview[0].UserID.String(), true, true)
		c := e.NewContext(req, rec)
		c.SetPath(base_path + "/review/enterprise/:id")
		c.SetParamNames("id")
		c.SetParamValues(dummyReview[0].EnterpriseID.String())
		reviewController := http2.NewReviewController(mockReviewUsecase, mockEnterpriseUsecase, mockAuthUsecase)
		mockReviewUsecase.On("GetReviewByUserIDAndEnterpriseID", mock.Anything, mock.Anything).Return(domain.Review{}, errors.New("error something")).Once()
		err := middlewareToken(reviewController.DeleteReviewEnterprise, c)
		responseBody := parseResponse(rec)
		assert.NoError(t, err)
		assert.Equal(t, float64(400), responseBody["code"])
		mockReviewUsecase.AssertExpectations(t)
	})
}

func TestReviewController_GetDetailReviewByID(t *testing.T) {
	mockEnterpriseUsecase := new(mocks.EnterpriseUsecase)
	mockReviewUsecase := new(mocks.ReviewUsecase)
	mockAuthUsecase := new(mocks.AuthUsecase)

	t.Run("success", func(t *testing.T) {
		e := echo.New()
		req, rec := makeRequestHttp("", echo.GET, "/review/"+dummyReview[0].ID.String(), true, true)
		c := e.NewContext(req, rec)
		c.SetPath(base_path + "/review/:id")
		c.SetParamNames("id")
		c.SetParamValues(dummyReview[0].ID.String())
		reviewController := http2.NewReviewController(mockReviewUsecase, mockEnterpriseUsecase, mockAuthUsecase)
		mockReviewUsecase.On("GetDetailReviewByID", mock.Anything).Return(dummyReview[0], nil).Once()
		mockEnterpriseUsecase.On("GetDetailEnterpriseByID", mock.Anything).Return(dummyEnterprise[0], nil).Once()
		err := middlewareToken(reviewController.GetDetailReviewByID, c)
		responseBody := parseResponse(rec)
		assert.NoError(t, err)
		assert.Equal(t, float64(200), responseBody["code"])
		mockReviewUsecase.AssertExpectations(t)
	})
	t.Run("error get detail review", func(t *testing.T) {
		e := echo.New()
		req, rec := makeRequestHttp("", echo.GET, "/review/"+dummyReview[0].ID.String(), true, true)
		c := e.NewContext(req, rec)
		c.SetPath(base_path + "/review/:id")
		c.SetParamNames("id")
		c.SetParamValues(dummyReview[0].ID.String())
		reviewController := http2.NewReviewController(mockReviewUsecase, mockEnterpriseUsecase, mockAuthUsecase)
		mockReviewUsecase.On("GetDetailReviewByID", mock.Anything).Return(domain.Review{}, errors.New("error something")).Once()
		err := middlewareToken(reviewController.GetDetailReviewByID, c)
		responseBody := parseResponse(rec)
		assert.NoError(t, err)
		assert.Equal(t, float64(404), responseBody["code"])
		mockReviewUsecase.AssertExpectations(t)
	})
	t.Run("error get enterprise", func(t *testing.T) {
		e := echo.New()
		req, rec := makeRequestHttp("", echo.GET, "/review/"+dummyReview[0].ID.String(), true, true)
		c := e.NewContext(req, rec)
		c.SetPath(base_path + "/review/:id")
		c.SetParamNames("id")
		c.SetParamValues(dummyReview[0].ID.String())
		reviewController := http2.NewReviewController(mockReviewUsecase, mockEnterpriseUsecase, mockAuthUsecase)
		mockReviewUsecase.On("GetDetailReviewByID", mock.Anything).Return(dummyReview[0], nil).Once()
		mockEnterpriseUsecase.On("GetDetailEnterpriseByID", mock.Anything).Return(domain.Enterprise{}, errors.New("error something")).Once()
		err := middlewareToken(reviewController.GetDetailReviewByID, c)
		responseBody := parseResponse(rec)
		assert.NoError(t, err)
		assert.Equal(t, float64(404), responseBody["code"])
		mockReviewUsecase.AssertExpectations(t)
	})
}
