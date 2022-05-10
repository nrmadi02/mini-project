package usecase_test

import (
	"errors"
	"github.com/nrmadi02/mini-project/domain"
	"github.com/nrmadi02/mini-project/domain/mocks"
	"github.com/nrmadi02/mini-project/internal/rating/usecase"
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

var dummyRating = domain.RatingEnterprises{
	domain.RatingEnterprise{
		ID:           uuid.FromStringOrNil("35d6a9a1-aa5e-41f1-9991-08878dfdf777"),
		Rating:       3,
		EnterpriseID: uuid.FromStringOrNil("35d6a9a1-aa5e-41f1-9991-08878dfdf89a"),
		UserID:       uuid.FromStringOrNil("0cf712fc-e631-40c7-8572-54772e698edf"),
	},
	domain.RatingEnterprise{
		ID:           uuid.FromStringOrNil("35d6a9a1-aa5e-41f1-9991-08878dfdf878"),
		Rating:       3,
		EnterpriseID: uuid.FromStringOrNil("35d6a9a1-aa5e-41f1-9991-08878dfdf89a"),
		UserID:       uuid.FromStringOrNil("0cf712fc-e631-40c7-8572-54772e698787"),
	},
}

func TestRatingUsecase_AddNewRanting(t *testing.T) {
	mockEnterpriseRepository := new(mocks.EnterpriseRepository)
	mockRatingRepository := new(mocks.RatingRepository)
	mockUserRepository := new(mocks.UserRepository)
	t.Run("success", func(t *testing.T) {
		uc := usecase.NewRatingUsecase(mockUserRepository, mockEnterpriseRepository, mockRatingRepository)
		mockUserRepository.On("FindUserById", mock.AnythingOfType("string")).Return(dummyUser[0], nil).Once()
		mockEnterpriseRepository.On("FindByID", mock.AnythingOfType("string")).Return(dummyEnterprise[0], nil).Once()
		mockRatingRepository.On("AddRating", mock.AnythingOfType("domain.RatingEnterprise")).Return(dummyRating[0], nil).Once()
		ranting, err := uc.AddNewRanting(dummyEnterprise[0].ID.String(), dummyUser[0].ID.String(), 3)
		assert.NoError(t, err)
		assert.NotNil(t, ranting)
	})
	t.Run("user not found", func(t *testing.T) {
		uc := usecase.NewRatingUsecase(mockUserRepository, mockEnterpriseRepository, mockRatingRepository)
		mockUserRepository.On("FindUserById", mock.AnythingOfType("string")).Return(domain.User{}, errors.New("error something")).Once()
		_, err := uc.AddNewRanting(dummyEnterprise[0].ID.String(), dummyUser[0].ID.String(), 3)
		assert.Error(t, err)
	})
	t.Run("enterprise not found", func(t *testing.T) {
		uc := usecase.NewRatingUsecase(mockUserRepository, mockEnterpriseRepository, mockRatingRepository)
		mockUserRepository.On("FindUserById", mock.AnythingOfType("string")).Return(dummyUser[0], nil).Once()
		mockEnterpriseRepository.On("FindByID", mock.AnythingOfType("string")).Return(domain.Enterprise{}, errors.New("error something")).Once()
		_, err := uc.AddNewRanting(dummyEnterprise[0].ID.String(), dummyUser[0].ID.String(), 3)
		assert.Error(t, err)
	})
	t.Run("failed", func(t *testing.T) {
		uc := usecase.NewRatingUsecase(mockUserRepository, mockEnterpriseRepository, mockRatingRepository)
		mockUserRepository.On("FindUserById", mock.AnythingOfType("string")).Return(dummyUser[0], nil).Once()
		mockEnterpriseRepository.On("FindByID", mock.AnythingOfType("string")).Return(dummyEnterprise[0], nil).Once()
		mockRatingRepository.On("AddRating", mock.AnythingOfType("domain.RatingEnterprise")).Return(domain.RatingEnterprise{}, errors.New("error something")).Once()
		_, err := uc.AddNewRanting(dummyEnterprise[0].ID.String(), dummyUser[0].ID.String(), 3)
		assert.Error(t, err)
	})
}

func TestRatingUsecase_FindRating(t *testing.T) {
	mockEnterpriseRepository := new(mocks.EnterpriseRepository)
	mockRatingRepository := new(mocks.RatingRepository)
	mockUserRepository := new(mocks.UserRepository)
	t.Run("success", func(t *testing.T) {
		uc := usecase.NewRatingUsecase(mockUserRepository, mockEnterpriseRepository, mockRatingRepository)
		mockUserRepository.On("FindUserById", mock.AnythingOfType("string")).Return(dummyUser[0], nil).Once()
		mockEnterpriseRepository.On("FindByID", mock.AnythingOfType("string")).Return(dummyEnterprise[0], nil).Once()
		mockRatingRepository.On("FindRatingByIDUserAndEnterprise", mock.AnythingOfType("string"),
			mock.AnythingOfType("string")).Return(dummyRating[0], nil).Once()
		ranting, err := uc.FindRating(dummyEnterprise[0].ID.String(), dummyUser[0].ID.String())
		assert.NoError(t, err)
		assert.NotNil(t, ranting)
	})
	t.Run("user not found", func(t *testing.T) {
		uc := usecase.NewRatingUsecase(mockUserRepository, mockEnterpriseRepository, mockRatingRepository)
		mockUserRepository.On("FindUserById", mock.AnythingOfType("string")).Return(domain.User{}, errors.New("error something")).Once()
		_, err := uc.FindRating(dummyEnterprise[0].ID.String(), dummyUser[0].ID.String())
		assert.Error(t, err)
	})
	t.Run("enterprise not found", func(t *testing.T) {
		uc := usecase.NewRatingUsecase(mockUserRepository, mockEnterpriseRepository, mockRatingRepository)
		mockUserRepository.On("FindUserById", mock.AnythingOfType("string")).Return(dummyUser[0], nil).Once()
		mockEnterpriseRepository.On("FindByID", mock.AnythingOfType("string")).Return(domain.Enterprise{}, errors.New("error something")).Once()
		_, err := uc.FindRating(dummyEnterprise[0].ID.String(), dummyUser[0].ID.String())
		assert.Error(t, err)
	})
	t.Run("failed", func(t *testing.T) {
		uc := usecase.NewRatingUsecase(mockUserRepository, mockEnterpriseRepository, mockRatingRepository)
		mockUserRepository.On("FindUserById", mock.AnythingOfType("string")).Return(dummyUser[0], nil).Once()
		mockEnterpriseRepository.On("FindByID", mock.AnythingOfType("string")).Return(dummyEnterprise[0], nil).Once()
		mockRatingRepository.On("FindRatingByIDUserAndEnterprise", mock.AnythingOfType("string"),
			mock.AnythingOfType("string")).Return(domain.RatingEnterprise{}, errors.New("error something")).Once()
		_, err := uc.FindRating(dummyEnterprise[0].ID.String(), dummyUser[0].ID.String())
		assert.Error(t, err)
	})
}

func TestRatingUsecase_GetAllRatingByEnterpriseID(t *testing.T) {
	mockEnterpriseRepository := new(mocks.EnterpriseRepository)
	mockRatingRepository := new(mocks.RatingRepository)
	mockUserRepository := new(mocks.UserRepository)
	t.Run("success", func(t *testing.T) {
		uc := usecase.NewRatingUsecase(mockUserRepository, mockEnterpriseRepository, mockRatingRepository)
		mockEnterpriseRepository.On("FindByID", mock.AnythingOfType("string")).Return(dummyEnterprise[0], nil).Once()
		mockRatingRepository.On("GetAllRatingByEnterpriseID", mock.AnythingOfType("string"),
			mock.AnythingOfType("string")).Return(domain.RatingEnterprises{
			dummyRating[0],
		}, nil).Once()
		rantings, err := uc.GetAllRatingByEnterpriseID(dummyEnterprise[0].ID.String())
		assert.NoError(t, err)
		assert.NotNil(t, rantings)
	})
	t.Run("enterprise not found", func(t *testing.T) {
		uc := usecase.NewRatingUsecase(mockUserRepository, mockEnterpriseRepository, mockRatingRepository)
		mockEnterpriseRepository.On("FindByID", mock.AnythingOfType("string")).Return(domain.Enterprise{}, errors.New("error something")).Once()
		_, err := uc.GetAllRatingByEnterpriseID(dummyEnterprise[0].ID.String())
		assert.Error(t, err)
	})
	t.Run("failed", func(t *testing.T) {
		uc := usecase.NewRatingUsecase(mockUserRepository, mockEnterpriseRepository, mockRatingRepository)
		mockEnterpriseRepository.On("FindByID", mock.AnythingOfType("string")).Return(dummyEnterprise[0], nil).Once()
		mockRatingRepository.On("GetAllRatingByEnterpriseID", mock.AnythingOfType("string"),
			mock.AnythingOfType("string")).Return(domain.RatingEnterprises{}, errors.New("error something")).Once()
		_, err := uc.GetAllRatingByEnterpriseID(dummyEnterprise[0].ID.String())
		assert.Error(t, err)
	})
}

func TestRatingUsecase_UpdateRating(t *testing.T) {
	mockEnterpriseRepository := new(mocks.EnterpriseRepository)
	mockRatingRepository := new(mocks.RatingRepository)
	mockUserRepository := new(mocks.UserRepository)
	t.Run("success", func(t *testing.T) {
		uc := usecase.NewRatingUsecase(mockUserRepository, mockEnterpriseRepository, mockRatingRepository)
		mockUserRepository.On("FindUserById", mock.AnythingOfType("string")).Return(dummyUser[0], nil).Once()
		mockEnterpriseRepository.On("FindByID", mock.AnythingOfType("string")).Return(dummyEnterprise[0], nil).Once()
		mockRatingRepository.On("UpdateRating", mock.AnythingOfType("string"),
			mock.AnythingOfType("string"), mock.AnythingOfType("int")).Return(dummyRating[0], nil).Once()
		ranting, err := uc.UpdateRating(dummyEnterprise[0].ID.String(), dummyUser[0].ID.String(), 3)
		assert.NoError(t, err)
		assert.NotNil(t, ranting)
	})
	t.Run("user not found", func(t *testing.T) {
		uc := usecase.NewRatingUsecase(mockUserRepository, mockEnterpriseRepository, mockRatingRepository)
		mockUserRepository.On("FindUserById", mock.AnythingOfType("string")).Return(domain.User{}, errors.New("error something")).Once()
		_, err := uc.UpdateRating(dummyEnterprise[0].ID.String(), dummyUser[0].ID.String(), 3)
		assert.Error(t, err)
	})
	t.Run("enterprise not found", func(t *testing.T) {
		uc := usecase.NewRatingUsecase(mockUserRepository, mockEnterpriseRepository, mockRatingRepository)
		mockUserRepository.On("FindUserById", mock.AnythingOfType("string")).Return(dummyUser[0], nil).Once()
		mockEnterpriseRepository.On("FindByID", mock.AnythingOfType("string")).Return(domain.Enterprise{}, errors.New("error something")).Once()
		_, err := uc.UpdateRating(dummyEnterprise[0].ID.String(), dummyUser[0].ID.String(), 3)
		assert.Error(t, err)
	})
	t.Run("failed", func(t *testing.T) {
		uc := usecase.NewRatingUsecase(mockUserRepository, mockEnterpriseRepository, mockRatingRepository)
		mockUserRepository.On("FindUserById", mock.AnythingOfType("string")).Return(dummyUser[0], nil).Once()
		mockEnterpriseRepository.On("FindByID", mock.AnythingOfType("string")).Return(dummyEnterprise[0], nil).Once()
		mockRatingRepository.On("UpdateRating", mock.AnythingOfType("string"),
			mock.AnythingOfType("string"), mock.AnythingOfType("int")).Return(domain.RatingEnterprise{}, errors.New("error something")).Once()
		_, err := uc.UpdateRating(dummyEnterprise[0].ID.String(), dummyUser[0].ID.String(), 3)
		assert.Error(t, err)
	})
}

func TestRatingUsecase_DeleteRating(t *testing.T) {
	mockEnterpriseRepository := new(mocks.EnterpriseRepository)
	mockRatingRepository := new(mocks.RatingRepository)
	mockUserRepository := new(mocks.UserRepository)
	t.Run("success", func(t *testing.T) {
		uc := usecase.NewRatingUsecase(mockUserRepository, mockEnterpriseRepository, mockRatingRepository)
		mockUserRepository.On("FindUserById", mock.AnythingOfType("string")).Return(dummyUser[0], nil).Once()
		mockEnterpriseRepository.On("FindByID", mock.AnythingOfType("string")).Return(dummyEnterprise[0], nil).Once()
		mockRatingRepository.On("FindRatingByIDUserAndEnterprise", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(dummyRating[0], nil).Once()
		mockRatingRepository.On("DeleteRating", mock.AnythingOfType("domain.RatingEnterprise")).Return(nil).Once()
		err := uc.DeleteRating(dummyEnterprise[0].ID.String(), dummyUser[0].ID.String())
		assert.NoError(t, err)
	})
	t.Run("user not found", func(t *testing.T) {
		uc := usecase.NewRatingUsecase(mockUserRepository, mockEnterpriseRepository, mockRatingRepository)
		mockUserRepository.On("FindUserById", mock.AnythingOfType("string")).Return(domain.User{}, errors.New("error something")).Once()
		err := uc.DeleteRating(dummyEnterprise[0].ID.String(), dummyUser[0].ID.String())
		assert.Error(t, err)
	})
	t.Run("enterprise not found", func(t *testing.T) {
		uc := usecase.NewRatingUsecase(mockUserRepository, mockEnterpriseRepository, mockRatingRepository)
		mockUserRepository.On("FindUserById", mock.AnythingOfType("string")).Return(dummyUser[0], nil).Once()
		mockEnterpriseRepository.On("FindByID", mock.AnythingOfType("string")).Return(domain.Enterprise{}, errors.New("error something")).Once()
		err := uc.DeleteRating(dummyEnterprise[0].ID.String(), dummyUser[0].ID.String())
		assert.Error(t, err)
	})
	t.Run("rating not found", func(t *testing.T) {
		uc := usecase.NewRatingUsecase(mockUserRepository, mockEnterpriseRepository, mockRatingRepository)
		mockUserRepository.On("FindUserById", mock.AnythingOfType("string")).Return(dummyUser[0], nil).Once()
		mockEnterpriseRepository.On("FindByID", mock.AnythingOfType("string")).Return(dummyEnterprise[0], nil).Once()
		mockRatingRepository.On("FindRatingByIDUserAndEnterprise", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(domain.RatingEnterprise{}, errors.New("error something")).Once()
		err := uc.DeleteRating(dummyEnterprise[0].ID.String(), dummyUser[0].ID.String())
		assert.Error(t, err)
	})
	t.Run("failed", func(t *testing.T) {
		uc := usecase.NewRatingUsecase(mockUserRepository, mockEnterpriseRepository, mockRatingRepository)
		mockUserRepository.On("FindUserById", mock.AnythingOfType("string")).Return(dummyUser[0], nil).Once()
		mockEnterpriseRepository.On("FindByID", mock.AnythingOfType("string")).Return(dummyEnterprise[0], nil).Once()
		mockRatingRepository.On("FindRatingByIDUserAndEnterprise", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(dummyRating[0], nil).Once()
		mockRatingRepository.On("DeleteRating", mock.AnythingOfType("domain.RatingEnterprise")).Return(errors.New("error something")).Once()
		err := uc.DeleteRating(dummyEnterprise[0].ID.String(), dummyUser[0].ID.String())
		assert.Error(t, err)
	})
}
