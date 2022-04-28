package response

import uuid "github.com/satori/go.uuid"

type TagsListResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type TagCreateResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}
