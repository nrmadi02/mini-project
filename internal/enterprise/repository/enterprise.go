package repository

import (
	"github.com/nrmadi02/mini-project/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type enterpriseRepository struct {
	DB *gorm.DB
}

func NewEnterpriseRepository(db *gorm.DB) domain.EnterpriseRepository {
	return enterpriseRepository{
		DB: db,
	}
}

func (e enterpriseRepository) FindAll(search string, page, length int) (enterprises domain.Enterprises, totalData int, err error) {
	if page == 0 {
		page = 1
	}
	offset := (page - 1) * length
	if search != "" {
		totalData = int(e.DB.Preload("Tags").Where("name LIKE ?", "%"+search+"%").Find(&enterprises).RowsAffected)
		err = e.DB.Preload("Tags").Offset(offset).Limit(length).Where("name LIKE ?", "%"+search+"%").Find(&enterprises).Error
	} else {
		totalData = int(e.DB.Preload("Tags").Find(&enterprises).RowsAffected)
		err = e.DB.Preload("Tags").Offset(offset).Limit(length).Find(&enterprises).Error
	}
	return enterprises, totalData, err
}

func (e enterpriseRepository) FindByStatusDraft() (enterprises domain.Enterprises, err error) {
	err = e.DB.Preload("Tags").Where("status = ? ", 0).Find(&enterprises).Error
	return enterprises, err
}

func (e enterpriseRepository) FindByStatusPublish() (enterprises domain.Enterprises, err error) {
	err = e.DB.Preload("Tags").Where("status = ? ", 1).Find(&enterprises).Error
	return enterprises, err
}

func (e enterpriseRepository) FindByID(id string) (enterprise domain.Enterprise, err error) {
	err = e.DB.Preload("Tags").Where("id = ?", id).Find(&enterprise).Error
	return enterprise, err
}

func (e enterpriseRepository) UpdateStatusByID(id string, status int) (enterprise domain.Enterprise, err error) {
	err = e.DB.Model(&enterprise).Where("id = ? ", id).Update("status", status).Error
	return enterprise, err
}

func (e enterpriseRepository) Save(enterprise domain.Enterprise) (domain.Enterprise, error) {
	err := e.DB.Create(&enterprise).Error
	return enterprise, err
}

func (e enterpriseRepository) Update(enterprise domain.Enterprise) (domain.Enterprise, error) {
	err := e.DB.Model(&enterprise).Where("id = ? ", enterprise.ID).Updates(&enterprise).Error
	err = e.DB.Model(&enterprise).Association("Tags").Replace(&enterprise.Tags)
	return enterprise, err
}

func (e enterpriseRepository) Delete(enterprise domain.Enterprise) error {
	err := e.DB.Select(clause.Associations).Where("id = ?", enterprise.ID).Delete(&enterprise).Error
	return err
}

func (e enterpriseRepository) FindByIDs(ids []string) (enterprises domain.Enterprises, err error) {
	err = e.DB.Preload("Tags").Where("id IN ? ", ids).Find(&enterprises).Error
	return enterprises, err
}

func (e enterpriseRepository) FindByUserID(id string) (enterprises domain.Enterprises, err error) {
	err = e.DB.Preload("Tags").Where("user_id = ? ", id).Find(&enterprises).Error
	return enterprises, err
}
