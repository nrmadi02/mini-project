package response

type SuccessLogin struct {
	ID    int    `json:"id" form:"id"`
	Email string `json:"email" form:"email"`
	Token string `json:"token" form:"token"`
}
