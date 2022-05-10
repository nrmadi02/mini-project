package usecase_test

import (
	"errors"
	"github.com/nrmadi02/mini-project/domain"
	"github.com/nrmadi02/mini-project/domain/mocks"
	"github.com/nrmadi02/mini-project/internal/user/usecase"
	"github.com/nrmadi02/mini-project/web/request"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func TestAuthUsecase_Login(t *testing.T) {
	mockUserRepository := new(mocks.UserRepository)
	mockRoleRepository := new(mocks.RoleRepository)
	mockFavoriteRepository := new(mocks.FavoriteRepository)
	mockEnterpriseRepository := new(mocks.EnterpriseRepository)

	t.Run("success", func(t *testing.T) {
		req := request.LoginRequest{
			Email:    "satu@email.com",
			Password: "12345678",
		}
		uc := usecase.NewAuthUsecase(mockUserRepository, mockRoleRepository, mockFavoriteRepository, mockEnterpriseRepository)
		mockUserRepository.On("FindUserByEmail", mock.AnythingOfType("string")).Return(dummyUser[0], nil).Once()
		res, err := uc.Login(req)
		assert.NoError(t, err)
		assert.NotNil(t, res)
		mockUserRepository.AssertExpectations(t)
	})

	t.Run("error wrong password", func(t *testing.T) {
		req := request.LoginRequest{
			Email:    "satu@email.com",
			Password: "1234",
		}
		uc := usecase.NewAuthUsecase(mockUserRepository, mockRoleRepository, mockFavoriteRepository, mockEnterpriseRepository)
		mockUserRepository.On("FindUserByEmail", mock.AnythingOfType("string")).Return(dummyUser[0], nil).Once()
		_, err := uc.Login(req)
		assert.Error(t, err)
		mockUserRepository.AssertExpectations(t)
	})

	t.Run("error token jwt", func(t *testing.T) {
		req := request.LoginRequest{
			Email:    "sasstu@email.com",
			Password: "12345678",
		}
		uc := usecase.NewAuthUsecase(mockUserRepository, mockRoleRepository, mockFavoriteRepository, mockEnterpriseRepository)
		mockUserRepository.On("FindUserByEmail", mock.AnythingOfType("string")).Return(domain.User{Username: ""}, errors.New("")).Once()
		_, err := uc.Login(req)
		assert.Error(t, err)
		mockUserRepository.AssertExpectations(t)
	})
}

func TestAuthUsecase_Register(t *testing.T) {
	mockUserRepository := new(mocks.UserRepository)
	mockRoleRepository := new(mocks.RoleRepository)
	mockFavoriteRepository := new(mocks.FavoriteRepository)
	mockEnterpriseRepository := new(mocks.EnterpriseRepository)

	t.Run("success", func(t *testing.T) {
		req := request.UserCreateRequest{
			Fullname: "user1",
			Username: "usr1",
			Email:    "satu@email.com",
			Password: "12345678",
		}
		uc := usecase.NewAuthUsecase(mockUserRepository, mockRoleRepository, mockFavoriteRepository, mockEnterpriseRepository)
		mockUserRepository.On("FindUserByEmail", mock.AnythingOfType("string")).Return(domain.User{}, nil).Once()
		mockRoleRepository.On("FindByName", mock.AnythingOfType("string")).Return(domain.Role{Name: "admin", ID: 1}, nil).Once()
		mockUserRepository.On("Save", mock.AnythingOfType("domain.User")).Return(dummyUser[0], nil).Once()
		mockFavoriteRepository.On("Add", mock.AnythingOfType("domain.Favorite")).Return(domain.Favorite{
			ID:     uuid.NewV4(),
			UserID: dummyUser[0].ID,
		}, nil).Once()
		res, err := uc.Register(req)
		assert.NoError(t, err)
		assert.NotNil(t, res)
		mockUserRepository.AssertExpectations(t)
	})
	t.Run("user already register", func(t *testing.T) {
		req := request.UserCreateRequest{
			Fullname: "user1",
			Username: "usr1",
			Email:    "satu@email.com",
			Password: "12345678",
		}
		uc := usecase.NewAuthUsecase(mockUserRepository, mockRoleRepository, mockFavoriteRepository, mockEnterpriseRepository)
		mockUserRepository.On("FindUserByEmail", mock.AnythingOfType("string")).Return(dummyUser[1], nil).Once()
		_, err := uc.Register(req)
		assert.Error(t, err)
		mockUserRepository.AssertExpectations(t)
	})

	t.Run("error role not found", func(t *testing.T) {
		req := request.UserCreateRequest{
			Fullname: "user1",
			Username: "usr1",
			Email:    "satu@email.com",
			Password: "12345678",
		}
		uc := usecase.NewAuthUsecase(mockUserRepository, mockRoleRepository, mockFavoriteRepository, mockEnterpriseRepository)
		mockUserRepository.On("FindUserByEmail", mock.AnythingOfType("string")).Return(domain.User{}, nil).Once()
		mockRoleRepository.On("FindByName", mock.AnythingOfType("string")).Return(domain.Role{}, errors.New("role not found")).Once()
		_, err := uc.Register(req)
		assert.Error(t, err)
		mockUserRepository.AssertExpectations(t)
	})

	t.Run("error save user", func(t *testing.T) {
		req := request.UserCreateRequest{
			Fullname: "user1",
			Username: "usr1",
			Email:    "satu@email.com",
			Password: "12345678",
		}
		uc := usecase.NewAuthUsecase(mockUserRepository, mockRoleRepository, mockFavoriteRepository, mockEnterpriseRepository)
		mockUserRepository.On("FindUserByEmail", mock.AnythingOfType("string")).Return(domain.User{}, nil).Once()
		mockRoleRepository.On("FindByName", mock.AnythingOfType("string")).Return(domain.Role{Name: "admin", ID: 1}, nil).Once()
		mockUserRepository.On("Save", mock.AnythingOfType("domain.User")).Return(domain.User{}, errors.New("error save")).Once()
		_, err := uc.Register(req)
		assert.Error(t, err)
		mockUserRepository.AssertExpectations(t)
	})
}

func TestAuthUsecase_GetUserDetails(t *testing.T) {
	mockUserRepository := new(mocks.UserRepository)
	mockRoleRepository := new(mocks.RoleRepository)
	mockFavoriteRepository := new(mocks.FavoriteRepository)
	mockEnterpriseRepository := new(mocks.EnterpriseRepository)

	t.Run("success", func(t *testing.T) {
		id := uuid.FromStringOrNil("0cf712fc-e631-40c7-8572-54772e698edg")
		uc := usecase.NewAuthUsecase(mockUserRepository, mockRoleRepository, mockFavoriteRepository, mockEnterpriseRepository)
		mockUserRepository.On("FindUserById", mock.AnythingOfType("string")).Return(dummyUser[0], nil).Once()
		mockFavoriteRepository.On("FindByUserID", mock.AnythingOfType("string")).Return(domain.Favorite{
			ID:     uuid.NewV4(),
			UserID: dummyUser[0].ID,
		}, nil).Once()
		mockEnterpriseRepository.On("FindByUserID", mock.AnythingOfType("string")).Return(domain.Enterprises{
			domain.Enterprise{
				ID:               uuid.FromStringOrNil("35d6a9a1-aa5e-41f1-9991-08878dfdf89a"),
				UserID:           uuid.FromStringOrNil("35d6a9a1-aa5e-41f1-9991-08878dfdf891"),
				Name:             "enterprise satu",
				NumberPhone:      "0012798232",
				Address:          "bjb",
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
		}, nil).Once()
		res, favorite, enterprises, err := uc.GetUserDetails(id.String())
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.NotNil(t, favorite)
		assert.NotNil(t, enterprises)
		mockUserRepository.AssertExpectations(t)
	})

	t.Run("user null", func(t *testing.T) {
		id := uuid.FromStringOrNil("0cf712fc-e631-40c7-8572-54772e698edg")
		uc := usecase.NewAuthUsecase(mockUserRepository, mockRoleRepository, mockFavoriteRepository, mockEnterpriseRepository)
		mockUserRepository.On("FindUserById", mock.AnythingOfType("string")).Return(domain.User{}, errors.New("user not found")).Once()
		_, _, _, err := uc.GetUserDetails(id.String())
		assert.Error(t, err)
		mockUserRepository.AssertExpectations(t)
	})
}

func TestAuthUsecase_CheckIfUserIsAdmin(t *testing.T) {
	mockUserRepository := new(mocks.UserRepository)
	mockRoleRepository := new(mocks.RoleRepository)
	mockFavoriteRepository := new(mocks.FavoriteRepository)
	mockEnterpriseRepository := new(mocks.EnterpriseRepository)

	t.Run("success", func(t *testing.T) {
		id := uuid.FromStringOrNil("0cf712fc-e631-40c7-8572-54772e698edg")
		uc := usecase.NewAuthUsecase(mockUserRepository, mockRoleRepository, mockFavoriteRepository, mockEnterpriseRepository)
		mockUserRepository.On("FindUserById", mock.AnythingOfType("string")).Return(dummyUser[0], nil).Once()
		res, err := uc.CheckIfUserIsAdmin(id.String())
		assert.NoError(t, err)
		assert.NotNil(t, res)
		mockUserRepository.AssertExpectations(t)
	})

	t.Run("user null", func(t *testing.T) {
		id := uuid.FromStringOrNil("0cf712fc-e631-40c7-8572-54772e698edg")
		uc := usecase.NewAuthUsecase(mockUserRepository, mockRoleRepository, mockFavoriteRepository, mockEnterpriseRepository)
		mockUserRepository.On("FindUserById", mock.AnythingOfType("string")).Return(domain.User{}, errors.New("error not found")).Once()
		_, err := uc.CheckIfUserIsAdmin(id.String())
		assert.Error(t, err)
		mockUserRepository.AssertExpectations(t)
	})

	t.Run("role not admin", func(t *testing.T) {
		id := uuid.FromStringOrNil("0cf712fc-e631-40c7-8572-54772e698edf")
		uc := usecase.NewAuthUsecase(mockUserRepository, mockRoleRepository, mockFavoriteRepository, mockEnterpriseRepository)
		mockUserRepository.On("FindUserById", mock.AnythingOfType("string")).Return(dummyUser[1], nil).Once()
		res, err := uc.CheckIfUserIsAdmin(id.String())
		assert.NoError(t, err)
		assert.NotNil(t, res)
		mockUserRepository.AssertExpectations(t)
	})
}
