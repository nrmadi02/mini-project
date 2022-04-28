package request

import "errors"

type UserCreateRequest struct {
	Fullname string `json:"fullname"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func ValidateCreation(userRequest UserCreateRequest) (bool, error) {
	if userRequest.Fullname == "" || len(userRequest.Fullname) == 0 {
		return false, errors.New("fullname empty")
	}
	if userRequest.Username == "" || len(userRequest.Username) == 0 {
		return false, errors.New("username empty")
	}
	if userRequest.Email == "" || len(userRequest.Email) < 6 {
		return false, errors.New("email invalid")
	}
	if userRequest.Password == "" || len(userRequest.Password) < 8 {
		return false, errors.New("password, minimum 8 words")
	}
	return true, nil
}
