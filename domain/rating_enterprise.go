package domain

import uuid "github.com/satori/go.uuid"

type RatingEnterprise struct {
	ID           uuid.UUID `json:"id" gorm:"PrimaryKey"`
	Rating       int       `json:"rating"`
	EnterpriseID uuid.UUID `json:"enterprise_id" gorm:"notnull;type:varchar;size:256"`
	UserID       uuid.UUID `json:"user_id" gorm:"notnull;type:varchar;size:256"`
}

type RatingEnterprises []RatingEnterprise

type RatingRepository interface {
	GetAllRatingByEnterpriseID(id string) (RatingEnterprises, error)
	FindAvgByEnterpriseID(id string) float64
	FindRatingByIDUserAndEnterprise(id string, userid string) (RatingEnterprise, error)
	UpdateRating(id string, userid string, value int) (RatingEnterprise, error)
	DeleteRating(rating RatingEnterprise) error
	AddRating(rating RatingEnterprise) (RatingEnterprise, error)
}

type RatingUsecase interface {
	GetAllRatingByEnterpriseID(id string) (RatingEnterprises, error)
	GetAverageRatingEnterprise(id string) float64
	FindRating(id, userid string) (RatingEnterprise, error)
	UpdateRating(id, userid string, value int) (RatingEnterprise, error)
	DeleteRating(id, userid string) error
	AddNewRanting(id, userid string, value int) (RatingEnterprise, error)
}
