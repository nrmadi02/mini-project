package usecase

import (
	"errors"
	"github.com/nrmadi02/mini-project/domain"
	"github.com/nrmadi02/mini-project/domain/web/request"
	"github.com/nrmadi02/mini-project/domain/web/response"
	role "github.com/nrmadi02/mini-project/role/utils"
	"golang.org/x/crypto/bcrypt"
)

type authUsecase struct {
	userRepository domain.UserRepository
	roleRepository domain.RoleRepository
}

func NewAuthUsecase(ur domain.UserRepository, rr domain.RoleRepository) domain.AuthUsecase {
	return authUsecase{
		userRepository: ur,
		roleRepository: rr,
	}
}

func (a authUsecase) Login(request request.LoginRequest) (response.SuccessLogin, error) {
	//TODO implement me
	panic("implement me")
}

func (a authUsecase) Register(request request.UserCreateRequest) (domain.User, error) {
	var existingUser domain.User
	existingUser, _ = a.userRepository.FindUserByEmail(request.Email)
	if existingUser.ID != "" {
		return domain.User{}, errors.New("user already exist")
	}

	password, _ := bcrypt.GenerateFromPassword([]byte(request.Password), 12)

	clientRole, err := a.roleRepository.FindByName(role.Client.String())
	if err != nil {
		return domain.User{}, errors.New("role not found - " + role.Client.String())
	}

	user := domain.User{
		Fullname: request.Fullname,
		Username: request.Username,
		Email:    request.Email,
		Password: password,
		Roles:    []domain.Role{clientRole},
	}

	user, err = a.userRepository.Save(user)
	if err != nil {
		return domain.User{}, err
	}

	return user, nil

}

func (a authUsecase) GetUserDetails(id string) (domain.User, error) {
	//TODO implement me
	panic("implement me")
}

func (a authUsecase) CheckIfUserIsAdmin(id string) (bool, error) {
	//TODO implement me
	panic("implement me")
}
