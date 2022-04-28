package response

import uuid "github.com/satori/go.uuid"

type SuccessLogin struct {
	ID       uuid.UUID `json:"id"`
	Email    string    `json:"email"`
	Fullname string    `json:"fullname"`
	Username string    `json:"username"`
	Token    string    `json:"token" form:"token"`
}
