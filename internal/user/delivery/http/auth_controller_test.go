package http_test

import (
	"encoding/json"
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/nrmadi02/mini-project/domain"
	"github.com/nrmadi02/mini-project/domain/mocks"
	"github.com/nrmadi02/mini-project/internal/user/delivery/http"
	"github.com/nrmadi02/mini-project/web/request"
	"github.com/nrmadi02/mini-project/web/response"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestAuthController_Login(t *testing.T) {
	mockAuthUsecase := new(mocks.AuthUsecase)
	reqBody := request.LoginRequest{
		Email:    "satu@gmail.com",
		Password: "12345678",
	}
	requestLogin, _ := json.Marshal(reqBody)
	t.Run("success", func(t *testing.T) {
		e := echo.New()
		req, rec := makeRequestHttp(string(requestLogin), echo.POST, "/login", true, true)
		c := e.NewContext(req, rec)
		authController := http.NewAuthController(mockAuthUsecase)
		mockAuthUsecase.On("Login", mock.Anything).Return(response.SuccessLogin{
			ID:       dummyUser[0].ID,
			Email:    "satu@gmail.com",
			Fullname: "user1",
			Username: "usr1",
			Token:    createToken(),
		}, nil).Once()
		err := authController.Login(c)
		responseBody := parseResponse(rec)
		assert.NoError(t, err)
		assert.Equal(t, float64(200), responseBody["code"])
		mockAuthUsecase.AssertExpectations(t)
	})
	t.Run("error bind echo", func(t *testing.T) {
		e := echo.New()
		req, rec := makeRequestHttp(string(requestLogin), echo.POST, "/login", true, false)
		c := e.NewContext(req, rec)
		authController := http.NewAuthController(mockAuthUsecase)
		err := authController.Login(c)
		responseBody := parseResponse(rec)
		assert.NoError(t, err)
		assert.Equal(t, float64(400), responseBody["code"])
		mockAuthUsecase.AssertExpectations(t)
	})
	t.Run("success", func(t *testing.T) {
		e := echo.New()
		req, rec := makeRequestHttp(string(requestLogin), echo.POST, "/login", true, true)
		c := e.NewContext(req, rec)
		authController := http.NewAuthController(mockAuthUsecase)
		mockAuthUsecase.On("Login", mock.Anything).Return(response.SuccessLogin{}, errors.New("error something")).Once()
		err := authController.Login(c)
		responseBody := parseResponse(rec)
		assert.NoError(t, err)
		assert.Equal(t, float64(401), responseBody["code"])
		mockAuthUsecase.AssertExpectations(t)
	})
}

func TestAuthController_Register(t *testing.T) {
	mockAuthUsecase := new(mocks.AuthUsecase)
	reqBody := request.UserCreateRequest{
		Fullname: "user1",
		Username: "usr1",
		Email:    "satu@gmail.com",
		Password: "12345678",
	}
	requestRegister, _ := json.Marshal(reqBody)
	t.Run("success", func(t *testing.T) {
		e := echo.New()
		req, rec := makeRequestHttp(string(requestRegister), echo.POST, "/register ", true, true)
		c := e.NewContext(req, rec)
		authController := http.NewAuthController(mockAuthUsecase)
		mockAuthUsecase.On("Register", mock.Anything).Return(dummyUser[0], nil).Once()
		err := authController.Register(c)
		responseBody := parseResponse(rec)
		assert.NoError(t, err)
		assert.Equal(t, float64(201), responseBody["code"])
		mockAuthUsecase.AssertExpectations(t)
	})
	t.Run("success", func(t *testing.T) {
		reqBody2 := request.UserCreateRequest{
			Fullname: "",
			Username: "",
			Email:    "satu@gmail.com",
			Password: "12345678",
		}
		requestRegister2, _ := json.Marshal(reqBody2)
		e := echo.New()
		req, rec := makeRequestHttp(string(requestRegister2), echo.POST, "/register ", true, true)
		c := e.NewContext(req, rec)
		authController := http.NewAuthController(mockAuthUsecase)
		//mockAuthUsecase.On("Register", mock.Anything).Return(dummyUser[0], nil).Once()
		err := authController.Register(c)
		responseBody := parseResponse(rec)
		assert.NoError(t, err)
		assert.Equal(t, float64(400), responseBody["code"])
		mockAuthUsecase.AssertExpectations(t)
	})
	t.Run("error register", func(t *testing.T) {
		e := echo.New()
		req, rec := makeRequestHttp(string(requestRegister), echo.POST, "/register ", true, true)
		c := e.NewContext(req, rec)
		authController := http.NewAuthController(mockAuthUsecase)
		mockAuthUsecase.On("Register", mock.Anything).Return(domain.User{}, errors.New("error something")).Once()
		err := authController.Register(c)
		responseBody := parseResponse(rec)
		assert.NoError(t, err)
		assert.Equal(t, float64(400), responseBody["code"])
		mockAuthUsecase.AssertExpectations(t)
	})
	t.Run("error bind echo", func(t *testing.T) {
		e := echo.New()
		req, rec := makeRequestHttp(string(requestRegister), echo.POST, "/register ", true, false)
		c := e.NewContext(req, rec)
		authController := http.NewAuthController(mockAuthUsecase)
		err := authController.Register(c)
		responseBody := parseResponse(rec)
		assert.NoError(t, err)
		assert.Equal(t, float64(400), responseBody["code"])
		mockAuthUsecase.AssertExpectations(t)
	})
}
