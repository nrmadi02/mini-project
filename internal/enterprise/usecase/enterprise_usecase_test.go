package usecase_test

import (
	"errors"
	"github.com/nrmadi02/mini-project/domain"
	"github.com/nrmadi02/mini-project/domain/mocks"
	"github.com/nrmadi02/mini-project/internal/enterprise/usecase"
	"github.com/nrmadi02/mini-project/web/request"
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

func TestEnterpriseUsecase_CreateNewEnterprise(t *testing.T) {
	mockEnterpriseRepository := new(mocks.EnterpriseRepository)
	mockTagRepository := new(mocks.TagRepository)
	mockUserRepository := new(mocks.UserRepository)

	t.Run("success", func(t *testing.T) {
		req := request.CreateEnterpriseRequest{
			Name:        "enterprise satu",
			NumberPhone: "0012798232",
			Address:     "bjb",
			Postcode:    707722,
			Latitude:    "4235,23",
			Longitude:   "4225,233231",
			Tags:        []string{uuid.FromStringOrNil("35d6a9a1-aa5e-41f1-9991-08878dfdf89b").String()},
		}
		uc := usecase.NewEnterpriseUsecase(mockEnterpriseRepository, mockTagRepository, mockUserRepository)
		mockTagRepository.On("FindByIDs", mock.AnythingOfType("[]string")).Return(domain.Tags{
			domain.Tag{
				ID:   uuid.FromStringOrNil("35d6a9a1-aa5e-41f1-9991-08878dfdf89b"),
				Name: "Tag satu",
			},
		}, nil).Once()
		mockUserRepository.On("FindUserById", mock.AnythingOfType("string")).Return(dummyUser[0], nil).Once()
		mockEnterpriseRepository.On("Save", mock.AnythingOfType("domain.Enterprise")).Return(dummyEnterprise[0], nil).Once()
		enterprise, err := uc.CreateNewEnterprise(req, dummyEnterprise[0].UserID.String())
		assert.NoError(t, err)
		assert.NotNil(t, enterprise)
		mockEnterpriseRepository.AssertExpectations(t)
	})

	t.Run("error get list tags", func(t *testing.T) {
		req := request.CreateEnterpriseRequest{
			Name:        "enterprise satu",
			NumberPhone: "0012798232",
			Address:     "bjb",
			Postcode:    707722,
			Latitude:    "4235,23",
			Longitude:   "4225,233231",
			Tags:        []string{uuid.FromStringOrNil("35d6a9a1-aa5e-41f1-9991-08878dfdf89b").String()},
		}
		uc := usecase.NewEnterpriseUsecase(mockEnterpriseRepository, mockTagRepository, mockUserRepository)
		mockTagRepository.On("FindByIDs", mock.AnythingOfType("[]string")).Return(domain.Tags{}, errors.New("tag not found")).Once()
		_, err := uc.CreateNewEnterprise(req, dummyEnterprise[0].UserID.String())
		assert.Error(t, err)
		mockEnterpriseRepository.AssertExpectations(t)
	})

	t.Run("error user not found", func(t *testing.T) {
		req := request.CreateEnterpriseRequest{
			Name:        "enterprise satu",
			NumberPhone: "0012798232",
			Address:     "bjb",
			Postcode:    707722,
			Latitude:    "4235,23",
			Longitude:   "4225,233231",
			Tags:        []string{uuid.FromStringOrNil("35d6a9a1-aa5e-41f1-9991-08878dfdf89b").String()},
		}
		uc := usecase.NewEnterpriseUsecase(mockEnterpriseRepository, mockTagRepository, mockUserRepository)
		mockTagRepository.On("FindByIDs", mock.AnythingOfType("[]string")).Return(domain.Tags{
			domain.Tag{
				ID:   uuid.FromStringOrNil("35d6a9a1-aa5e-41f1-9991-08878dfdf89b"),
				Name: "Tag satu",
			},
		}, nil).Once()
		mockUserRepository.On("FindUserById", mock.AnythingOfType("string")).Return(domain.User{}, errors.New("user not found")).Once()
		_, err := uc.CreateNewEnterprise(req, dummyEnterprise[0].UserID.String())
		assert.Error(t, err)
		mockEnterpriseRepository.AssertExpectations(t)
	})

	t.Run("failed to save", func(t *testing.T) {
		req := request.CreateEnterpriseRequest{
			Name:        "enterprise satu",
			NumberPhone: "0012798232",
			Address:     "bjb",
			Postcode:    707722,
			Latitude:    "4235,23",
			Longitude:   "4225,233231",
			Tags:        []string{uuid.FromStringOrNil("35d6a9a1-aa5e-41f1-9991-08878dfdf89b").String()},
		}
		uc := usecase.NewEnterpriseUsecase(mockEnterpriseRepository, mockTagRepository, mockUserRepository)
		mockTagRepository.On("FindByIDs", mock.AnythingOfType("[]string")).Return(domain.Tags{
			domain.Tag{
				ID:   uuid.FromStringOrNil("35d6a9a1-aa5e-41f1-9991-08878dfdf89b"),
				Name: "Tag satu",
			},
		}, nil).Once()
		mockUserRepository.On("FindUserById", mock.AnythingOfType("string")).Return(dummyUser[0], nil).Once()
		mockEnterpriseRepository.On("Save", mock.AnythingOfType("domain.Enterprise")).Return(domain.Enterprise{}, errors.New("failed to save")).Once()
		_, err := uc.CreateNewEnterprise(req, dummyEnterprise[0].UserID.String())
		assert.Error(t, err)
		mockEnterpriseRepository.AssertExpectations(t)
	})

}

func TestEnterpriseUsecase_DeleteEnterpriseByID(t *testing.T) {
	mockEnterpriseRepository := new(mocks.EnterpriseRepository)
	mockTagRepository := new(mocks.TagRepository)
	mockUserRepository := new(mocks.UserRepository)

	t.Run("success", func(t *testing.T) {
		uc := usecase.NewEnterpriseUsecase(mockEnterpriseRepository, mockTagRepository, mockUserRepository)
		mockEnterpriseRepository.On("FindByID", mock.AnythingOfType("string")).Return(dummyEnterprise[0], nil).Once()
		mockEnterpriseRepository.On("Delete", mock.AnythingOfType("domain.Enterprise")).Return(nil).Once()
		err := uc.DeleteEnterpriseByID(dummyEnterprise[0].UserID.String())
		assert.NoError(t, err)
		mockEnterpriseRepository.AssertExpectations(t)
	})

	t.Run("enterprise not found", func(t *testing.T) {
		uc := usecase.NewEnterpriseUsecase(mockEnterpriseRepository, mockTagRepository, mockUserRepository)
		mockEnterpriseRepository.On("FindByID", mock.AnythingOfType("string")).Return(domain.Enterprise{}, errors.New("enterprise not found")).Once()
		err := uc.DeleteEnterpriseByID(dummyEnterprise[0].UserID.String())
		assert.Error(t, err)
		mockEnterpriseRepository.AssertExpectations(t)
	})

	t.Run("failed delete", func(t *testing.T) {
		uc := usecase.NewEnterpriseUsecase(mockEnterpriseRepository, mockTagRepository, mockUserRepository)
		mockEnterpriseRepository.On("FindByID", mock.AnythingOfType("string")).Return(dummyEnterprise[0], nil).Once()
		mockEnterpriseRepository.On("Delete", mock.AnythingOfType("domain.Enterprise")).Return(errors.New("failed delete")).Once()
		err := uc.DeleteEnterpriseByID(dummyEnterprise[0].UserID.String())
		assert.Error(t, err)
		mockEnterpriseRepository.AssertExpectations(t)
	})
}

func TestEnterpriseUsecase_GetDetailEnterpriseByID(t *testing.T) {
	mockEnterpriseRepository := new(mocks.EnterpriseRepository)
	mockTagRepository := new(mocks.TagRepository)
	mockUserRepository := new(mocks.UserRepository)

	t.Run("success", func(t *testing.T) {
		uc := usecase.NewEnterpriseUsecase(mockEnterpriseRepository, mockTagRepository, mockUserRepository)
		mockEnterpriseRepository.On("FindByID", mock.AnythingOfType("string")).Return(dummyEnterprise[0], nil).Once()
		enterprise, err := uc.GetDetailEnterpriseByID(dummyEnterprise[0].ID.String())
		assert.NoError(t, err)
		assert.NotNil(t, enterprise)
		mockEnterpriseRepository.AssertExpectations(t)
	})

	t.Run("failed", func(t *testing.T) {
		uc := usecase.NewEnterpriseUsecase(mockEnterpriseRepository, mockTagRepository, mockUserRepository)
		mockEnterpriseRepository.On("FindByID", mock.AnythingOfType("string")).Return(domain.Enterprise{}, errors.New("error something")).Once()
		_, err := uc.GetDetailEnterpriseByID(dummyEnterprise[0].ID.String())
		assert.Error(t, err)
		mockEnterpriseRepository.AssertExpectations(t)
	})
}

func TestEnterpriseUsecase_GetListAllEnterprise(t *testing.T) {
	mockEnterpriseRepository := new(mocks.EnterpriseRepository)
	mockTagRepository := new(mocks.TagRepository)
	mockUserRepository := new(mocks.UserRepository)

	t.Run("success", func(t *testing.T) {
		uc := usecase.NewEnterpriseUsecase(mockEnterpriseRepository, mockTagRepository, mockUserRepository)
		mockEnterpriseRepository.On("FindAll", mock.AnythingOfType("string"), mock.AnythingOfType("int"), mock.AnythingOfType("int")).Return(domain.Enterprises{
			dummyEnterprise[0],
		}, 1, nil).Once()
		enterprise, totalData, err := uc.GetListAllEnterprise("satu", 1, 1)
		assert.NoError(t, err)
		assert.NotNil(t, enterprise)
		assert.NotNil(t, totalData)
		mockEnterpriseRepository.AssertExpectations(t)
	})

	t.Run("failed", func(t *testing.T) {
		uc := usecase.NewEnterpriseUsecase(mockEnterpriseRepository, mockTagRepository, mockUserRepository)
		mockEnterpriseRepository.On("FindAll", mock.AnythingOfType("string"), mock.AnythingOfType("int"), mock.AnythingOfType("int")).Return(domain.Enterprises{}, 1, errors.New("error something")).Once()
		_, _, err := uc.GetListAllEnterprise("satu", 1, 1)
		assert.Error(t, err)
		mockEnterpriseRepository.AssertExpectations(t)
	})
}

func TestEnterpriseUsecase_GetListEnterpriseByStatus(t *testing.T) {
	mockEnterpriseRepository := new(mocks.EnterpriseRepository)
	mockTagRepository := new(mocks.TagRepository)
	mockUserRepository := new(mocks.UserRepository)

	t.Run("success get list draft", func(t *testing.T) {
		uc := usecase.NewEnterpriseUsecase(mockEnterpriseRepository, mockTagRepository, mockUserRepository)
		mockEnterpriseRepository.On("FindByStatusDraft").Return(dummyEnterprise, nil).Once()
		enterprises, err := uc.GetListEnterpriseByStatus(0)
		assert.NoError(t, err)
		assert.NotNil(t, enterprises)
		mockEnterpriseRepository.AssertExpectations(t)
	})
	t.Run("failed get list draft", func(t *testing.T) {
		uc := usecase.NewEnterpriseUsecase(mockEnterpriseRepository, mockTagRepository, mockUserRepository)
		mockEnterpriseRepository.On("FindByStatusDraft").Return(domain.Enterprises{}, errors.New("error something")).Once()
		_, err := uc.GetListEnterpriseByStatus(0)
		assert.Error(t, err)
		mockEnterpriseRepository.AssertExpectations(t)
	})

	t.Run("success get list publish", func(t *testing.T) {
		uc := usecase.NewEnterpriseUsecase(mockEnterpriseRepository, mockTagRepository, mockUserRepository)
		mockEnterpriseRepository.On("FindByStatusPublish").Return(dummyEnterprise, nil).Once()
		enterprises, err := uc.GetListEnterpriseByStatus(1)
		assert.NoError(t, err)
		assert.NotNil(t, enterprises)
		mockEnterpriseRepository.AssertExpectations(t)
	})
	t.Run("failed get list publish", func(t *testing.T) {
		uc := usecase.NewEnterpriseUsecase(mockEnterpriseRepository, mockTagRepository, mockUserRepository)
		mockEnterpriseRepository.On("FindByStatusPublish").Return(domain.Enterprises{}, errors.New("error something")).Once()
		_, err := uc.GetListEnterpriseByStatus(1)
		assert.Error(t, err)
		mockEnterpriseRepository.AssertExpectations(t)
	})

	t.Run("failed get list", func(t *testing.T) {
		uc := usecase.NewEnterpriseUsecase(mockEnterpriseRepository, mockTagRepository, mockUserRepository)
		_, err := uc.GetListEnterpriseByStatus(2)
		assert.Error(t, err)
		mockEnterpriseRepository.AssertExpectations(t)
	})
}

func TestEnterpriseUsecase_UpdateEnterpriseByID(t *testing.T) {
	mockEnterpriseRepository := new(mocks.EnterpriseRepository)
	mockTagRepository := new(mocks.TagRepository)
	mockUserRepository := new(mocks.UserRepository)

	t.Run("success", func(t *testing.T) {
		req := request.CreateEnterpriseRequest{
			Name:        "enterprise satu",
			NumberPhone: "0012798232",
			Address:     "bjb",
			Postcode:    707722,
			Latitude:    "4235,23",
			Longitude:   "4225,233231",
			Tags:        []string{uuid.FromStringOrNil("35d6a9a1-aa5e-41f1-9991-08878dfdf89b").String()},
		}
		uc := usecase.NewEnterpriseUsecase(mockEnterpriseRepository, mockTagRepository, mockUserRepository)
		mockTagRepository.On("FindByIDs", mock.AnythingOfType("[]string")).Return(domain.Tags{
			domain.Tag{
				ID:   uuid.FromStringOrNil("35d6a9a1-aa5e-41f1-9991-08878dfdf89b"),
				Name: "Tag satu",
			},
		}, nil).Once()
		mockEnterpriseRepository.On("FindByID", mock.AnythingOfType("string")).Return(dummyEnterprise[0], nil).Once()
		mockEnterpriseRepository.On("Update", mock.AnythingOfType("domain.Enterprise")).Return(dummyEnterprise[0], nil).Once()
		enterprise, err := uc.UpdateEnterpriseByID(dummyEnterprise[0].ID.String(), dummyEnterprise[0].UserID.String(), req)
		assert.NoError(t, err)
		assert.NotNil(t, enterprise)
		mockEnterpriseRepository.AssertExpectations(t)
	})

	t.Run("error not current user", func(t *testing.T) {
		req := request.CreateEnterpriseRequest{
			Name:        "enterprise satu",
			NumberPhone: "0012798232",
			Address:     "bjb",
			Postcode:    707722,
			Latitude:    "4235,23",
			Longitude:   "4225,233231",
			Tags:        []string{uuid.FromStringOrNil("35d6a9a1-aa5e-41f1-9991-08878dfdf89b").String()},
		}
		uc := usecase.NewEnterpriseUsecase(mockEnterpriseRepository, mockTagRepository, mockUserRepository)
		mockTagRepository.On("FindByIDs", mock.AnythingOfType("[]string")).Return(domain.Tags{
			domain.Tag{
				ID:   uuid.FromStringOrNil("35d6a9a1-aa5e-41f1-9991-08878dfdf89b"),
				Name: "Tag satu",
			},
		}, nil).Once()
		mockEnterpriseRepository.On("FindByID", mock.AnythingOfType("string")).Return(dummyEnterprise[1], nil).Once()
		_, err := uc.UpdateEnterpriseByID(dummyEnterprise[0].ID.String(), dummyEnterprise[0].UserID.String(), req)
		assert.Error(t, err)
		mockEnterpriseRepository.AssertExpectations(t)
	})

	t.Run("not found list tags", func(t *testing.T) {
		req := request.CreateEnterpriseRequest{
			Name:        "enterprise satu",
			NumberPhone: "0012798232",
			Address:     "bjb",
			Postcode:    707722,
			Latitude:    "4235,23",
			Longitude:   "4225,233231",
			Tags:        []string{uuid.FromStringOrNil("35d6a9a1-aa5e-41f1-9991-08878dfdf89b").String()},
		}
		uc := usecase.NewEnterpriseUsecase(mockEnterpriseRepository, mockTagRepository, mockUserRepository)
		mockTagRepository.On("FindByIDs", mock.AnythingOfType("[]string")).Return(domain.Tags{}, errors.New("not found list tags")).Once()
		_, err := uc.UpdateEnterpriseByID(dummyEnterprise[0].ID.String(), dummyEnterprise[0].UserID.String(), req)
		assert.Error(t, err)
		mockEnterpriseRepository.AssertExpectations(t)
	})

	t.Run("enterprise not found", func(t *testing.T) {
		req := request.CreateEnterpriseRequest{
			Name:        "enterprise satu",
			NumberPhone: "0012798232",
			Address:     "bjb",
			Postcode:    707722,
			Latitude:    "4235,23",
			Longitude:   "4225,233231",
			Tags:        []string{uuid.FromStringOrNil("35d6a9a1-aa5e-41f1-9991-08878dfdf89b").String()},
		}
		uc := usecase.NewEnterpriseUsecase(mockEnterpriseRepository, mockTagRepository, mockUserRepository)
		mockTagRepository.On("FindByIDs", mock.AnythingOfType("[]string")).Return(domain.Tags{
			domain.Tag{
				ID:   uuid.FromStringOrNil("35d6a9a1-aa5e-41f1-9991-08878dfdf89b"),
				Name: "Tag satu",
			},
		}, nil).Once()
		mockEnterpriseRepository.On("FindByID", mock.AnythingOfType("string")).Return(domain.Enterprise{}, errors.New("enterprise not found")).Once()
		_, err := uc.UpdateEnterpriseByID(dummyEnterprise[0].ID.String(), dummyEnterprise[0].UserID.String(), req)
		assert.Error(t, err)
		mockEnterpriseRepository.AssertExpectations(t)
	})
}

func TestEnterpriseUsecase_UpdateStatusEnterprise(t *testing.T) {
	mockEnterpriseRepository := new(mocks.EnterpriseRepository)
	mockTagRepository := new(mocks.TagRepository)
	mockUserRepository := new(mocks.UserRepository)

	t.Run("success", func(t *testing.T) {
		uc := usecase.NewEnterpriseUsecase(mockEnterpriseRepository, mockTagRepository, mockUserRepository)
		mockEnterpriseRepository.On("UpdateStatusByID", mock.AnythingOfType("string"), mock.AnythingOfType("int")).Return(dummyEnterprise[0], nil).Once()
		enterprise, err := uc.UpdateStatusEnterprise(dummyEnterprise[0].ID.String(), 1)
		assert.NoError(t, err)
		assert.NotNil(t, enterprise)
		mockEnterpriseRepository.AssertExpectations(t)
	})

	t.Run("failed", func(t *testing.T) {
		uc := usecase.NewEnterpriseUsecase(mockEnterpriseRepository, mockTagRepository, mockUserRepository)
		mockEnterpriseRepository.On("UpdateStatusByID", mock.AnythingOfType("string"), mock.AnythingOfType("int")).Return(domain.Enterprise{}, errors.New("error change status")).Once()
		_, err := uc.UpdateStatusEnterprise(dummyEnterprise[0].ID.String(), 1)
		assert.Error(t, err)
		mockEnterpriseRepository.AssertExpectations(t)
	})
}
