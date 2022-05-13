package http_test

import (
	"encoding/json"
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/nrmadi02/mini-project/domain"
	"github.com/nrmadi02/mini-project/domain/mocks"
	http2 "github.com/nrmadi02/mini-project/internal/user/delivery/http"
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

func TestAdminController_GetUserList(t *testing.T) {
	mockAuthUsecase := new(mocks.AuthUsecase)
	mockUsercase := new(mocks.UserUsecase)

	t.Run("success", func(t *testing.T) {
		e := echo.New()
		req, rec := makeRequestHttp("", echo.GET, "/users", true, true)
		c := e.NewContext(req, rec)
		adminController := http2.NewAdminController(mockAuthUsecase, mockUsercase)
		mockAuthUsecase.On("CheckIfUserIsAdmin", mock.Anything).Return(true, nil).Once()
		mockUsercase.On("GetAllUsers").Return(domain.Users{dummyUser[0]}, nil).Once()
		err := middlewareToken(adminController.GetUserList, c)
		responseBody := parseResponse(rec)
		assert.NoError(t, err)
		assert.Equal(t, float64(200), responseBody["code"])
		mockAuthUsecase.AssertExpectations(t)
		mockAuthUsecase.AssertExpectations(t)
	})
	t.Run("error not admin", func(t *testing.T) {
		e := echo.New()
		req, rec := makeRequestHttp("", echo.GET, "/users", true, true)
		c := e.NewContext(req, rec)
		adminController := http2.NewAdminController(mockAuthUsecase, mockUsercase)
		mockAuthUsecase.On("CheckIfUserIsAdmin", mock.Anything).Return(false, errors.New("error something")).Once()
		err := middlewareToken(adminController.GetUserList, c)
		responseBody := parseResponse(rec)
		assert.NoError(t, err)
		assert.Equal(t, float64(401), responseBody["code"])
		mockAuthUsecase.AssertExpectations(t)
		mockAuthUsecase.AssertExpectations(t)
	})
	t.Run("error get list users", func(t *testing.T) {
		e := echo.New()
		req, rec := makeRequestHttp("", echo.GET, "/users", true, true)
		c := e.NewContext(req, rec)
		adminController := http2.NewAdminController(mockAuthUsecase, mockUsercase)
		mockAuthUsecase.On("CheckIfUserIsAdmin", mock.Anything).Return(true, nil).Once()
		mockUsercase.On("GetAllUsers").Return(domain.Users{}, errors.New("error something")).Once()
		err := middlewareToken(adminController.GetUserList, c)
		responseBody := parseResponse(rec)
		assert.NoError(t, err)
		assert.Equal(t, float64(400), responseBody["code"])
		mockAuthUsecase.AssertExpectations(t)
		mockAuthUsecase.AssertExpectations(t)
	})
}
