package usecase_test

import (
	"errors"
	"github.com/nrmadi02/mini-project/domain"
	"github.com/nrmadi02/mini-project/domain/mocks"
	"github.com/nrmadi02/mini-project/internal/review/usecase"
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

var dummyReview = domain.Reviews{
	domain.Review{
		ID:           uuid.FromStringOrNil("35d6a9a1-aa5e-41f1-9991-08878dfdf888"),
		Review:       "baguss",
		EnterpriseID: uuid.FromStringOrNil("35d6a9a1-aa5e-41f1-9991-08878dfdf89a"),
		UserID:       uuid.FromStringOrNil("35d6a9a1-aa5e-41f1-9991-08878dfdf891"),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	},
	domain.Review{
		ID:           uuid.FromStringOrNil("35d6a9a1-aa5e-41f1-9991-08878dfdf999"),
		Review:       "baguss",
		EnterpriseID: uuid.FromStringOrNil("35d6a9a1-aa5e-41f1-9991-08878dfdf89a"),
		UserID:       uuid.FromStringOrNil("35d6a9a1-aa5e-41f1-9991-08878dfdf892"),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	},
}

func TestReviewUsecase_AddReview(t *testing.T) {
	mockEnterpriseRepository := new(mocks.EnterpriseRepository)
	mockReviewRepository := new(mocks.ReviewRepository)
	mockUserRepository := new(mocks.UserRepository)
	mockAuthUsecase := new(mocks.AuthUsecase)
	t.Run("success", func(t *testing.T) {
		uc := usecase.NewReviewUsecase(mockEnterpriseRepository, mockUserRepository, mockReviewRepository, mockAuthUsecase)
		mockEnterpriseRepository.On("FindByID", mock.AnythingOfType("string")).Return(dummyEnterprise[0], nil).Once()
		mockUserRepository.On("FindUserById", mock.AnythingOfType("string")).Return(dummyUser[0], nil).Once()
		mockReviewRepository.On("Add", mock.AnythingOfType("domain.Review")).Return(dummyReview[0], nil).Once()
		review, err := uc.AddReview(dummyEnterprise[0].ID.String(), dummyUser[0].ID.String(), "baguss")
		assert.NoError(t, err)
		assert.NotNil(t, review)
	})
	t.Run("enterprise not found", func(t *testing.T) {
		uc := usecase.NewReviewUsecase(mockEnterpriseRepository, mockUserRepository, mockReviewRepository, mockAuthUsecase)
		mockEnterpriseRepository.On("FindByID", mock.AnythingOfType("string")).Return(domain.Enterprise{}, errors.New("error something")).Once()
		_, err := uc.AddReview(dummyEnterprise[0].ID.String(), dummyUser[0].ID.String(), "baguss")
		assert.Error(t, err)
	})
	t.Run("user not found", func(t *testing.T) {
		uc := usecase.NewReviewUsecase(mockEnterpriseRepository, mockUserRepository, mockReviewRepository, mockAuthUsecase)
		mockEnterpriseRepository.On("FindByID", mock.AnythingOfType("string")).Return(dummyEnterprise[0], nil).Once()
		mockUserRepository.On("FindUserById", mock.AnythingOfType("string")).Return(domain.User{}, errors.New("error something")).Once()
		_, err := uc.AddReview(dummyEnterprise[0].ID.String(), dummyUser[0].ID.String(), "baguss")
		assert.Error(t, err)
	})
	t.Run("failed", func(t *testing.T) {
		uc := usecase.NewReviewUsecase(mockEnterpriseRepository, mockUserRepository, mockReviewRepository, mockAuthUsecase)
		mockEnterpriseRepository.On("FindByID", mock.AnythingOfType("string")).Return(dummyEnterprise[0], nil).Once()
		mockUserRepository.On("FindUserById", mock.AnythingOfType("string")).Return(dummyUser[0], nil).Once()
		mockReviewRepository.On("Add", mock.AnythingOfType("domain.Review")).Return(domain.Review{}, errors.New("error something")).Once()
		_, err := uc.AddReview(dummyEnterprise[0].ID.String(), dummyUser[0].ID.String(), "baguss")
		assert.Error(t, err)
	})
}

func TestReviewUsecase_UpdateReview(t *testing.T) {
	mockEnterpriseRepository := new(mocks.EnterpriseRepository)
	mockReviewRepository := new(mocks.ReviewRepository)
	mockUserRepository := new(mocks.UserRepository)
	mockAuthUsecase := new(mocks.AuthUsecase)
	t.Run("success", func(t *testing.T) {
		uc := usecase.NewReviewUsecase(mockEnterpriseRepository, mockUserRepository, mockReviewRepository, mockAuthUsecase)
		mockReviewRepository.On("FindByUserIDAndEnterpriseID", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(dummyReview[0], nil).Once()
		mockReviewRepository.On("Update", mock.AnythingOfType("string"),
			mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(dummyReview[0], nil).Once()
		review, err := uc.UpdateReview(dummyEnterprise[0].ID.String(), dummyUser[0].ID.String(), "baguss")
		assert.NoError(t, err)
		assert.NotNil(t, review)
	})
	t.Run("request user and enterprise not found", func(t *testing.T) {
		uc := usecase.NewReviewUsecase(mockEnterpriseRepository, mockUserRepository, mockReviewRepository, mockAuthUsecase)
		mockReviewRepository.On("FindByUserIDAndEnterpriseID", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(domain.Review{}, errors.New("error something")).Once()
		_, err := uc.UpdateReview(dummyEnterprise[0].ID.String(), dummyUser[0].ID.String(), "baguss")
		assert.Error(t, err)
	})
	t.Run("failed", func(t *testing.T) {
		uc := usecase.NewReviewUsecase(mockEnterpriseRepository, mockUserRepository, mockReviewRepository, mockAuthUsecase)
		mockReviewRepository.On("FindByUserIDAndEnterpriseID", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(dummyReview[0], nil).Once()
		mockReviewRepository.On("Update", mock.AnythingOfType("string"),
			mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(domain.Review{}, errors.New("error something")).Once()
		_, err := uc.UpdateReview(dummyEnterprise[0].ID.String(), dummyUser[0].ID.String(), "baguss")
		assert.Error(t, err)
	})
}

func TestReviewUsecase_DeleteReview(t *testing.T) {
	mockEnterpriseRepository := new(mocks.EnterpriseRepository)
	mockReviewRepository := new(mocks.ReviewRepository)
	mockUserRepository := new(mocks.UserRepository)
	mockAuthUsecase := new(mocks.AuthUsecase)
	t.Run("success", func(t *testing.T) {
		uc := usecase.NewReviewUsecase(mockEnterpriseRepository, mockUserRepository, mockReviewRepository, mockAuthUsecase)
		mockReviewRepository.On("FindByUserIDAndEnterpriseID", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(dummyReview[0], nil).Once()
		mockReviewRepository.On("Delete", mock.AnythingOfType("domain.Review")).Return(nil).Once()
		err := uc.DeleteReview(dummyEnterprise[0].ID.String(), dummyUser[0].ID.String())
		assert.NoError(t, err)
	})
	t.Run("request user and enterprise not found", func(t *testing.T) {
		uc := usecase.NewReviewUsecase(mockEnterpriseRepository, mockUserRepository, mockReviewRepository, mockAuthUsecase)
		mockReviewRepository.On("FindByUserIDAndEnterpriseID", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(domain.Review{}, errors.New("error something")).Once()
		err := uc.DeleteReview(dummyEnterprise[0].ID.String(), dummyUser[0].ID.String())
		assert.Error(t, err)
	})
}

func TestReviewUsecase_GetListReviewsByEnterpriseID(t *testing.T) {
	mockEnterpriseRepository := new(mocks.EnterpriseRepository)
	mockReviewRepository := new(mocks.ReviewRepository)
	mockUserRepository := new(mocks.UserRepository)
	mockAuthUsecase := new(mocks.AuthUsecase)
	t.Run("success", func(t *testing.T) {
		uc := usecase.NewReviewUsecase(mockEnterpriseRepository, mockUserRepository, mockReviewRepository, mockAuthUsecase)
		mockReviewRepository.On("FindByEnterpriseID", mock.AnythingOfType("string")).Return(domain.Reviews{dummyReview[0]}, nil).Once()
		mockAuthUsecase.On("GetUserDetails", mock.Anything).Return(dummyUser[0], domain.Favorite{}, domain.Enterprises{}, nil).Once()
		reviews, err := uc.GetListReviewsByEnterpriseID(dummyEnterprise[0].ID.String())
		assert.NoError(t, err)
		assert.NotNil(t, reviews)
	})
	t.Run("failed", func(t *testing.T) {
		uc := usecase.NewReviewUsecase(mockEnterpriseRepository, mockUserRepository, mockReviewRepository, mockAuthUsecase)
		mockReviewRepository.On("FindByEnterpriseID", mock.AnythingOfType("string")).Return(domain.Reviews{}, errors.New("error something")).Once()
		_, err := uc.GetListReviewsByEnterpriseID(dummyEnterprise[0].ID.String())
		assert.Error(t, err)
	})
}

func TestReviewUsecase_GetDetailReviewByID(t *testing.T) {
	mockEnterpriseRepository := new(mocks.EnterpriseRepository)
	mockReviewRepository := new(mocks.ReviewRepository)
	mockUserRepository := new(mocks.UserRepository)
	mockAuthUsecase := new(mocks.AuthUsecase)
	t.Run("success", func(t *testing.T) {
		uc := usecase.NewReviewUsecase(mockEnterpriseRepository, mockUserRepository, mockReviewRepository, mockAuthUsecase)
		mockReviewRepository.On("FindByID", mock.AnythingOfType("string")).Return(dummyReview[0], nil).Once()
		review, err := uc.GetDetailReviewByID(dummyReview[0].ID.String())
		assert.NoError(t, err)
		assert.NotNil(t, review)
	})
	t.Run("failed", func(t *testing.T) {
		uc := usecase.NewReviewUsecase(mockEnterpriseRepository, mockUserRepository, mockReviewRepository, mockAuthUsecase)
		mockReviewRepository.On("FindByID", mock.AnythingOfType("string")).Return(domain.Review{}, errors.New("error something")).Once()
		_, err := uc.GetDetailReviewByID(dummyReview[0].ID.String())
		assert.Error(t, err)
	})
}

func TestReviewUsecase_GetReviewByUserIDAndEnterpriseID(t *testing.T) {
	mockEnterpriseRepository := new(mocks.EnterpriseRepository)
	mockReviewRepository := new(mocks.ReviewRepository)
	mockUserRepository := new(mocks.UserRepository)
	mockAuthUsecase := new(mocks.AuthUsecase)
	t.Run("success", func(t *testing.T) {
		uc := usecase.NewReviewUsecase(mockEnterpriseRepository, mockUserRepository, mockReviewRepository, mockAuthUsecase)
		mockReviewRepository.On("FindByUserIDAndEnterpriseID", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(dummyReview[0], nil).Once()
		review, err := uc.GetReviewByUserIDAndEnterpriseID(dummyEnterprise[0].ID.String(), dummyUser[0].ID.String())
		assert.NoError(t, err)
		assert.NotNil(t, review)
	})
	t.Run("failed", func(t *testing.T) {
		uc := usecase.NewReviewUsecase(mockEnterpriseRepository, mockUserRepository, mockReviewRepository, mockAuthUsecase)
		mockReviewRepository.On("FindByUserIDAndEnterpriseID", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(domain.Review{}, errors.New("error something")).Once()
		_, err := uc.GetReviewByUserIDAndEnterpriseID(dummyEnterprise[0].ID.String(), dummyUser[0].ID.String())
		assert.Error(t, err)
	})

}
