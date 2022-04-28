package response

import uuid "github.com/satori/go.uuid"

type UserCreateResponse struct {
	ID       uuid.UUID `json:"id"`
	Email    string    `json:"email"`
	Fullname string    `json:"fullname"`
	Username string    `json:"username"`
}

type UserDetailResponse struct {
	ID       uuid.UUID `json:"id"`
	Email    string    `json:"email"`
	Fullname string    `json:"fullname"`
	Username string    `json:"username"`
}

type UsersListResponse struct {
	ID       uuid.UUID `json:"id"`
	Email    string    `json:"email"`
	Fullname string    `json:"fullname"`
	Username string    `json:"username"`
}
