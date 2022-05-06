package repository

import (
	"github.com/nrmadi02/mini-project/domain"
	"gorm.io/gorm"
)

type reviewRepository struct {
	DB *gorm.DB
}

func NewReviewRepository(db *gorm.DB) domain.ReviewRepository {
	return reviewRepository{
		DB: db,
	}
}

func (r reviewRepository) FindByUserIDAndEnterpriseID(enterpriseid, userid string) (review domain.Review, err error) {
	err = r.DB.Where("enterprise_id = ? AND user_id = ?", enterpriseid, userid).Find(&review).Error
	return review, err
}

func (r reviewRepository) FindByEnterpriseID(id string) (reviews domain.Reviews, err error) {
	err = r.DB.Where("enterprise_id = ?", id).Find(&reviews).Error
	return reviews, err
}

func (r reviewRepository) FindByID(id string) (review domain.Review, err error) {
	err = r.DB.Where("id = ?", id).Find(&review).Error
	return review, err
}

func (r reviewRepository) Update(enterpriseid, userid string, value string) (review domain.Review, err error) {
	err = r.DB.Model(&review).Where("enterprise_id = ? AND user_id = ?", enterpriseid, userid).Update("review", value).Error
	return review, err
}

func (r reviewRepository) Delete(review domain.Review) error {
	err := r.DB.Delete(&review).Error
	return err
}

func (r reviewRepository) Add(review domain.Review) (domain.Review, error) {
	err := r.DB.Create(&review).Error
	return review, err
}
