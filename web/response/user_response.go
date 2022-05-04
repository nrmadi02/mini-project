package response

import (
	uuid "github.com/satori/go.uuid"
	"time"
)

type UserCreateResponse struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Fullname  string    `json:"fullname"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserDetailResponse struct {
	ID          uuid.UUID     `json:"id"`
	Email       string        `json:"email"`
	Fullname    string        `json:"fullname"`
	Username    string        `json:"username"`
	Favorite    interface{}   `json:"favorite"`
	Enterprises []interface{} `json:"enterprises"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
}

type UsersListResponse struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Fullname  string    `json:"fullname"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
