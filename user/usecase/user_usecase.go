package usecase

import (
	"errors"
	"go_cicd/domain"
	"go_cicd/domain/web/request"
	"go_cicd/domain/web/response"
	"go_cicd/user/delivery/http/helper"
)

type userUsecase struct {
	UserRepo domain.UserRepository
}

func NewUserUsecase(ur domain.UserRepository) domain.UserUsecase {
	return &userUsecase{UserRepo: ur}
}

func (u *userUsecase) Login(request request.LoginRequest) (*response.SuccessLogin, error) {
	if request.Email == "" || request.Password == "" {
		return nil, errors.New("email or password empty")
	}
	user := &domain.User{
		Email:    request.Email,
		Password: request.Password,
	}

	resUser, _, err := u.UserRepo.CheckLogin(user)
	if err != nil {
		return nil, errors.New("email or password wrong")
	}

	jwt := helper.NewGoJWT()

	token, err := jwt.CreateTokenJWT(int(resUser.ID), resUser.Email)
	if err != nil {
		return nil, err
	}

	resLogin := &response.SuccessLogin{ID: int(resUser.ID), Email: resUser.Email, Token: token}

	return resLogin, nil

}

func (u *userUsecase) Create(request request.UserCreateRequest) (*domain.User, error) {
	if request.Email == "" {
		return nil, errors.New("email empty")
	}
	user := &domain.User{
		Email:    request.Email,
		Password: request.Password,
	}

	createdUser, err := u.UserRepo.Create(user)
	if err != nil {
		return nil, err
	}

	return createdUser, nil
}

func (u *userUsecase) ReadByID(id int) (*domain.User, error) {
	user, err := u.UserRepo.ReadByID(id)
	if err != nil {
		return nil, err
	}

	return user, err
}

func (u *userUsecase) ReadAll() (*domain.Users, error) {
	foundUsers, err := u.UserRepo.ReadAll()
	if err != nil {
		return nil, err
	}

	return foundUsers, nil
}
