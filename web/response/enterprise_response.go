package response

import (
	uuid "github.com/satori/go.uuid"
	"time"
)

type GetListByStatusResponse struct {
	ID          uuid.UUID          `json:"id"`
	UserID      uuid.UUID          `json:"user_id"`
	Name        string             `json:"name"`
	NumberPhone string             `json:"number_phone"`
	Address     string             `json:"address"`
	Postcode    int                `json:"postcode"`
	Description string             `json:"description"`
	Latitude    string             `json:"latitude"`
	Longitude   string             `json:"longitude"`
	Status      int                `json:"status"`
	Tags        interface{}        `json:"tags,omitempty"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
	Rating      float64            `json:"rating"`
	Owner       UserDetailResponse `json:"owner"`
}
