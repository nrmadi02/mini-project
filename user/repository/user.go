package repository

import (
	"go_cicd/domain"
	"gorm.io/gorm"
)

type userRepository struct {
	Conn *gorm.DB
}

func NewUserRepository(Conn *gorm.DB) domain.UserRepository {
	return &userRepository{Conn: Conn}
}

func (u *userRepository) CheckLogin(user *domain.User) (*domain.User, bool, error) {
	if err := u.Conn.Where("email = ? AND password = ?", user.Email, user.Password).First(&user).Error; err != nil {
		return nil, false, err
	}

	return user, true, nil
}

func (u *userRepository) Create(user *domain.User) (*domain.User, error) {
	if err := u.Conn.Create(&user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (u *userRepository) ReadByID(id int) (*domain.User, error) {
	user := &domain.User{ID: id}
	if err := u.Conn.First(&user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (u *userRepository) ReadAll() (*domain.Users, error) {
	users := &domain.Users{}
	u.Conn.Find(&users)

	return users, nil
}
