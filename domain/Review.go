package domain

import (
	uuid "github.com/satori/go.uuid"
	"time"
)

type Review struct {
	ID           uuid.UUID `json:"id" gorm:"PrimaryKey"`
	Review       string    `json:"review"`
	EnterpriseID uuid.UUID `json:"enterprise_id" gorm:"notnull;type:varchar;size:191"`
	UserID       uuid.UUID `json:"user_id" gorm:"notnull;type:varchar;size:191"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Reviews []Review

type ReviewRepository interface {
	FindByUserIDAndEnterpriseID(enterpriseid, userid string) (Review, error)
	FindByEnterpriseID(id string) (Reviews, error)
	FindByID(id string) (Review, error)
	Update(enterpriseid, userid string, value string) (Review, error)
	Delete(review Review) error
	Add(review Review) (Review, error)
}

type ReviewUsecase interface {
	AddReview(enterpriseid, userid string, value string) (Review, error)
	UpdateReview(enterpriseid, userid string, value string) (Review, error)
	DeleteReview(enterpriseid, userid string) error
	GetListReviewsByEnterpriseID(id string) (Reviews, error)
	GetReviewByUserIDAndEnterpriseID(enterpriseid, userid string) (Review, error)
	GetDetailReviewByID(id string) (Review, error)
}
