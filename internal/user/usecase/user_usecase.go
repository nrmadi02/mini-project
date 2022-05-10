package usecase

import (
	"github.com/nrmadi02/mini-project/domain"
)

type userUsecase struct {
	UserRepo domain.UserRepository
}

func NewUserUsecase(ur domain.UserRepository) domain.UserUsecase {
	return userUsecase{
		UserRepo: ur,
	}
}

func (u userUsecase) GetAllUsers() (domain.Users, error) {
	return u.UserRepo.FindAllUsers()
}
