package repository

import (
	"github.com/nrmadi02/mini-project/domain"
	"gorm.io/gorm"
)

type roleRepository struct {
	DB *gorm.DB
}

func NewRoleRepository(db *gorm.DB) domain.RoleRepository {
	return roleRepository{
		DB: db,
	}
}

func (r roleRepository) FindByName(name string) (role domain.Role, err error) {
	err = r.DB.Where("name = ?", name).First(&role).Error
	return role, err
}
