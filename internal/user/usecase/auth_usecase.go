package usecase

import (
	"errors"
	"github.com/nrmadi02/mini-project/domain"
	"github.com/nrmadi02/mini-project/internal/role/utils"
	"github.com/nrmadi02/mini-project/internal/user/delivery/http/helper"
	request2 "github.com/nrmadi02/mini-project/web/request"
	"github.com/nrmadi02/mini-project/web/response"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

type authUsecase struct {
	userRepository       domain.UserRepository
	roleRepository       domain.RoleRepository
	favoriteRepository   domain.FavoriteRepository
	enterpriseRepository domain.EnterpriseRepository
}

func NewAuthUsecase(ur domain.UserRepository, rr domain.RoleRepository, fr domain.FavoriteRepository, er domain.EnterpriseRepository) domain.AuthUsecase {
	return authUsecase{
		userRepository:       ur,
		roleRepository:       rr,
		favoriteRepository:   fr,
		enterpriseRepository: er,
	}
}

func (a authUsecase) Login(request request2.LoginRequest) (response.SuccessLogin, error) {
	user, err := a.userRepository.FindUserByEmail(request.Email)
	if err != nil {
		return response.SuccessLogin{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))
	if err != nil {
		return response.SuccessLogin{}, errors.New("password wrong")
	}

	jwt := helper.NewGoJWT()
	token, err := jwt.CreateTokenJWT(&user)
	if err != nil {
		return response.SuccessLogin{}, err
	}

	responseLogin := response.SuccessLogin{
		ID:       user.ID,
		Username: user.Username,
		Fullname: user.Fullname,
		Email:    user.Email,
		Token:    token,
	}

	return responseLogin, nil

}

func (a authUsecase) Register(request request2.UserCreateRequest) (domain.User, error) {
	var existingUser domain.User
	existingUser, _ = a.userRepository.FindUserByEmail(request.Email)
	if existingUser.ID != uuid.FromStringOrNil("") {
		return domain.User{}, errors.New("user already exist")
	}

	password, _ := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)

	clientRole, err := a.roleRepository.FindByName(role.Client.String())
	if err != nil {
		return domain.User{}, errors.New("role not found - " + role.Client.String())
	}

	user := domain.User{
		ID:       uuid.NewV4(),
		Fullname: request.Fullname,
		Username: request.Username,
		Email:    request.Email,
		Password: string(password),
		Roles:    []domain.Role{clientRole},
	}
	favorite := domain.Favorite{
		ID:     uuid.NewV4(),
		UserID: user.ID,
	}

	user, err = a.userRepository.Save(user)
	if err != nil {
		return domain.User{}, err
	}
	if err != nil {
		return domain.User{}, err
	}
	_, _ = a.favoriteRepository.Add(favorite)

	return user, nil

}

func (a authUsecase) GetUserDetails(id string) (domain.User, domain.Favorite, domain.Enterprises, error) {
	var user domain.User

	user, err := a.userRepository.FindUserById(id)
	if err != nil {
		return domain.User{}, domain.Favorite{}, nil, err
	}
	favorite, _ := a.favoriteRepository.FindByUserID(user.ID.String())
	enterprises, _ := a.enterpriseRepository.FindByUserID(user.ID.String())

	return user, favorite, enterprises, nil
}

func (a authUsecase) CheckIfUserIsAdmin(id string) (bool, error) {
	user, err := a.userRepository.FindUserById(id)
	if err != nil {
		return false, err
	}

	for _, a := range user.Roles {
		if a.Name == "ROLE_ADMIN" {
			return true, nil
		}
	}
	return false, nil

}
