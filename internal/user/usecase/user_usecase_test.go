package usecase_test

import (
	"errors"
	"github.com/nrmadi02/mini-project/domain"
	"github.com/nrmadi02/mini-project/domain/mocks"
	"github.com/nrmadi02/mini-project/internal/user/usecase"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"testing"
	"time"
)

var password, _ = bcrypt.GenerateFromPassword([]byte("12345678"), bcrypt.DefaultCost)
var dummyUser = domain.Users{
	domain.User{
		ID:       uuid.FromStringOrNil("0cf712fc-e631-40c7-8572-54772e698edg"),
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

func TestUserUsecase_GetAllUsers(t *testing.T) {
	mockUserRepository := new(mocks.UserRepository)

	t.Run("success", func(t *testing.T) {
		mockUserRepository.On("FindAllUsers").Return(dummyUser, nil).Once()
		uc := usecase.NewUserUsecase(mockUserRepository)
		res, err := uc.GetAllUsers()
		assert.NoError(t, err)
		assert.Len(t, res, len(dummyUser))
		mockUserRepository.AssertExpectations(t)
	})

	t.Run("error-failed", func(t *testing.T) {
		mockUserRepository.On("FindAllUsers").Return(nil, errors.New("error something")).Once()
		uc := usecase.NewUserUsecase(mockUserRepository)
		_, err := uc.GetAllUsers()
		assert.Error(t, err)
		mockUserRepository.AssertExpectations(t)
	})
}
