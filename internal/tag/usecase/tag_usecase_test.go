package usecase_test

import (
	"errors"
	"github.com/nrmadi02/mini-project/domain"
	"github.com/nrmadi02/mini-project/domain/mocks"
	"github.com/nrmadi02/mini-project/internal/tag/usecase"
	"github.com/nrmadi02/mini-project/web/request"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

var dummyTag = domain.Tags{
	domain.Tag{
		ID:   uuid.FromStringOrNil("0cf712fc-e631-40c7-8572-54772e698edf"),
		Name: "Tag Satu",
	},
	domain.Tag{
		ID:   uuid.FromStringOrNil("0cf712fc-e631-40c7-8572-54772e698edd"),
		Name: "Tag Dua",
	},
}

func TestTagUsecase_GetAllTags(t *testing.T) {
	mockTagRepository := new(mocks.TagRepository)

	t.Run("success", func(t *testing.T) {
		uc := usecase.NewTagUsecase(mockTagRepository)
		mockTagRepository.On("FindAllTags").Return(domain.Tags{dummyTag[0]}, nil).Once()
		tags, err := uc.GetAllTags()
		assert.NoError(t, err)
		assert.NotNil(t, tags)
	})
}

func TestTagUsecase_CreateNewTag(t *testing.T) {
	mockTagRepository := new(mocks.TagRepository)

	t.Run("success", func(t *testing.T) {
		req := request.CreateTagRequest{
			Name: "Tag Satu",
		}
		uc := usecase.NewTagUsecase(mockTagRepository)
		mockTagRepository.On("FindByName", mock.AnythingOfType("string")).Return(domain.Tag{}, nil).Once()
		mockTagRepository.On("Save", mock.AnythingOfType("domain.Tag")).Return(dummyTag[0], nil).Once()
		tags, err := uc.CreateNewTag(req)
		assert.NoError(t, err)
		assert.NotNil(t, tags)
	})
	t.Run("error same tag name", func(t *testing.T) {
		req := request.CreateTagRequest{
			Name: "Tag Satu",
		}
		uc := usecase.NewTagUsecase(mockTagRepository)
		mockTagRepository.On("FindByName", mock.AnythingOfType("string")).Return(dummyTag[0], nil).Once()
		_, err := uc.CreateNewTag(req)
		assert.Error(t, err)
	})
	t.Run("error create tag", func(t *testing.T) {
		req := request.CreateTagRequest{
			Name: "Tag Satu",
		}
		uc := usecase.NewTagUsecase(mockTagRepository)
		mockTagRepository.On("FindByName", mock.AnythingOfType("string")).Return(domain.Tag{}, nil).Once()
		mockTagRepository.On("Save", mock.AnythingOfType("domain.Tag")).Return(domain.Tag{}, errors.New("error something")).Once()
		_, err := uc.CreateNewTag(req)
		assert.Error(t, err)
	})
}

func TestTagUsecase_DeleteTag(t *testing.T) {
	mockTagRepository := new(mocks.TagRepository)

	t.Run("success", func(t *testing.T) {
		uc := usecase.NewTagUsecase(mockTagRepository)
		mockTagRepository.On("FindByID", mock.AnythingOfType("string")).Return(dummyTag[0], nil).Once()
		mockTagRepository.On("Delete", mock.AnythingOfType("domain.Tag"), mock.AnythingOfType("string")).Return(nil).Once()
		err := uc.DeleteTag(dummyTag[0].ID.String())
		assert.NoError(t, err)
	})
	t.Run("error tag not found", func(t *testing.T) {
		uc := usecase.NewTagUsecase(mockTagRepository)
		mockTagRepository.On("FindByID", mock.AnythingOfType("string")).Return(domain.Tag{}, errors.New("error something")).Once()
		err := uc.DeleteTag(dummyTag[0].ID.String())
		assert.Error(t, err)
	})
	t.Run("error delete tag", func(t *testing.T) {
		uc := usecase.NewTagUsecase(mockTagRepository)
		mockTagRepository.On("FindByID", mock.AnythingOfType("string")).Return(dummyTag[0], nil).Once()
		mockTagRepository.On("Delete", mock.AnythingOfType("domain.Tag"), mock.AnythingOfType("string")).Return(errors.New("error something")).Once()
		err := uc.DeleteTag(dummyTag[0].ID.String())
		assert.Error(t, err)
	})
}
