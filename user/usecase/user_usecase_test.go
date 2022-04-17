package usecase_test

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go_cicd/domain"
	"go_cicd/domain/mocks"
	"go_cicd/domain/web/request"
	usecase "go_cicd/user/usecase"
	"testing"
)

func TestUserUsecase_ReadAll(t *testing.T) {
	mockUserRepo := new(mocks.UserRepository)
	mockListUser := &domain.Users{
		domain.User{
			ID:       1,
			Email:    "test@gmail.com",
			Password: "12345",
		},
		domain.User{
			ID:       2,
			Email:    "testtest@gmail.com",
			Password: "12345",
		},
	}

	t.Run("success", func(t *testing.T) {
		mockUserRepo.On("ReadAll").Return(mockListUser, nil).Once()
		uc := usecase.NewUserUsecase(mockUserRepo)
		res, err := uc.ReadAll()
		assert.NoError(t, err)
		assert.Len(t, *res, len(*mockListUser))
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("error-failed", func(t *testing.T) {
		mockUserRepo.On("ReadAll").Return(nil, errors.New("error something")).Once()
		uc := usecase.NewUserUsecase(mockUserRepo)
		_, err := uc.ReadAll()
		assert.Error(t, err)
		mockUserRepo.AssertExpectations(t)
	})
}

func TestUserUsecase_ReadByID(t *testing.T) {
	mockUserRepo := new(mocks.UserRepository)
	mockUser := &domain.User{
		ID:       1,
		Email:    "test@gmail.com",
		Password: "12345",
	}

	t.Run("success", func(t *testing.T) {
		mockUserRepo.On("ReadByID", mock.AnythingOfType("int")).Return(mockUser, nil).Once()
		uc := usecase.NewUserUsecase(mockUserRepo)
		res, err := uc.ReadByID(1)
		assert.NoError(t, err)
		assert.Equal(t, res, mockUser)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("error-failed", func(t *testing.T) {
		mockUserRepo.On("ReadByID", mock.AnythingOfType("int")).Return(nil, errors.New("error something")).Once()
		uc := usecase.NewUserUsecase(mockUserRepo)
		_, err := uc.ReadByID(2)
		assert.Error(t, err)
		mockUserRepo.AssertExpectations(t)
	})
}

func TestUserUsecase_Create(t *testing.T) {
	mockUserRepo := new(mocks.UserRepository)
	mockUser := &domain.User{
		ID:       1,
		Email:    "test@gmail.com",
		Password: "12345",
	}

	t.Run("success", func(t *testing.T) {
		mockRequest := &request.UserCreateRequest{
			Email:    "test@gmail.com",
			Password: "12345",
		}

		mockUserRepo.On("Create", mock.AnythingOfType("*domain.User")).Return(mockUser, nil).Once()
		uc := usecase.NewUserUsecase(mockUserRepo)
		res, err := uc.Create(*mockRequest)
		assert.NoError(t, err)
		assert.Equal(t, res, mockUser)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("error-failed", func(t *testing.T) {
		mockRequest := &request.UserCreateRequest{
			Email:    "test@gmail.com",
			Password: "12345",
		}

		mockUserRepo.On("Create", mock.AnythingOfType("*domain.User")).Return(nil, errors.New("error something")).Once()
		uc := usecase.NewUserUsecase(mockUserRepo)
		_, err := uc.Create(*mockRequest)
		assert.Error(t, err)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("email-empty-failed", func(t *testing.T) {
		mockRequest := &request.UserCreateRequest{
			Email:    "test@gmail.com",
			Password: "12345",
		}

		mockUserRepo.On("Create", mock.AnythingOfType("*domain.User")).Return(nil, errors.New("email empty")).Once()
		uc := usecase.NewUserUsecase(mockUserRepo)
		_, err := uc.Create(*mockRequest)
		assert.Error(t, err)
		mockUserRepo.AssertExpectations(t)
	})
}

func TestUserUsecase_Login(t *testing.T) {
	mockUserRepo := new(mocks.UserRepository)
	mockUser := &domain.User{
		ID:       1,
		Email:    "test@gmail.com",
		Password: "12345",
	}

	t.Run("success", func(t *testing.T) {
		mockRequest := &request.LoginRequest{
			Email:    "test@gmail.com",
			Password: "12345",
		}

		mockUserRepo.On("CheckLogin", mock.AnythingOfType("*domain.User")).Return(mockUser, true, nil).Once()
		uc := usecase.NewUserUsecase(mockUserRepo)
		res, err := uc.Login(*mockRequest)
		assert.NoError(t, err)
		assert.NotNil(t, res)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("error-failed", func(t *testing.T) {
		mockRequest := &request.LoginRequest{
			Email:    "test@gmail.com",
			Password: "12345",
		}

		mockUserRepo.On("CheckLogin", mock.AnythingOfType("*domain.User")).Return(nil, false, errors.New("error something")).Once()
		uc := usecase.NewUserUsecase(mockUserRepo)
		_, err := uc.Login(*mockRequest)
		assert.Error(t, err)
		mockUserRepo.AssertExpectations(t)
	})

}
