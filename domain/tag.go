package domain

import (
	"github.com/nrmadi02/mini-project/web/request"
	uuid "github.com/satori/go.uuid"
)

type Tag struct {
	ID   uuid.UUID `json:"id" gorm:"PrimaryKey"`
	Name string    `json:"name" gorm:"unique;notnull"`
}

type Tags []Tag

type TagRepository interface {
	FindByName(name string) (Tag, error)
	FindByID(id string) (Tag, error)
	FindByIDs(ids []string) (Tags, error)
	FindAllTags() (Tags, error)
	Delete(tag Tag, id string) error
	Save(tag Tag) (Tag, error)
}

type TagUsecase interface {
	GetAllTags() (Tags, error)
	DeleteTag(id string) error
	CreateNewTag(request request.CreateTagRequest) (Tag, error)
}
