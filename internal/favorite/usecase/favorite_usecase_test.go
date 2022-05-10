package usecase_test

import (
	"errors"
	"github.com/nrmadi02/mini-project/domain"
	"github.com/nrmadi02/mini-project/domain/mocks"
	"github.com/nrmadi02/mini-project/internal/favorite/usecase"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
	"testing"
	"time"
)

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
		ID:     uuid.FromStringOrNil("35d6a9a1-aa5e-41f1-9991-08878dfdf888"),
		UserID: uuid.FromStringOrNil("35d6a9a1-aa5e-41f1-9991-08878dfdf891"),
		Enterprises: domain.Enterprises{
			dummyEnterprise[0],
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	},
	domain.Favorite{
		ID:     uuid.FromStringOrNil("35d6a9a1-aa5e-41f1-9991-08878dfdf877"),
		UserID: uuid.FromStringOrNil("35d6a9a1-aa5e-41f1-9991-08878dfdf891"),
		Enterprises: domain.Enterprises{
			dummyEnterprise[0],
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	},
}

func TestFavoriteUsecase_AddFavorite(t *testing.T) {
	mockEnterpriseRepository := new(mocks.EnterpriseRepository)
	mockFavoriteRepository := new(mocks.FavoriteRepository)

	t.Run("success", func(t *testing.T) {
		uc := usecase.NewFavoriteUsecase(mockEnterpriseRepository, mockFavoriteRepository)
		mockEnterpriseRepository.On("FindByIDs", mock.AnythingOfType("[]string")).Return(domain.Enterprises{
			dummyEnterprise[0],
		}, nil).Once()
		mockFavoriteRepository.On("FindByUserID", mock.AnythingOfType("string")).Return(dummyFavorite[0], nil).Once()
		mockFavoriteRepository.On("Update",
			mock.AnythingOfType("domain.Favorite"), mock.AnythingOfType("domain.Enterprises"), mock.AnythingOfType("string")).
			Return(dummyFavorite[0], nil).Once()
		favorite, err := uc.AddFavorite([]string{dummyEnterprise[0].ID.String()}, dummyUser[0].ID.String())
		assert.NoError(t, err)
		assert.NotNil(t, favorite)
		mockFavoriteRepository.AssertExpectations(t)
	})

	t.Run("list enterprises not found", func(t *testing.T) {
		uc := usecase.NewFavoriteUsecase(mockEnterpriseRepository, mockFavoriteRepository)
		mockEnterpriseRepository.On("FindByIDs", mock.AnythingOfType("[]string")).Return(domain.Enterprises{}, nil).Once()
		_, err := uc.AddFavorite([]string{dummyEnterprise[0].ID.String()}, dummyUser[0].ID.String())
		assert.Error(t, err)
		mockFavoriteRepository.AssertExpectations(t)
	})

	t.Run("user not register favorite", func(t *testing.T) {
		uc := usecase.NewFavoriteUsecase(mockEnterpriseRepository, mockFavoriteRepository)
		mockEnterpriseRepository.On("FindByIDs", mock.AnythingOfType("[]string")).Return(domain.Enterprises{
			dummyEnterprise[0],
		}, nil).Once()
		mockFavoriteRepository.On("FindByUserID", mock.AnythingOfType("string")).Return(domain.Favorite{}, errors.New("error something")).Once()
		_, err := uc.AddFavorite([]string{dummyEnterprise[0].ID.String()}, dummyUser[0].ID.String())
		assert.Error(t, err)
		mockFavoriteRepository.AssertExpectations(t)
	})

	t.Run("failed add", func(t *testing.T) {
		uc := usecase.NewFavoriteUsecase(mockEnterpriseRepository, mockFavoriteRepository)
		mockEnterpriseRepository.On("FindByIDs", mock.AnythingOfType("[]string")).Return(domain.Enterprises{
			dummyEnterprise[0],
		}, nil).Once()
		mockFavoriteRepository.On("FindByUserID", mock.AnythingOfType("string")).Return(dummyFavorite[0], nil).Once()
		mockFavoriteRepository.On("Update",
			mock.AnythingOfType("domain.Favorite"), mock.AnythingOfType("domain.Enterprises"), mock.AnythingOfType("string")).
			Return(domain.Favorite{}, errors.New("error something")).Once()
		_, err := uc.AddFavorite([]string{dummyEnterprise[0].ID.String()}, dummyUser[0].ID.String())
		assert.Error(t, err)
		mockFavoriteRepository.AssertExpectations(t)
	})
}

func TestFavoriteUsecase_RemoveFavorite(t *testing.T) {
	mockEnterpriseRepository := new(mocks.EnterpriseRepository)
	mockFavoriteRepository := new(mocks.FavoriteRepository)

	t.Run("success", func(t *testing.T) {
		uc := usecase.NewFavoriteUsecase(mockEnterpriseRepository, mockFavoriteRepository)
		mockEnterpriseRepository.On("FindByIDs", mock.AnythingOfType("[]string")).Return(domain.Enterprises{
			dummyEnterprise[0],
		}, nil).Once()
		mockFavoriteRepository.On("FindByUserID", mock.AnythingOfType("string")).Return(dummyFavorite[0], nil).Once()
		mockFavoriteRepository.On("Update",
			mock.AnythingOfType("domain.Favorite"), mock.AnythingOfType("domain.Enterprises"), mock.AnythingOfType("string")).
			Return(dummyFavorite[0], nil).Once()
		favorite, err := uc.RemoveFavorite([]string{dummyEnterprise[0].ID.String()}, dummyUser[0].ID.String())
		assert.NoError(t, err)
		assert.NotNil(t, favorite)
		mockFavoriteRepository.AssertExpectations(t)
	})

	t.Run("list enterprises not found", func(t *testing.T) {
		uc := usecase.NewFavoriteUsecase(mockEnterpriseRepository, mockFavoriteRepository)
		mockEnterpriseRepository.On("FindByIDs", mock.AnythingOfType("[]string")).Return(domain.Enterprises{}, errors.New("error something")).Once()
		_, err := uc.RemoveFavorite([]string{dummyEnterprise[0].ID.String()}, dummyUser[0].ID.String())
		assert.Error(t, err)
		mockFavoriteRepository.AssertExpectations(t)
	})

	t.Run("user not register favorite", func(t *testing.T) {
		uc := usecase.NewFavoriteUsecase(mockEnterpriseRepository, mockFavoriteRepository)
		mockEnterpriseRepository.On("FindByIDs", mock.AnythingOfType("[]string")).Return(domain.Enterprises{
			dummyEnterprise[0],
		}, nil).Once()
		mockFavoriteRepository.On("FindByUserID", mock.AnythingOfType("string")).Return(domain.Favorite{}, errors.New("error something")).Once()
		_, err := uc.RemoveFavorite([]string{dummyEnterprise[0].ID.String()}, dummyUser[0].ID.String())
		assert.Error(t, err)
		mockFavoriteRepository.AssertExpectations(t)
	})

	t.Run("failed remove", func(t *testing.T) {
		uc := usecase.NewFavoriteUsecase(mockEnterpriseRepository, mockFavoriteRepository)
		mockEnterpriseRepository.On("FindByIDs", mock.AnythingOfType("[]string")).Return(domain.Enterprises{
			dummyEnterprise[0],
		}, nil).Once()
		mockFavoriteRepository.On("FindByUserID", mock.AnythingOfType("string")).Return(dummyFavorite[0], nil).Once()
		mockFavoriteRepository.On("Update",
			mock.AnythingOfType("domain.Favorite"), mock.AnythingOfType("domain.Enterprises"), mock.AnythingOfType("string")).
			Return(domain.Favorite{}, errors.New("error something")).Once()
		_, err := uc.RemoveFavorite([]string{dummyEnterprise[0].ID.String()}, dummyUser[0].ID.String())
		assert.Error(t, err)
		mockFavoriteRepository.AssertExpectations(t)
	})
}

func TestFavoriteUsecase_GetDetailByID(t *testing.T) {
	mockEnterpriseRepository := new(mocks.EnterpriseRepository)
	mockFavoriteRepository := new(mocks.FavoriteRepository)
	t.Run("success", func(t *testing.T) {
		uc := usecase.NewFavoriteUsecase(mockEnterpriseRepository, mockFavoriteRepository)
		mockFavoriteRepository.On("FindByID", mock.AnythingOfType("string")).Return(dummyFavorite[0], nil).Once()
		favorite, err := uc.GetDetailByID(dummyFavorite[0].ID.String())
		assert.NoError(t, err)
		assert.NotNil(t, favorite)
		mockFavoriteRepository.AssertExpectations(t)
	})

	t.Run("failed", func(t *testing.T) {
		uc := usecase.NewFavoriteUsecase(mockEnterpriseRepository, mockFavoriteRepository)
		mockFavoriteRepository.On("FindByID", mock.AnythingOfType("string")).Return(domain.Favorite{}, errors.New("not found")).Once()
		_, err := uc.GetDetailByID(dummyFavorite[0].ID.String())
		assert.Error(t, err)
		mockFavoriteRepository.AssertExpectations(t)
	})
}

func TestFavoriteUsecase_GetDetailByUserID(t *testing.T) {
	mockEnterpriseRepository := new(mocks.EnterpriseRepository)
	mockFavoriteRepository := new(mocks.FavoriteRepository)
	t.Run("success", func(t *testing.T) {
		uc := usecase.NewFavoriteUsecase(mockEnterpriseRepository, mockFavoriteRepository)
		mockFavoriteRepository.On("FindByUserID", mock.AnythingOfType("string")).Return(dummyFavorite[0], nil).Once()
		favorite, err := uc.GetDetailByUserID(dummyFavorite[0].UserID.String())
		assert.NoError(t, err)
		assert.NotNil(t, favorite)
		mockFavoriteRepository.AssertExpectations(t)
	})

	t.Run("failed", func(t *testing.T) {
		uc := usecase.NewFavoriteUsecase(mockEnterpriseRepository, mockFavoriteRepository)
		mockFavoriteRepository.On("FindByUserID", mock.AnythingOfType("string")).Return(domain.Favorite{}, errors.New("not found")).Once()
		_, err := uc.GetDetailByUserID(dummyFavorite[0].UserID.String())
		assert.Error(t, err)
		mockFavoriteRepository.AssertExpectations(t)
	})
}
