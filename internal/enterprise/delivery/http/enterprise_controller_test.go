package http_test

import (
	"encoding/json"
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/nrmadi02/mini-project/domain"
	"github.com/nrmadi02/mini-project/domain/mocks"
	http2 "github.com/nrmadi02/mini-project/internal/enterprise/delivery/http"
	"github.com/nrmadi02/mini-project/internal/user/delivery/http/helper"
	"github.com/nrmadi02/mini-project/web/request"
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

func TestEnterpriseController_CreateNewEnterprise(t *testing.T) {
	mockEnterpriseUsecase := new(mocks.EnterpriseUsecase)
	mockRatingUsecase := new(mocks.RatingUsecase)
	mockAuthUsecase := new(mocks.AuthUsecase)
	reqCreate := request.CreateEnterpriseRequest{
		Name:        "enterprise satu",
		NumberPhone: "0012798232",
		Address:     "bjb",
		Postcode:    707722,
		Latitude:    "4235,23",
		Longitude:   "4225,233231",
		Description: "testing1",
		Tags:        []string{uuid.FromStringOrNil("35d6a9a1-aa5e-41f1-9991-08878dfdf89b").String()},
	}
	requestCreate, _ := json.Marshal(reqCreate)
	t.Run("success", func(t *testing.T) {
		e := echo.New()
		req, rec := makeRequestHttp(string(requestCreate), echo.POST, "/enterprise", true, true)
		c := e.NewContext(req, rec)
		mockEnterpriseUsecase.On("CreateNewEnterprise", mock.Anything, mock.Anything).Return(dummyEnterprise[0], nil).Once()
		mockEnterpriseUsecase.On("GetDetailEnterpriseByID", mock.Anything).Return(dummyEnterprise[0], nil).Once()
		enterpriseController := http2.NewEnterpriseController(mockAuthUsecase, mockEnterpriseUsecase, mockRatingUsecase)
		err := middlewareToken(enterpriseController.CreateNewEnterprise, c)
		responseBody := parseResponse(rec)
		assert.NoError(t, err)
		assert.Equal(t, 201, int(responseBody["code"].(float64)))
		assert.NotNil(t, responseBody["data"])
		mockEnterpriseUsecase.AssertExpectations(t)
	})
	t.Run("failed bind", func(t *testing.T) {
		e := echo.New()
		req, rec := makeRequestHttp(string(requestCreate), echo.POST, "/enterprise", true, false)
		c := e.NewContext(req, rec)
		enterpriseController := http2.NewEnterpriseController(mockAuthUsecase, mockEnterpriseUsecase, mockRatingUsecase)
		err := middlewareToken(enterpriseController.CreateNewEnterprise, c)
		assert.NoError(t, err)
		mockEnterpriseUsecase.AssertExpectations(t)
	})
	t.Run("failed create", func(t *testing.T) {
		e := echo.New()
		req, rec := makeRequestHttp(string(requestCreate), echo.POST, "/enterprise", true, true)
		c := e.NewContext(req, rec)
		enterpriseController := http2.NewEnterpriseController(mockAuthUsecase, mockEnterpriseUsecase, mockRatingUsecase)
		mockEnterpriseUsecase.On("CreateNewEnterprise", mock.Anything, mock.Anything).Return(domain.Enterprise{}, errors.New("error something")).Once()
		err := middlewareToken(enterpriseController.CreateNewEnterprise, c)
		assert.NoError(t, err)
		mockEnterpriseUsecase.AssertExpectations(t)
	})
	t.Run("failed get by id", func(t *testing.T) {
		e := echo.New()
		req, rec := makeRequestHttp(string(requestCreate), echo.POST, "/enterprise", true, true)
		c := e.NewContext(req, rec)
		enterpriseController := http2.NewEnterpriseController(mockAuthUsecase, mockEnterpriseUsecase, mockRatingUsecase)
		mockEnterpriseUsecase.On("CreateNewEnterprise", mock.Anything, mock.Anything).Return(dummyEnterprise[0], nil).Once()
		mockEnterpriseUsecase.On("GetDetailEnterpriseByID", mock.Anything).Return(domain.Enterprise{}, errors.New("error something")).Once()
		err := middlewareToken(enterpriseController.CreateNewEnterprise, c)
		assert.NoError(t, err)
		mockEnterpriseUsecase.AssertExpectations(t)
	})
}

func TestEnterpriseController_GetEnterpriseByStatus(t *testing.T) {
	mockEnterpriseUsecase := new(mocks.EnterpriseUsecase)
	mockRatingUsecase := new(mocks.RatingUsecase)
	mockAuthUsecase := new(mocks.AuthUsecase)

	t.Run("get draft list", func(t *testing.T) {
		e := echo.New()
		req, rec := makeRequestHttp("", echo.GET, "/enterprise/draft", true, true)
		c := e.NewContext(req, rec)
		c.SetPath(base_path + "enterprise/:status")
		c.SetParamNames("status")
		c.SetParamValues("draft")
		enterpriseController := http2.NewEnterpriseController(mockAuthUsecase, mockEnterpriseUsecase, mockRatingUsecase)
		mockEnterpriseUsecase.On("GetListEnterpriseByStatus", mock.Anything).Return(domain.Enterprises{dummyEnterprise[0]}, nil).Once()
		mockAuthUsecase.On("GetUserDetails", mock.Anything).Return(dummyUser[0], domain.Favorite{}, domain.Enterprises{}, nil).Once()
		mockRatingUsecase.On("GetAllRatingByEnterpriseID", mock.Anything).Return(domain.RatingEnterprises{dummyRating[0]}, nil).Once()
		err := middlewareToken(enterpriseController.GetEnterpriseByStatus, c)
		responseBody := parseResponse(rec)
		assert.NoError(t, err)
		assert.Equal(t, 200, int(responseBody["code"].(float64)))
		assert.NotNil(t, responseBody["data"])
		mockEnterpriseUsecase.AssertExpectations(t)
	})

	t.Run("get publish list", func(t *testing.T) {
		e := echo.New()
		req, rec := makeRequestHttp("", echo.GET, "/enterprise/publish", true, true)
		c := e.NewContext(req, rec)
		c.SetPath(base_path + "enterprise/:status")
		c.SetParamNames("status")
		c.SetParamValues("publish")
		enterpriseController := http2.NewEnterpriseController(mockAuthUsecase, mockEnterpriseUsecase, mockRatingUsecase)
		mockEnterpriseUsecase.On("GetListEnterpriseByStatus", mock.Anything).Return(domain.Enterprises{dummyEnterprise[0]}, nil).Once()
		mockAuthUsecase.On("GetUserDetails", mock.Anything).Return(dummyUser[0], domain.Favorite{}, domain.Enterprises{}, nil).Once()
		mockRatingUsecase.On("GetAllRatingByEnterpriseID", mock.Anything).Return(domain.RatingEnterprises{dummyRating[0]}, nil).Once()
		err := middlewareToken(enterpriseController.GetEnterpriseByStatus, c)
		responseBody := parseResponse(rec)
		assert.NoError(t, err)
		assert.Equal(t, 200, int(responseBody["code"].(float64)))
		assert.NotNil(t, responseBody["data"])
		mockEnterpriseUsecase.AssertExpectations(t)
	})

	t.Run("param invalid", func(t *testing.T) {
		e := echo.New()
		req, rec := makeRequestHttp("", echo.GET, "/enterprise/ststtss", true, true)
		c := e.NewContext(req, rec)
		c.SetPath(base_path + "enterprise/:status")
		c.SetParamNames("status")
		c.SetParamValues("ststtss")
		enterpriseController := http2.NewEnterpriseController(mockAuthUsecase, mockEnterpriseUsecase, mockRatingUsecase)
		err := middlewareToken(enterpriseController.GetEnterpriseByStatus, c)
		responseBody := parseResponse(rec)
		assert.NoError(t, err)
		assert.Equal(t, float64(400), responseBody["code"])
		mockEnterpriseUsecase.AssertExpectations(t)
	})
	t.Run("error get by status", func(t *testing.T) {
		e := echo.New()
		req, rec := makeRequestHttp("", echo.GET, "/enterprise/draft", true, true)
		c := e.NewContext(req, rec)
		c.SetPath(base_path + "enterprise/:status")
		c.SetParamNames("status")
		c.SetParamValues("draft")
		enterpriseController := http2.NewEnterpriseController(mockAuthUsecase, mockEnterpriseUsecase, mockRatingUsecase)
		mockEnterpriseUsecase.On("GetListEnterpriseByStatus", mock.Anything).Return(domain.Enterprises{}, errors.New("error something")).Once()
		err := middlewareToken(enterpriseController.GetEnterpriseByStatus, c)
		responseBody := parseResponse(rec)
		assert.NoError(t, err)
		assert.Equal(t, float64(400), responseBody["code"])
		mockEnterpriseUsecase.AssertExpectations(t)
	})
	t.Run("error get rantings", func(t *testing.T) {
		e := echo.New()
		req, rec := makeRequestHttp("", echo.GET, "/enterprise/draft", true, true)
		c := e.NewContext(req, rec)
		c.SetPath(base_path + "enterprise/:status")
		c.SetParamNames("status")
		c.SetParamValues("draft")
		enterpriseController := http2.NewEnterpriseController(mockAuthUsecase, mockEnterpriseUsecase, mockRatingUsecase)
		mockEnterpriseUsecase.On("GetListEnterpriseByStatus", mock.Anything).Return(domain.Enterprises{dummyEnterprise[0]}, nil).Once()
		mockAuthUsecase.On("GetUserDetails", mock.Anything).Return(dummyUser[0], domain.Favorite{}, domain.Enterprises{}, nil).Once()
		mockRatingUsecase.On("GetAllRatingByEnterpriseID", mock.Anything).Return(domain.RatingEnterprises{}, errors.New("error something")).Once()
		err := middlewareToken(enterpriseController.GetEnterpriseByStatus, c)
		responseBody := parseResponse(rec)
		assert.NoError(t, err)
		assert.Equal(t, float64(400), responseBody["code"])
		mockEnterpriseUsecase.AssertExpectations(t)
	})
	t.Run("get list all without rating", func(t *testing.T) {
		e := echo.New()
		req, rec := makeRequestHttp("", echo.GET, "/enterprise/draft", true, true)
		c := e.NewContext(req, rec)
		c.SetPath(base_path + "enterprise/:status")
		c.SetParamNames("status")
		c.SetParamValues("draft")
		enterpriseController := http2.NewEnterpriseController(mockAuthUsecase, mockEnterpriseUsecase, mockRatingUsecase)
		mockEnterpriseUsecase.On("GetListEnterpriseByStatus", mock.Anything).Return(domain.Enterprises{dummyEnterprise[0]}, nil).Once()
		mockAuthUsecase.On("GetUserDetails", mock.Anything).Return(dummyUser[0], domain.Favorite{}, domain.Enterprises{}, nil).Once()
		mockRatingUsecase.On("GetAllRatingByEnterpriseID", mock.Anything).Return(domain.RatingEnterprises{}, nil).Once()
		err := middlewareToken(enterpriseController.GetEnterpriseByStatus, c)
		responseBody := parseResponse(rec)
		assert.NoError(t, err)
		assert.Equal(t, 200, int(responseBody["code"].(float64)))
		assert.NotNil(t, responseBody["data"])
		mockEnterpriseUsecase.AssertExpectations(t)
	})

}

func TestEnterpriseController_UpdateStatusEnterprise(t *testing.T) {
	mockEnterpriseUsecase := new(mocks.EnterpriseUsecase)
	mockRatingUsecase := new(mocks.RatingUsecase)
	mockAuthUsecase := new(mocks.AuthUsecase)

	t.Run("success", func(t *testing.T) {
		e := echo.New()
		req, rec := makeRequestHttp("", echo.PUT, "/enterprise/"+dummyEnterprise[0].ID.String()+"/status?status=1", true, true)
		c := e.NewContext(req, rec)
		c.SetPath(base_path + "enterprise/:id/status")
		c.SetParamNames("id")
		c.SetParamValues(dummyEnterprise[0].ID.String())
		enterpriseController := http2.NewEnterpriseController(mockAuthUsecase, mockEnterpriseUsecase, mockRatingUsecase)
		mockAuthUsecase.On("CheckIfUserIsAdmin", mock.Anything).Return(true, nil).Once()
		mockEnterpriseUsecase.On("UpdateStatusEnterprise", mock.Anything, mock.Anything).Return(dummyEnterprise[0], nil).Once()
		mockEnterpriseUsecase.On("GetDetailEnterpriseByID", mock.Anything).Return(dummyEnterprise[0], nil).Once()
		err := middlewareToken(enterpriseController.UpdateStatusEnterprise, c)
		responseBody := parseResponse(rec)
		assert.NoError(t, err)
		assert.Equal(t, 200, int(responseBody["code"].(float64)))
		assert.NotNil(t, responseBody["data"])
		mockEnterpriseUsecase.AssertExpectations(t)
	})

	t.Run("not admin", func(t *testing.T) {
		e := echo.New()
		req, rec := makeRequestHttp("", echo.PUT, "/enterprise/"+dummyEnterprise[0].ID.String()+"/status?status=1", true, true)
		c := e.NewContext(req, rec)
		c.SetPath(base_path + "enterprise/:id/status")
		c.SetParamNames("id")
		c.SetParamValues(dummyEnterprise[0].ID.String())
		enterpriseController := http2.NewEnterpriseController(mockAuthUsecase, mockEnterpriseUsecase, mockRatingUsecase)
		mockAuthUsecase.On("CheckIfUserIsAdmin", mock.Anything).Return(bool(true), errors.New("error something")).Once()
		err := middlewareToken(enterpriseController.UpdateStatusEnterprise, c)
		_ = parseResponse(rec)
		assert.NoError(t, err)
		mockEnterpriseUsecase.AssertExpectations(t)
	})

	t.Run("failed update status", func(t *testing.T) {
		e := echo.New()
		req, rec := makeRequestHttp("", echo.PUT, "/enterprise/"+dummyEnterprise[0].ID.String()+"/status?status=1", true, true)
		c := e.NewContext(req, rec)
		c.SetPath(base_path + "enterprise/:id/status")
		c.SetParamNames("id")
		c.SetParamValues(dummyEnterprise[0].ID.String())
		enterpriseController := http2.NewEnterpriseController(mockAuthUsecase, mockEnterpriseUsecase, mockRatingUsecase)
		mockAuthUsecase.On("CheckIfUserIsAdmin", mock.Anything).Return(true, nil).Once()
		mockEnterpriseUsecase.On("UpdateStatusEnterprise", mock.Anything, mock.Anything).Return(domain.Enterprise{}, errors.New("error something")).Once()
		err := middlewareToken(enterpriseController.UpdateStatusEnterprise, c)
		responseBody := parseResponse(rec)
		assert.NoError(t, err)
		assert.Equal(t, 400, int(responseBody["code"].(float64)))
		mockEnterpriseUsecase.AssertExpectations(t)
	})

	t.Run("Failed get detail", func(t *testing.T) {
		e := echo.New()
		req, rec := makeRequestHttp("", echo.PUT, "/enterprise/"+dummyEnterprise[0].ID.String()+"/status?status=1", true, true)
		c := e.NewContext(req, rec)
		c.SetPath(base_path + "enterprise/:id/status")
		c.SetParamNames("id")
		c.SetParamValues(dummyEnterprise[0].ID.String())
		enterpriseController := http2.NewEnterpriseController(mockAuthUsecase, mockEnterpriseUsecase, mockRatingUsecase)
		mockAuthUsecase.On("CheckIfUserIsAdmin", mock.Anything).Return(true, nil).Once()
		mockEnterpriseUsecase.On("UpdateStatusEnterprise", mock.Anything, mock.Anything).Return(dummyEnterprise[0], nil).Once()
		mockEnterpriseUsecase.On("GetDetailEnterpriseByID", mock.Anything).Return(domain.Enterprise{}, errors.New("error something")).Once()
		err := middlewareToken(enterpriseController.UpdateStatusEnterprise, c)
		responseBody := parseResponse(rec)
		assert.NoError(t, err)
		assert.Equal(t, 400, int(responseBody["code"].(float64)))
		mockEnterpriseUsecase.AssertExpectations(t)
	})

}

func TestEnterpriseController_UpdateEnterpriseByID(t *testing.T) {
	mockEnterpriseUsecase := new(mocks.EnterpriseUsecase)
	mockRatingUsecase := new(mocks.RatingUsecase)
	mockAuthUsecase := new(mocks.AuthUsecase)
	reqCreate := request.CreateEnterpriseRequest{
		Name:        "enterprise satu",
		NumberPhone: "0012798232",
		Address:     "bjb",
		Postcode:    707722,
		Latitude:    "4235,23",
		Longitude:   "4225,233231",
		Description: "testing1",
		Tags:        []string{uuid.FromStringOrNil("35d6a9a1-aa5e-41f1-9991-08878dfdf89b").String()},
	}
	requestCreate, _ := json.Marshal(reqCreate)
	t.Run("success", func(t *testing.T) {
		e := echo.New()
		req, rec := makeRequestHttp(string(requestCreate), echo.PUT, "/enterprise/"+dummyEnterprise[0].ID.String(), true, true)
		c := e.NewContext(req, rec)
		c.SetPath(base_path + "enterprise/:id")
		c.SetParamNames("id")
		c.SetParamValues(dummyEnterprise[0].ID.String())
		mockEnterpriseUsecase.On("UpdateEnterpriseByID", mock.Anything, mock.Anything, mock.Anything).Return(dummyEnterprise[0], nil).Once()
		mockEnterpriseUsecase.On("GetDetailEnterpriseByID", mock.Anything).Return(dummyEnterprise[0], nil).Once()
		enterpriseController := http2.NewEnterpriseController(mockAuthUsecase, mockEnterpriseUsecase, mockRatingUsecase)
		err := middlewareToken(enterpriseController.UpdateEnterpriseByID, c)
		responseBody := parseResponse(rec)
		assert.NoError(t, err)
		assert.Equal(t, 200, int(responseBody["code"].(float64)))
		assert.NotNil(t, responseBody["data"])
		mockEnterpriseUsecase.AssertExpectations(t)
	})
	t.Run("failed bind", func(t *testing.T) {
		e := echo.New()
		req, rec := makeRequestHttp(string(requestCreate), echo.PUT, "/enterprise/"+dummyEnterprise[0].ID.String(), true, false)
		c := e.NewContext(req, rec)
		c.SetPath(base_path + "enterprise/:id")
		c.SetParamNames("id")
		c.SetParamValues(dummyEnterprise[0].ID.String())
		enterpriseController := http2.NewEnterpriseController(mockAuthUsecase, mockEnterpriseUsecase, mockRatingUsecase)
		err := middlewareToken(enterpriseController.UpdateEnterpriseByID, c)
		assert.NoError(t, err)
		mockEnterpriseUsecase.AssertExpectations(t)
	})
	t.Run("failed update", func(t *testing.T) {
		e := echo.New()
		req, rec := makeRequestHttp(string(requestCreate), echo.PUT, "/enterprise/"+dummyEnterprise[0].ID.String(), true, true)
		c := e.NewContext(req, rec)
		c.SetPath(base_path + "enterprise/:id")
		c.SetParamNames("id")
		c.SetParamValues(dummyEnterprise[0].ID.String())
		enterpriseController := http2.NewEnterpriseController(mockAuthUsecase, mockEnterpriseUsecase, mockRatingUsecase)
		mockEnterpriseUsecase.On("UpdateEnterpriseByID", mock.Anything, mock.Anything, mock.Anything).Return(domain.Enterprise{}, errors.New("error something")).Once()
		err := middlewareToken(enterpriseController.UpdateEnterpriseByID, c)
		responseBody := parseResponse(rec)
		assert.NoError(t, err)
		assert.Equal(t, float64(400), responseBody["code"])
		mockEnterpriseUsecase.AssertExpectations(t)
	})
	t.Run("failed get by id", func(t *testing.T) {
		e := echo.New()
		req, rec := makeRequestHttp(string(requestCreate), echo.PUT, "/enterprise/"+dummyEnterprise[0].ID.String(), true, true)
		c := e.NewContext(req, rec)
		c.SetPath(base_path + "enterprise/:id")
		c.SetParamNames("id")
		c.SetParamValues(dummyEnterprise[0].ID.String())
		enterpriseController := http2.NewEnterpriseController(mockAuthUsecase, mockEnterpriseUsecase, mockRatingUsecase)
		mockEnterpriseUsecase.On("UpdateEnterpriseByID", mock.Anything, mock.Anything, mock.Anything).Return(dummyEnterprise[0], nil).Once()
		mockEnterpriseUsecase.On("GetDetailEnterpriseByID", mock.Anything).Return(domain.Enterprise{}, errors.New("error something")).Once()
		err := middlewareToken(enterpriseController.UpdateEnterpriseByID, c)
		responseBody := parseResponse(rec)
		assert.NoError(t, err)
		assert.Equal(t, float64(400), responseBody["code"])
		mockEnterpriseUsecase.AssertExpectations(t)
	})
}

func TestEnterpriseController_GetAllEnterprises(t *testing.T) {
	mockEnterpriseUsecase := new(mocks.EnterpriseUsecase)
	mockRatingUsecase := new(mocks.RatingUsecase)
	mockAuthUsecase := new(mocks.AuthUsecase)

	t.Run("success", func(t *testing.T) {
		e := echo.New()
		req, rec := makeRequestHttp("", echo.GET, "/enterprises?search=&length=1&page=1", true, true)
		c := e.NewContext(req, rec)
		enterpriseController := http2.NewEnterpriseController(mockAuthUsecase, mockEnterpriseUsecase, mockRatingUsecase)
		mockEnterpriseUsecase.On("GetListAllEnterprise", mock.Anything, mock.Anything, mock.Anything).Return(domain.Enterprises{dummyEnterprise[0]}, 1, nil).Once()
		mockAuthUsecase.On("GetUserDetails", mock.Anything).Return(dummyUser[0], domain.Favorite{}, domain.Enterprises{}, nil).Once()
		mockRatingUsecase.On("GetAllRatingByEnterpriseID", mock.Anything).Return(domain.RatingEnterprises{dummyRating[0]}, nil).Once()
		err := middlewareToken(enterpriseController.GetAllEnterprises, c)
		responseBody := parseResponse(rec)
		assert.NoError(t, err)
		assert.Equal(t, 200, int(responseBody["code"].(float64)))
		assert.NotNil(t, responseBody["data"])
		mockEnterpriseUsecase.AssertExpectations(t)
	})

	t.Run("success get list empty", func(t *testing.T) {
		e := echo.New()
		req, rec := makeRequestHttp("", echo.GET, "/enterprises?search=&length=1&page=1", true, true)
		c := e.NewContext(req, rec)
		enterpriseController := http2.NewEnterpriseController(mockAuthUsecase, mockEnterpriseUsecase, mockRatingUsecase)
		mockEnterpriseUsecase.On("GetListAllEnterprise", mock.Anything, mock.Anything, mock.Anything).Return(domain.Enterprises{}, 1, nil).Once()
		mockAuthUsecase.On("GetUserDetails", mock.Anything).Return(dummyUser[0], domain.Favorite{}, domain.Enterprises{}, nil).Once()
		mockRatingUsecase.On("GetAllRatingByEnterpriseID", mock.Anything).Return(domain.RatingEnterprises{dummyRating[0]}, nil).Once()
		err := middlewareToken(enterpriseController.GetAllEnterprises, c)
		responseBody := parseResponse(rec)
		assert.NoError(t, err)
		assert.Equal(t, 200, int(responseBody["code"].(float64)))
		mockEnterpriseUsecase.AssertExpectations(t)
	})
	t.Run("error get list", func(t *testing.T) {
		e := echo.New()
		req, rec := makeRequestHttp("", echo.GET, "/enterprises?search=&length=1&page=1", true, true)
		c := e.NewContext(req, rec)
		enterpriseController := http2.NewEnterpriseController(mockAuthUsecase, mockEnterpriseUsecase, mockRatingUsecase)
		mockEnterpriseUsecase.On("GetListAllEnterprise", mock.Anything, mock.Anything, mock.Anything).Return(domain.Enterprises{}, 1, errors.New("error something")).Once()
		err := middlewareToken(enterpriseController.GetAllEnterprises, c)
		responseBody := parseResponse(rec)
		assert.NoError(t, err)
		assert.Equal(t, 400, int(responseBody["code"].(float64)))
		mockEnterpriseUsecase.AssertExpectations(t)
	})
	t.Run("error get ratings", func(t *testing.T) {
		e := echo.New()
		req, rec := makeRequestHttp("", echo.GET, "/enterprises?search=&length=1&page=1", true, true)
		c := e.NewContext(req, rec)
		enterpriseController := http2.NewEnterpriseController(mockAuthUsecase, mockEnterpriseUsecase, mockRatingUsecase)
		mockEnterpriseUsecase.On("GetListAllEnterprise", mock.Anything, mock.Anything, mock.Anything).Return(domain.Enterprises{dummyEnterprise[0], dummyEnterprise[1]}, 2, nil).Once()
		mockAuthUsecase.On("GetUserDetails", mock.Anything).Return(dummyUser[0], domain.Favorite{}, domain.Enterprises{}, nil).Once()
		mockRatingUsecase.On("GetAllRatingByEnterpriseID", mock.Anything).Return(domain.RatingEnterprises{}, errors.New("error something")).Once()
		err := middlewareToken(enterpriseController.GetAllEnterprises, c)
		responseBody := parseResponse(rec)
		assert.NoError(t, err)
		assert.Equal(t, 400, int(responseBody["code"].(float64)))
		mockEnterpriseUsecase.AssertExpectations(t)
	})

	t.Run("get list page 0", func(t *testing.T) {
		e := echo.New()
		req, rec := makeRequestHttp("", echo.GET, "/enterprises?search=&length=1&page=0", true, true)
		c := e.NewContext(req, rec)
		enterpriseController := http2.NewEnterpriseController(mockAuthUsecase, mockEnterpriseUsecase, mockRatingUsecase)
		mockEnterpriseUsecase.On("GetListAllEnterprise", mock.Anything, mock.Anything, mock.Anything).Return(domain.Enterprises{dummyEnterprise[0]}, 1, nil).Once()
		mockAuthUsecase.On("GetUserDetails", mock.Anything).Return(dummyUser[0], domain.Favorite{}, domain.Enterprises{}, nil).Once()
		mockRatingUsecase.On("GetAllRatingByEnterpriseID", mock.Anything).Return(domain.RatingEnterprises{}, nil).Once()
		err := middlewareToken(enterpriseController.GetAllEnterprises, c)
		responseBody := parseResponse(rec)
		assert.NoError(t, err)
		assert.Equal(t, 200, int(responseBody["code"].(float64)))
		mockEnterpriseUsecase.AssertExpectations(t)
	})
}

func TestEnterpriseController_DeleteEnterpriseByID(t *testing.T) {
	mockEnterpriseUsecase := new(mocks.EnterpriseUsecase)
	mockRatingUsecase := new(mocks.RatingUsecase)
	mockAuthUsecase := new(mocks.AuthUsecase)

	t.Run("success", func(t *testing.T) {
		e := echo.New()
		req, rec := makeRequestHttp("", echo.DELETE, "/enterprise/"+dummyEnterprise[0].ID.String(), true, true)
		c := e.NewContext(req, rec)
		c.SetPath(base_path + "enterprise/:id")
		c.SetParamNames("id")
		enterpriseController := http2.NewEnterpriseController(mockAuthUsecase, mockEnterpriseUsecase, mockRatingUsecase)
		mockAuthUsecase.On("CheckIfUserIsAdmin", mock.Anything).Return(true, nil).Once()
		mockEnterpriseUsecase.On("GetDetailEnterpriseByID", mock.Anything).Return(dummyEnterprise[0], nil).Once()
		mockEnterpriseUsecase.On("DeleteEnterpriseByID", mock.Anything).Return(nil).Once()
		err := middlewareToken(enterpriseController.DeleteEnterpriseByID, c)
		responseBody := parseResponse(rec)
		assert.NoError(t, err)
		assert.Equal(t, 200, int(responseBody["code"].(float64)))
		mockEnterpriseUsecase.AssertExpectations(t)
	})

	t.Run("not current user", func(t *testing.T) {
		e := echo.New()
		req, rec := makeRequestHttp("", echo.DELETE, "/enterprise/"+dummyEnterprise[0].ID.String(), true, true)
		c := e.NewContext(req, rec)
		c.SetPath(base_path + "enterprise/:id")
		c.SetParamNames("id")
		enterpriseController := http2.NewEnterpriseController(mockAuthUsecase, mockEnterpriseUsecase, mockRatingUsecase)
		mockAuthUsecase.On("CheckIfUserIsAdmin", mock.Anything).Return(false, nil).Once()
		mockEnterpriseUsecase.On("GetDetailEnterpriseByID", mock.Anything).Return(dummyEnterprise[1], nil).Once()
		err := middlewareToken(enterpriseController.DeleteEnterpriseByID, c)
		responseBody := parseResponse(rec)
		assert.NoError(t, err)
		assert.Equal(t, 401, int(responseBody["code"].(float64)))
		mockEnterpriseUsecase.AssertExpectations(t)
	})

	t.Run("error checking admin", func(t *testing.T) {
		e := echo.New()
		req, rec := makeRequestHttp("", echo.DELETE, "/enterprise/"+dummyEnterprise[0].ID.String(), true, true)
		c := e.NewContext(req, rec)
		c.SetPath(base_path + "enterprise/:id")
		c.SetParamNames("id")
		enterpriseController := http2.NewEnterpriseController(mockAuthUsecase, mockEnterpriseUsecase, mockRatingUsecase)
		mockAuthUsecase.On("CheckIfUserIsAdmin", mock.Anything).Return(false, errors.New("error something")).Once()
		err := middlewareToken(enterpriseController.DeleteEnterpriseByID, c)
		responseBody := parseResponse(rec)
		assert.NoError(t, err)
		assert.Equal(t, 400, int(responseBody["code"].(float64)))
		mockEnterpriseUsecase.AssertExpectations(t)
	})

	t.Run("error get detail user", func(t *testing.T) {
		e := echo.New()
		req, rec := makeRequestHttp("", echo.DELETE, "/enterprise/"+dummyEnterprise[0].ID.String(), true, true)
		c := e.NewContext(req, rec)
		c.SetPath(base_path + "enterprise/:id")
		c.SetParamNames("id")
		enterpriseController := http2.NewEnterpriseController(mockAuthUsecase, mockEnterpriseUsecase, mockRatingUsecase)
		mockAuthUsecase.On("CheckIfUserIsAdmin", mock.Anything).Return(true, nil).Once()
		mockEnterpriseUsecase.On("GetDetailEnterpriseByID", mock.Anything).Return(domain.Enterprise{}, errors.New("error something")).Once()
		err := middlewareToken(enterpriseController.DeleteEnterpriseByID, c)
		responseBody := parseResponse(rec)
		assert.NoError(t, err)
		assert.Equal(t, 400, int(responseBody["code"].(float64)))
		mockEnterpriseUsecase.AssertExpectations(t)
	})

	t.Run("error delete", func(t *testing.T) {
		e := echo.New()
		req, rec := makeRequestHttp("", echo.DELETE, "/enterprise/"+dummyEnterprise[0].ID.String(), true, true)
		c := e.NewContext(req, rec)
		c.SetPath(base_path + "enterprise/:id")
		c.SetParamNames("id")
		enterpriseController := http2.NewEnterpriseController(mockAuthUsecase, mockEnterpriseUsecase, mockRatingUsecase)
		mockAuthUsecase.On("CheckIfUserIsAdmin", mock.Anything).Return(true, nil).Once()
		mockEnterpriseUsecase.On("GetDetailEnterpriseByID", mock.Anything).Return(dummyEnterprise[0], nil).Once()
		mockEnterpriseUsecase.On("DeleteEnterpriseByID", mock.Anything).Return(errors.New("error something")).Once()
		err := middlewareToken(enterpriseController.DeleteEnterpriseByID, c)
		responseBody := parseResponse(rec)
		assert.NoError(t, err)
		assert.Equal(t, 400, int(responseBody["code"].(float64)))
		mockEnterpriseUsecase.AssertExpectations(t)
	})
}

func TestEnterpriseController_GetDetailEnterpriseByID(t *testing.T) {
	mockEnterpriseUsecase := new(mocks.EnterpriseUsecase)
	mockRatingUsecase := new(mocks.RatingUsecase)
	mockAuthUsecase := new(mocks.AuthUsecase)
	t.Run("success", func(t *testing.T) {
		e := echo.New()
		req, rec := makeRequestHttp("", echo.GET, "/enterprise/"+dummyEnterprise[0].ID.String(), true, true)
		c := e.NewContext(req, rec)
		c.SetPath(base_path + "enterprise/:id")
		c.SetParamNames("id")
		c.SetParamValues(dummyEnterprise[0].ID.String())
		enterpriseController := http2.NewEnterpriseController(mockAuthUsecase, mockEnterpriseUsecase, mockRatingUsecase)
		mockEnterpriseUsecase.On("GetDetailEnterpriseByID", mock.Anything).Return(dummyEnterprise[0], nil).Once()
		mockRatingUsecase.On("GetAllRatingByEnterpriseID", mock.Anything).Return(domain.RatingEnterprises{dummyRating[0]}, nil).Once()
		mockAuthUsecase.On("GetUserDetails", mock.Anything).Return(dummyUser[0], domain.Favorite{}, domain.Enterprises{}, nil).Once()
		err := middlewareToken(enterpriseController.GetDetailEnterpriseByID, c)
		responseBody := parseResponse(rec)
		assert.NoError(t, err)
		assert.Equal(t, 200, int(responseBody["code"].(float64)))
		mockEnterpriseUsecase.AssertExpectations(t)
	})
	t.Run("error not found enterprise", func(t *testing.T) {
		e := echo.New()
		req, rec := makeRequestHttp("", echo.GET, "/enterprise/"+dummyEnterprise[0].ID.String(), true, true)
		c := e.NewContext(req, rec)
		c.SetPath(base_path + "enterprise/:id")
		c.SetParamNames("id")
		c.SetParamValues(dummyEnterprise[0].ID.String())
		enterpriseController := http2.NewEnterpriseController(mockAuthUsecase, mockEnterpriseUsecase, mockRatingUsecase)
		mockEnterpriseUsecase.On("GetDetailEnterpriseByID", mock.Anything).Return(dummyEnterprise[1], nil).Once()
		mockRatingUsecase.On("GetAllRatingByEnterpriseID", mock.Anything).Return(domain.RatingEnterprises{dummyRating[0]}, nil).Once()
		mockAuthUsecase.On("GetUserDetails", mock.Anything).Return(dummyUser[0], domain.Favorite{}, domain.Enterprises{}, nil).Once()
		err := middlewareToken(enterpriseController.GetDetailEnterpriseByID, c)
		responseBody := parseResponse(rec)
		assert.NoError(t, err)
		assert.Equal(t, 404, int(responseBody["code"].(float64)))
		mockEnterpriseUsecase.AssertExpectations(t)
	})
	t.Run("get without ratings", func(t *testing.T) {
		e := echo.New()
		req, rec := makeRequestHttp("", echo.GET, "/enterprise/"+dummyEnterprise[0].ID.String(), true, true)
		c := e.NewContext(req, rec)
		c.SetPath(base_path + "enterprise/:id")
		c.SetParamNames("id")
		c.SetParamValues(dummyEnterprise[0].ID.String())
		enterpriseController := http2.NewEnterpriseController(mockAuthUsecase, mockEnterpriseUsecase, mockRatingUsecase)
		mockEnterpriseUsecase.On("GetDetailEnterpriseByID", mock.Anything).Return(dummyEnterprise[0], nil).Once()
		mockRatingUsecase.On("GetAllRatingByEnterpriseID", mock.Anything).Return(domain.RatingEnterprises{}, nil).Once()
		mockAuthUsecase.On("GetUserDetails", mock.Anything).Return(dummyUser[0], domain.Favorite{}, domain.Enterprises{}, nil).Once()
		err := middlewareToken(enterpriseController.GetDetailEnterpriseByID, c)
		responseBody := parseResponse(rec)
		assert.NoError(t, err)
		assert.Equal(t, 200, int(responseBody["code"].(float64)))
		mockEnterpriseUsecase.AssertExpectations(t)
	})
}

func TestEnterpriseController_GetDistance(t *testing.T) {
	mockEnterpriseUsecase := new(mocks.EnterpriseUsecase)
	mockRatingUsecase := new(mocks.RatingUsecase)
	mockAuthUsecase := new(mocks.AuthUsecase)

	t.Run("success", func(t *testing.T) {
		e := echo.New()
		req, rec := makeRequestHttp("", echo.GET, "/enterprise/"+dummyEnterprise[0].ID.String()+"distance?longitude=114.740106&latitude=-3.442821", true, true)
		c := e.NewContext(req, rec)
		c.SetPath(base_path + "enterprise/:id")
		c.SetParamNames("id")
		c.SetParamValues(dummyEnterprise[0].ID.String())
		enterpriseController := http2.NewEnterpriseController(mockAuthUsecase, mockEnterpriseUsecase, mockRatingUsecase)
		mockEnterpriseUsecase.On("GetDetailEnterpriseByID", mock.Anything).Return(dummyEnterprise[0], nil).Once()
		err := middlewareToken(enterpriseController.GetDistance, c)
		responseBody := parseResponse(rec)
		assert.NoError(t, err)
		assert.Equal(t, 200, int(responseBody["code"].(float64)))
		mockEnterpriseUsecase.AssertExpectations(t)
	})

	t.Run("error get enterprise", func(t *testing.T) {
		e := echo.New()
		req, rec := makeRequestHttp("", echo.GET, "/enterprise/"+dummyEnterprise[0].ID.String()+"distance?longitude=114.740106&latitude=-3.442821", true, true)
		c := e.NewContext(req, rec)
		c.SetPath(base_path + "enterprise/:id")
		c.SetParamNames("id")
		c.SetParamValues(dummyEnterprise[0].ID.String())
		enterpriseController := http2.NewEnterpriseController(mockAuthUsecase, mockEnterpriseUsecase, mockRatingUsecase)
		mockEnterpriseUsecase.On("GetDetailEnterpriseByID", mock.Anything).Return(dummyEnterprise[0], errors.New("error something")).Once()
		err := middlewareToken(enterpriseController.GetDistance, c)
		responseBody := parseResponse(rec)
		assert.NoError(t, err)
		assert.Equal(t, 400, int(responseBody["code"].(float64)))
		mockEnterpriseUsecase.AssertExpectations(t)
	})
}

func TestEnterpriseController_AddNewRanting(t *testing.T) {
	mockEnterpriseUsecase := new(mocks.EnterpriseUsecase)
	mockRatingUsecase := new(mocks.RatingUsecase)
	mockAuthUsecase := new(mocks.AuthUsecase)
	t.Run("success", func(t *testing.T) {
		e := echo.New()
		req, rec := makeRequestHttp("", echo.POST, "/enterprise/"+dummyEnterprise[0].ID.String()+"/rating?value=3", true, true)
		c := e.NewContext(req, rec)
		c.SetPath(base_path + "enterprise/:id/rating")
		c.SetParamNames("id")
		c.SetParamValues(dummyEnterprise[0].ID.String())
		enterpriseController := http2.NewEnterpriseController(mockAuthUsecase, mockEnterpriseUsecase, mockRatingUsecase)
		mockRatingUsecase.On("FindRating", mock.Anything, mock.Anything).Return(domain.RatingEnterprise{}, nil).Once()
		mockRatingUsecase.On("AddNewRanting", mock.Anything, mock.Anything, mock.Anything).Return(dummyRating[0], nil).Once()
		mockRatingUsecase.On("GetAllRatingByEnterpriseID", mock.Anything).Return(domain.RatingEnterprises{dummyRating[0]}, nil).Once()
		err := middlewareToken(enterpriseController.AddNewRanting, c)
		responseBody := parseResponse(rec)
		assert.NoError(t, err)
		assert.Equal(t, 201, int(responseBody["code"].(float64)))
		mockRatingUsecase.AssertExpectations(t)
	})
	t.Run("error rating found", func(t *testing.T) {
		e := echo.New()
		req, rec := makeRequestHttp("", echo.POST, "/enterprise/"+dummyEnterprise[0].ID.String()+"/rating?value=3", true, true)
		c := e.NewContext(req, rec)
		c.SetPath(base_path + "enterprise/:id/rating")
		c.SetParamNames("id")
		c.SetParamValues(dummyEnterprise[0].ID.String())
		enterpriseController := http2.NewEnterpriseController(mockAuthUsecase, mockEnterpriseUsecase, mockRatingUsecase)
		mockRatingUsecase.On("FindRating", mock.Anything, mock.Anything).Return(dummyRating[0], nil).Once()
		err := middlewareToken(enterpriseController.AddNewRanting, c)
		responseBody := parseResponse(rec)
		assert.NoError(t, err)
		assert.Equal(t, 400, int(responseBody["code"].(float64)))
		mockRatingUsecase.AssertExpectations(t)
	})
	t.Run("failed error", func(t *testing.T) {
		e := echo.New()
		req, rec := makeRequestHttp("", echo.POST, "/enterprise/"+dummyEnterprise[0].ID.String()+"/rating?value=3", true, true)
		c := e.NewContext(req, rec)
		c.SetPath(base_path + "enterprise/:id/rating")
		c.SetParamNames("id")
		c.SetParamValues(dummyEnterprise[0].ID.String())
		enterpriseController := http2.NewEnterpriseController(mockAuthUsecase, mockEnterpriseUsecase, mockRatingUsecase)
		mockRatingUsecase.On("FindRating", mock.Anything, mock.Anything).Return(domain.RatingEnterprise{}, nil).Once()
		mockRatingUsecase.On("AddNewRanting", mock.Anything, mock.Anything, mock.Anything).Return(dummyRating[0], errors.New("error something")).Once()
		err := middlewareToken(enterpriseController.AddNewRanting, c)
		responseBody := parseResponse(rec)
		assert.NoError(t, err)
		assert.Equal(t, 400, int(responseBody["code"].(float64)))
		mockRatingUsecase.AssertExpectations(t)
	})
	t.Run("error get rating", func(t *testing.T) {
		e := echo.New()
		req, rec := makeRequestHttp("", echo.POST, "/enterprise/"+dummyEnterprise[0].ID.String()+"/rating?value=3", true, true)
		c := e.NewContext(req, rec)
		c.SetPath(base_path + "enterprise/:id/rating")
		c.SetParamNames("id")
		c.SetParamValues(dummyEnterprise[0].ID.String())
		enterpriseController := http2.NewEnterpriseController(mockAuthUsecase, mockEnterpriseUsecase, mockRatingUsecase)
		mockRatingUsecase.On("FindRating", mock.Anything, mock.Anything).Return(domain.RatingEnterprise{}, nil).Once()
		mockRatingUsecase.On("AddNewRanting", mock.Anything, mock.Anything, mock.Anything).Return(dummyRating[0], nil).Once()
		mockRatingUsecase.On("GetAllRatingByEnterpriseID", mock.Anything).Return(domain.RatingEnterprises{}, errors.New("error something")).Once()
		err := middlewareToken(enterpriseController.AddNewRanting, c)
		responseBody := parseResponse(rec)
		assert.NoError(t, err)
		assert.Equal(t, 400, int(responseBody["code"].(float64)))
		mockRatingUsecase.AssertExpectations(t)
	})
}

func TestEnterpriseController_CekRatingUser(t *testing.T) {
	mockEnterpriseUsecase := new(mocks.EnterpriseUsecase)
	mockRatingUsecase := new(mocks.RatingUsecase)
	mockAuthUsecase := new(mocks.AuthUsecase)
	t.Run("success but 0 rating", func(t *testing.T) {
		e := echo.New()
		req, rec := makeRequestHttp("", echo.POST, "/enterprise/"+dummyEnterprise[0].ID.String()+"/rating/user/"+dummyEnterprise[0].UserID.String(), true, true)
		c := e.NewContext(req, rec)
		c.SetPath(base_path + "enterprise/:id/rating/user/:userid")
		c.SetParamNames("id")
		c.SetParamValues(dummyEnterprise[0].ID.String())
		enterpriseController := http2.NewEnterpriseController(mockAuthUsecase, mockEnterpriseUsecase, mockRatingUsecase)
		mockRatingUsecase.On("FindRating", mock.Anything, mock.Anything).Return(dummyRating[1], nil).Once()
		err := middlewareToken(enterpriseController.CekRatingUser, c)
		responseBody := parseResponse(rec)
		assert.NoError(t, err)
		assert.Equal(t, 200, int(responseBody["code"].(float64)))
		mockRatingUsecase.AssertExpectations(t)
	})

	t.Run("success", func(t *testing.T) {
		e := echo.New()
		req, rec := makeRequestHttp("", echo.POST, "/enterprise/"+dummyEnterprise[0].ID.String()+"/rating/user/"+dummyEnterprise[0].UserID.String(), true, true)
		c := e.NewContext(req, rec)
		c.SetPath(base_path + "enterprise/:id/rating/user/:userid")
		c.SetParamNames("id")
		c.SetParamValues(dummyEnterprise[0].ID.String())
		enterpriseController := http2.NewEnterpriseController(mockAuthUsecase, mockEnterpriseUsecase, mockRatingUsecase)
		mockRatingUsecase.On("FindRating", mock.Anything, mock.Anything).Return(dummyRating[0], nil).Once()
		err := middlewareToken(enterpriseController.CekRatingUser, c)
		responseBody := parseResponse(rec)
		assert.NoError(t, err)
		assert.Equal(t, 200, int(responseBody["code"].(float64)))
		mockRatingUsecase.AssertExpectations(t)
	})

	t.Run("rating not found", func(t *testing.T) {
		e := echo.New()
		req, rec := makeRequestHttp("", echo.POST, "/enterprise/"+dummyEnterprise[0].ID.String()+"/rating/user/"+dummyEnterprise[0].UserID.String(), true, true)
		c := e.NewContext(req, rec)
		c.SetPath(base_path + "enterprise/:id/rating/user/:userid")
		c.SetParamNames("id")
		c.SetParamValues(dummyEnterprise[0].ID.String())
		enterpriseController := http2.NewEnterpriseController(mockAuthUsecase, mockEnterpriseUsecase, mockRatingUsecase)
		mockRatingUsecase.On("FindRating", mock.Anything, mock.Anything).Return(domain.RatingEnterprise{}, errors.New("error something")).Once()
		err := middlewareToken(enterpriseController.CekRatingUser, c)
		responseBody := parseResponse(rec)
		assert.NoError(t, err)
		assert.Equal(t, 404, int(responseBody["code"].(float64)))
		mockRatingUsecase.AssertExpectations(t)
	})
}

func TestEnterpriseController_DeleteRatingUser(t *testing.T) {
	mockEnterpriseUsecase := new(mocks.EnterpriseUsecase)
	mockRatingUsecase := new(mocks.RatingUsecase)
	mockAuthUsecase := new(mocks.AuthUsecase)
	t.Run("success", func(t *testing.T) {
		e := echo.New()
		req, rec := makeRequestHttp("", echo.DELETE, "/enterprise/"+dummyEnterprise[0].ID.String()+"/rating/user/"+dummyEnterprise[0].UserID.String(), true, true)
		c := e.NewContext(req, rec)
		c.SetPath(base_path + "enterprise/:id/rating/user/:userid")
		c.SetParamNames("id", "userid")
		c.SetParamValues(dummyEnterprise[0].ID.String(), dummyEnterprise[0].UserID.String())
		enterpriseController := http2.NewEnterpriseController(mockAuthUsecase, mockEnterpriseUsecase, mockRatingUsecase)
		mockRatingUsecase.On("FindRating", mock.Anything, mock.Anything).Return(dummyRating[0], nil).Once()
		mockAuthUsecase.On("CheckIfUserIsAdmin", mock.Anything).Return(true, nil).Once()
		mockRatingUsecase.On("DeleteRating", mock.Anything, mock.Anything).Return(nil).Once()
		err := middlewareToken(enterpriseController.DeleteRatingUser, c)
		responseBody := parseResponse(rec)
		assert.NoError(t, err)
		assert.Equal(t, 200, int(responseBody["code"].(float64)))
		mockRatingUsecase.AssertExpectations(t)
		mockAuthUsecase.AssertExpectations(t)
	})
	t.Run("error find rating", func(t *testing.T) {
		e := echo.New()
		req, rec := makeRequestHttp("", echo.DELETE, "/enterprise/"+dummyEnterprise[0].ID.String()+"/rating/user/"+dummyEnterprise[0].UserID.String(), true, true)
		c := e.NewContext(req, rec)
		c.SetPath(base_path + "enterprise/:id/rating/user/:userid")
		c.SetParamNames("id", "userid")
		c.SetParamValues(dummyEnterprise[0].ID.String(), dummyEnterprise[0].UserID.String())
		enterpriseController := http2.NewEnterpriseController(mockAuthUsecase, mockEnterpriseUsecase, mockRatingUsecase)
		mockRatingUsecase.On("FindRating", mock.Anything, mock.Anything).Return(domain.RatingEnterprise{}, errors.New("error something")).Once()
		err := middlewareToken(enterpriseController.DeleteRatingUser, c)
		responseBody := parseResponse(rec)
		assert.NoError(t, err)
		assert.Equal(t, 404, int(responseBody["code"].(float64)))
		mockRatingUsecase.AssertExpectations(t)
	})
	t.Run("error find rating", func(t *testing.T) {
		e := echo.New()
		req, rec := makeRequestHttp("", echo.DELETE, "/enterprise/"+dummyEnterprise[0].ID.String()+"/rating/user/"+dummyEnterprise[0].UserID.String(), true, true)
		c := e.NewContext(req, rec)
		c.SetPath(base_path + "enterprise/:id/rating/user/:userid")
		c.SetParamNames("id", "userid")
		c.SetParamValues(dummyEnterprise[0].ID.String(), dummyEnterprise[0].UserID.String())
		enterpriseController := http2.NewEnterpriseController(mockAuthUsecase, mockEnterpriseUsecase, mockRatingUsecase)
		mockRatingUsecase.On("FindRating", mock.Anything, mock.Anything).Return(dummyRating[0], nil).Once()
		mockAuthUsecase.On("CheckIfUserIsAdmin", mock.Anything).Return(true, errors.New("error something")).Once()
		err := middlewareToken(enterpriseController.DeleteRatingUser, c)
		responseBody := parseResponse(rec)
		assert.NoError(t, err)
		assert.Equal(t, 400, int(responseBody["code"].(float64)))
		mockRatingUsecase.AssertExpectations(t)
	})

	t.Run("success", func(t *testing.T) {
		e := echo.New()
		req, rec := makeRequestHttp("", echo.DELETE, "/enterprise/"+dummyEnterprise[0].ID.String()+"/rating/user/"+dummyEnterprise[0].UserID.String(), true, true)
		c := e.NewContext(req, rec)
		c.SetPath(base_path + "enterprise/:id/rating/user/:userid")
		c.SetParamNames("id", "userid")
		c.SetParamValues(dummyEnterprise[0].ID.String(), dummyEnterprise[1].UserID.String())
		enterpriseController := http2.NewEnterpriseController(mockAuthUsecase, mockEnterpriseUsecase, mockRatingUsecase)
		mockRatingUsecase.On("FindRating", mock.Anything, mock.Anything).Return(dummyRating[0], nil).Once()
		mockAuthUsecase.On("CheckIfUserIsAdmin", mock.Anything).Return(true, nil).Once()
		mockRatingUsecase.On("DeleteRating", mock.Anything, mock.Anything).Return(errors.New("error something")).Once()
		err := middlewareToken(enterpriseController.DeleteRatingUser, c)
		responseBody := parseResponse(rec)
		assert.NoError(t, err)
		assert.Equal(t, 400, int(responseBody["code"].(float64)))
		mockRatingUsecase.AssertExpectations(t)
		mockAuthUsecase.AssertExpectations(t)
	})
}

func TestEnterpriseController_UpdateRating(t *testing.T) {
	mockEnterpriseUsecase := new(mocks.EnterpriseUsecase)
	mockRatingUsecase := new(mocks.RatingUsecase)
	mockAuthUsecase := new(mocks.AuthUsecase)
	t.Run("success", func(t *testing.T) {
		e := echo.New()
		req, rec := makeRequestHttp("", echo.PUT, "/enterprise/"+dummyEnterprise[0].ID.String()+"/rating/user/"+dummyEnterprise[0].UserID.String()+"?value=3", true, true)
		c := e.NewContext(req, rec)
		c.SetPath(base_path + "enterprise/:id/rating/user/:userid")
		c.SetParamNames("id", "userid")
		c.SetParamValues(dummyEnterprise[0].ID.String(), dummyEnterprise[0].UserID.String())
		enterpriseController := http2.NewEnterpriseController(mockAuthUsecase, mockEnterpriseUsecase, mockRatingUsecase)
		mockRatingUsecase.On("FindRating", mock.Anything, mock.Anything).Return(dummyRating[0], nil).Once()
		mockAuthUsecase.On("CheckIfUserIsAdmin", mock.Anything).Return(true, nil).Once()
		mockRatingUsecase.On("UpdateRating", mock.Anything, mock.Anything, mock.Anything).Return(domain.RatingEnterprise{}, nil).Once()
		mockRatingUsecase.On("FindRating", mock.Anything, mock.Anything).Return(dummyRating[0], nil).Once()
		mockRatingUsecase.On("GetAllRatingByEnterpriseID", mock.Anything).Return(domain.RatingEnterprises{dummyRating[0]}, nil).Once()
		err := middlewareToken(enterpriseController.UpdateRating, c)
		responseBody := parseResponse(rec)
		assert.NoError(t, err)
		assert.Equal(t, 200, int(responseBody["code"].(float64)))
		mockRatingUsecase.AssertExpectations(t)
		mockAuthUsecase.AssertExpectations(t)
	})
	t.Run("success", func(t *testing.T) {
		e := echo.New()
		req, rec := makeRequestHttp("", echo.PUT, "/enterprise/"+dummyEnterprise[0].ID.String()+"/rating/user/"+dummyEnterprise[0].UserID.String()+"?value=3", true, true)
		c := e.NewContext(req, rec)
		c.SetPath(base_path + "enterprise/:id/rating/user/:userid")
		c.SetParamNames("id", "userid")
		c.SetParamValues(dummyEnterprise[0].ID.String(), dummyEnterprise[0].UserID.String())
		enterpriseController := http2.NewEnterpriseController(mockAuthUsecase, mockEnterpriseUsecase, mockRatingUsecase)
		mockRatingUsecase.On("FindRating", mock.Anything, mock.Anything).Return(dummyRating[1], nil).Once()
		mockAuthUsecase.On("CheckIfUserIsAdmin", mock.Anything).Return(true, nil).Once()
		mockRatingUsecase.On("UpdateRating", mock.Anything, mock.Anything, mock.Anything).Return(domain.RatingEnterprise{}, nil).Once()
		mockRatingUsecase.On("FindRating", mock.Anything, mock.Anything).Return(dummyRating[0], nil).Once()
		mockRatingUsecase.On("GetAllRatingByEnterpriseID", mock.Anything).Return(domain.RatingEnterprises{dummyRating[0]}, nil).Once()
		err := middlewareToken(enterpriseController.UpdateRating, c)
		responseBody := parseResponse(rec)
		assert.NoError(t, err)
		assert.Equal(t, 200, int(responseBody["code"].(float64)))
		mockRatingUsecase.AssertExpectations(t)
		mockAuthUsecase.AssertExpectations(t)
	})
}
