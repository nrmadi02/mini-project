package usecase

import (
	"github.com/nrmadi02/mini-project/domain"
)

type userUsecase struct {
	UserRepo    domain.UserRepository
	authUsecase domain.AuthUsecase
}

func NewUserUsecase(ur domain.UserRepository, au domain.AuthUsecase) domain.UserUsecase {
	return userUsecase{
		UserRepo:    ur,
		authUsecase: au,
	}
}

func (u userUsecase) GetAllUsers() (domain.Users, error) {
	return u.UserRepo.FindAllUsers()
}
