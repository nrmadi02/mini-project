package response

import (
	uuid "github.com/satori/go.uuid"
	"time"
)

type GetListByStatusResponse struct {
	ID          uuid.UUID          `json:"id" gorm:"PrimaryKey"`
	UserID      uuid.UUID          `json:"user_id" gorm:"notnull;type:varchar;size:191"`
	Name        string             `json:"name" gorm:"notnull"`
	NumberPhone string             `json:"number_phone" gorm:"notnull"`
	Address     string             `json:"address" gorm:"notnull"`
	Postcode    int                `json:"postcode" gorm:"notnull"`
	Description string             `json:"description" gorm:"notnull;type:text"`
	Latitude    string             `json:"latitude" gorm:"null"`
	Longitude   string             `json:"longitude" gorm:"null"`
	Status      int                `json:"status" gorm:"notnull"`
	Tags        interface{}        `json:"tags" gorm:"many2many:enterprise_tags;"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
	Rating      float64            `json:"rating"`
	Owner       UserDetailResponse `json:"owner"`
}
