package domain

import (
	request2 "github.com/nrmadi02/mini-project/web/request"
	uuid "github.com/satori/go.uuid"
	"time"
)

type Enterprise struct {
	ID               uuid.UUID          `json:"id" gorm:"PrimaryKey"`
	UserID           uuid.UUID          `json:"user_id" gorm:"notnull;type:varchar;size:191"`
	Name             string             `json:"name" gorm:"notnull"`
	NumberPhone      string             `json:"number_phone" gorm:"notnull"`
	Address          string             `json:"address" gorm:"notnull"`
	Postcode         int                `json:"postcode" gorm:"notnull"`
	Latitude         string             `json:"latitude" gorm:"null"`
	Longitude        string             `json:"longitude" gorm:"null"`
	Description      string             `json:"description" gorm:"notnull;type:text"`
	Status           int                `json:"status" gorm:"notnull"`
	Tags             []Tag              `json:"tags" gorm:"many2many:enterprise_tags;"`
	RatingEnterprise []RatingEnterprise `json:"rating_enterprise" gorm:"foreignKey:EnterpriseID;references:ID"`
	CreatedAt        time.Time          `json:"created_at"`
	UpdatedAt        time.Time          `json:"updated_at"`
}

type Enterprises []Enterprise

type EnterpriseRepository interface {
	FindByID(id string) (Enterprise, error)
	FindAll() (Enterprises, error)
	FindByStatusDraft() (Enterprises, error)
	FindByStatusPublish() (Enterprises, error)
	UpdateStatusByID(id string, status int) (Enterprise, error)
	Update(enterprise Enterprise, req Enterprise) (Enterprise, error)
	Save(enterprise Enterprise) (Enterprise, error)
	Delete(enterprise Enterprise) error
}

type EnterpriseUsecase interface {
	CreateNewEnterprise(request request2.CreateEnterpriseRequest, userid string) (Enterprise, error)
	UpdateStatusEnterprise(id string, status int) (Enterprise, error)
	UpdateEnterpriseByID(id string, userid string, request request2.CreateEnterpriseRequest) (Enterprise, error)
	GetDetailEnterpriseByID(id string) (Enterprise, error)
	GetListEnterpriseByStatus(status int) (Enterprises, error)
	GetListAllEnterprise() (Enterprises, error)
	DeleteEnterpriseByID(id string) error
}
