package usecase

import (
	"errors"
	"github.com/nrmadi02/mini-project/domain"
	"github.com/nrmadi02/mini-project/web/request"
	uuid "github.com/satori/go.uuid"
)

type tagUsecase struct {
	tagRepository domain.TagRepository
}

func NewTagUsecase(tg domain.TagRepository) domain.TagUsecase {
	return tagUsecase{
		tagRepository: tg,
	}
}

func (t tagUsecase) GetAllTags() (domain.Tags, error) {
	return t.tagRepository.FindAllTags()
}

func (t tagUsecase) DeleteTag(id string) error {
	tag, err := t.tagRepository.FindByID(id)
	if err != nil {
		return err
	}
	err = t.tagRepository.Delete(tag, id)
	if err != nil {
		return err
	}

	return nil
}

func (t tagUsecase) CreateNewTag(request request.CreateTagRequest) (domain.Tag, error) {
	var existingTag domain.Tag
	existingTag, _ = t.tagRepository.FindByName(request.Name)
	if existingTag.ID != uuid.FromStringOrNil("") {
		return domain.Tag{}, errors.New("tag already exist")
	}
	tagBody := domain.Tag{
		ID:   uuid.NewV4(),
		Name: request.Name,
	}

	tag, err := t.tagRepository.Save(tagBody)
	if err != nil {
		return domain.Tag{}, err
	}

	return tag, nil
}
