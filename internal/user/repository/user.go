package repository

import (
	"github.com/nrmadi02/mini-project/domain"
	"gorm.io/gorm"
)

type userRepository struct {
	Conn *gorm.DB
}

func NewUserRepository(Conn *gorm.DB) domain.UserRepository {
	return &userRepository{Conn: Conn}
}

func (u userRepository) FindUserByEmail(email string) (user domain.User, err error) {
	err = u.Conn.Preload("Roles").Where("email = ?", email).First(&user).Error
	return user, err
}

func (u userRepository) FindUserById(id string) (user domain.User, err error) {
	err = u.Conn.Preload("Roles").Where("id = ?", id).First(&user).Error
	return user, err
}

func (u userRepository) Save(user domain.User) (domain.User, error) {
	err := u.Conn.Create(&user).Error
	return user, err
}

func (u userRepository) FindAllUsers() (users domain.Users, err error) {
	err = u.Conn.Preload("Roles").Find(&users).Error
	return users, err
}
