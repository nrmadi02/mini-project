package request

import "errors"

type CreateTagRequest struct {
	Name string `json:"name"`
}

func ValidateCreationTag(tagRequest CreateTagRequest) (bool, error) {
	if tagRequest.Name == "" || len(tagRequest.Name) == 0 {
		return false, errors.New("name empty")
	}
	return true, nil
}
