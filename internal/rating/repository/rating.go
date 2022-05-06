package repository

import (
	"github.com/nrmadi02/mini-project/domain"
	"gorm.io/gorm"
)

type ratingRepository struct {
	DB *gorm.DB
}

func NewRatingRepository(db *gorm.DB) domain.RatingRepository {
	return ratingRepository{
		DB: db,
	}
}

func (r ratingRepository) GetAllRatingByEnterpriseID(id string) (ratings domain.RatingEnterprises, err error) {
	err = r.DB.Where("enterprise_id = ? ", id).Find(&ratings).Error
	return ratings, err
}

func (r ratingRepository) FindRatingByIDUserAndEnterprise(id string, userid string) (rating domain.RatingEnterprise, err error) {
	err = r.DB.Where("enterprise_id = ? AND user_id = ?", id, userid).Find(&rating).Error
	return rating, err
}

func (r ratingRepository) AddRating(rating domain.RatingEnterprise) (domain.RatingEnterprise, error) {
	err := r.DB.Create(&rating).Error
	return rating, err
}

func (r ratingRepository) UpdateRating(id string, userid string, value int) (rating domain.RatingEnterprise, err error) {
	err = r.DB.Model(&rating).Where("enterprise_id = ? AND user_id = ? ", id, userid).Update("rating", value).Error
	return rating, err
}

func (r ratingRepository) DeleteRating(rating domain.RatingEnterprise) error {
	err := r.DB.Where("enterprise_id = ? AND user_id = ? ", rating.EnterpriseID, rating.UserID).Delete(&rating).Error
	return err
}
