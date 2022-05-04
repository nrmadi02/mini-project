package repository

import (
	"github.com/nrmadi02/mini-project/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type favoriteRepository struct {
	DB *gorm.DB
}

func NewFavoriteRepository(db *gorm.DB) domain.FavoriteRepository {
	return favoriteRepository{
		DB: db,
	}
}

func (f favoriteRepository) FindAll() (favorites domain.Favorites, err error) {
	err = f.DB.Preload("Enterprises").Find(&favorites).Error
	return favorites, err
}

func (f favoriteRepository) FindByUserID(id string) (favorite domain.Favorite, err error) {
	err = f.DB.Preload(clause.Associations).Where("user_id = ? ", id).Find(&favorite).Error
	return favorite, err
}

func (f favoriteRepository) FindByID(id string) (favorite domain.Favorite, err error) {
	err = f.DB.Preload("Enterprises").Where("id = ? ", id).Find(&favorite).Error
	return favorite, err
}

func (f favoriteRepository) Add(favorite domain.Favorite) (domain.Favorite, error) {
	err := f.DB.Create(&favorite).Error
	return favorite, err
}

func (f favoriteRepository) Update(favorite domain.Favorite, enterprises domain.Enterprises, types string) (domain.Favorite, error) {
	if types == "remove" {
		err := f.DB.Model(&favorite).Association("Enterprises").Delete(enterprises)
		return favorite, err
	} else {
		err := f.DB.Model(&favorite).Where("id = ? ", favorite.ID).Association("Enterprises").Append(enterprises)
		return favorite, err
	}

}
