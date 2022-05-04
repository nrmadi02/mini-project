package repository

import (
	"github.com/nrmadi02/mini-project/domain"
	"gorm.io/gorm"
)

type enterpriseRepository struct {
	DB *gorm.DB
}

func NewEnterpriseRepository(db *gorm.DB) domain.EnterpriseRepository {
	return enterpriseRepository{
		DB: db,
	}
}

func (e enterpriseRepository) FindAll() (enterprises domain.Enterprises, err error) {
	err = e.DB.Preload("Tags").Find(&enterprises).Error
	return enterprises, err
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

func (e enterpriseRepository) Update(enterprise domain.Enterprise, req domain.Enterprise) (domain.Enterprise, error) {
	if len(enterprise.Tags) != 0 {
		err := e.DB.Model(&enterprise).Association("Tags").Clear()
		if err != nil {
			return domain.Enterprise{}, err
		}
	}
	err := e.DB.Model(&enterprise).Where("id = ? ", req.ID).Updates(map[string]interface{}{
		"name":         req.Name,
		"number_phone": req.NumberPhone,
		"address":      req.Address,
		"postcode":     req.Postcode,
		"description":  req.Description,
		"latitude":     req.Latitude,
		"longitude":    req.Longitude,
	}).Association("Tags").Append(req.Tags)
	return enterprise, err
}

func (e enterpriseRepository) Delete(enterprise domain.Enterprise) error {
	err := e.DB.Model(&enterprise).Association("Tags").Clear()
	if err != nil {
		return err
	}
	err = e.DB.Where("id = ? ", enterprise.ID).Delete(&enterprise).Error
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
