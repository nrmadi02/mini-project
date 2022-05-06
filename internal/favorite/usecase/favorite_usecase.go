package usecase

import (
	"errors"
	"fmt"
	"github.com/nrmadi02/mini-project/domain"
	uuid "github.com/satori/go.uuid"
)

type favoriteUsecase struct {
	enterprisesRepository domain.EnterpriseRepository
	favoriteRepository    domain.FavoriteRepository
}

func NewFavoriteUsecase(er domain.EnterpriseRepository, fr domain.FavoriteRepository) domain.FavoriteUsecase {
	return favoriteUsecase{
		enterprisesRepository: er,
		favoriteRepository:    fr,
	}
}

func (f favoriteUsecase) GetDetailByID(id string) (domain.Favorite, error) {
	favorite, err := f.favoriteRepository.FindByID(id)
	if err != nil {
		return domain.Favorite{}, err
	}

	return favorite, err
}

func (f favoriteUsecase) GetDetailByUserID(id string) (domain.Favorite, error) {
	favorite, err := f.favoriteRepository.FindByUserID(id)
	if favorite.ID == uuid.FromStringOrNil("") {
		return domain.Favorite{}, errors.New("user not found")
	}

	return favorite, err
}

func (f favoriteUsecase) AddFavorite(ids []string, userid string) (domain.Favorite, error) {
	enterprises, _ := f.enterprisesRepository.FindByIDs(ids)
	if len(enterprises) == 0 {
		return domain.Favorite{}, errors.New("enterprise not found")
	}
	favorite, err := f.favoriteRepository.FindByUserID(userid)
	if err != nil {
		return domain.Favorite{}, err
	}
	update, err := f.favoriteRepository.Update(favorite, enterprises, "add")
	if err != nil {
		return domain.Favorite{}, err
	}
	return update, err
}

func (f favoriteUsecase) RemoveFavorite(ids []string, userid string) (domain.Favorite, error) {
	enterprises, err := f.enterprisesRepository.FindByIDs(ids)
	fmt.Println(enterprises)
	if err != nil {
		return domain.Favorite{}, err
	}
	favorite, err := f.favoriteRepository.FindByUserID(userid)
	if err != nil {
		return domain.Favorite{}, err
	}
	update, err := f.favoriteRepository.Update(favorite, enterprises, "remove")
	if err != nil {
		return domain.Favorite{}, err
	}
	return update, err
}
