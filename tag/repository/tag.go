package repository

import (
	"github.com/nrmadi02/mini-project/domain"
	"gorm.io/gorm"
)

type tagRepository struct {
	DB *gorm.DB
}

func NewTagRepository(db *gorm.DB) domain.TagRepository {
	return tagRepository{
		DB: db,
	}
}

func (t tagRepository) FindByName(name string) (tag domain.Tag, err error) {
	err = t.DB.Where("name = ? ", name).First(&tag).Error
	return tag, err
}

func (t tagRepository) FindByID(id string) (tag domain.Tag, err error) {
	err = t.DB.Where("id = ? ", id).First(&tag).Error
	return tag, err
}

func (t tagRepository) FindAllTags() (tags domain.Tags, err error) {
	err = t.DB.Find(&tags).Error
	return tags, err
}

func (t tagRepository) Delete(tag domain.Tag, id string) error {
	err := t.DB.Where("id = ? ", id).Delete(&tag).Error
	return err
}

func (t tagRepository) Save(tag domain.Tag) (domain.Tag, error) {
	err := t.DB.Create(&tag).Error
	return tag, err
}
