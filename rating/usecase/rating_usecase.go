package usecase

import (
	"github.com/nrmadi02/mini-project/domain"
	uuid "github.com/satori/go.uuid"
)

type ratingUsecase struct {
	userRepository       domain.UserRepository
	enterpriseRepository domain.EnterpriseRepository
	ratingRepository     domain.RatingRepository
}

func NewRatingUsecase(ur domain.UserRepository, er domain.EnterpriseRepository, rr domain.RatingRepository) domain.RatingUsecase {
	return ratingUsecase{
		userRepository:       ur,
		enterpriseRepository: er,
		ratingRepository:     rr,
	}
}

func (r ratingUsecase) UpdateRating(id, userid string, value int) (domain.RatingEnterprise, error) {
	user, err := r.userRepository.FindUserById(userid)
	if err != nil {
		return domain.RatingEnterprise{}, err
	}
	enterprise, err := r.enterpriseRepository.FindByID(id)
	if err != nil {
		return domain.RatingEnterprise{}, err
	}

	rating, err := r.ratingRepository.UpdateRating(enterprise.ID.String(), user.ID.String(), value)
	if err != nil {
		return domain.RatingEnterprise{}, err
	}

	return rating, err
}

func (r ratingUsecase) DeleteRating(id, userid string) error {
	user, err := r.userRepository.FindUserById(userid)
	if err != nil {
		return err
	}
	enterprise, err := r.enterpriseRepository.FindByID(id)
	if err != nil {
		return err
	}
	rating, err := r.ratingRepository.FindRatingByIDUserAndEnterprise(enterprise.ID.String(), user.ID.String())
	if err != nil {
		return err
	}

	err = r.ratingRepository.DeleteRating(rating)
	if err != nil {
		return err
	}

	return nil
}

func (r ratingUsecase) AddNewRanting(id, userid string, value int) (domain.RatingEnterprise, error) {
	user, err := r.userRepository.FindUserById(userid)
	if err != nil {
		return domain.RatingEnterprise{}, err
	}
	enterprise, err := r.enterpriseRepository.FindByID(id)
	if err != nil {
		return domain.RatingEnterprise{}, err
	}

	req := domain.RatingEnterprise{
		ID:           uuid.NewV4(),
		UserID:       user.ID,
		EnterpriseID: enterprise.ID,
		Rating:       value,
	}

	rating, err := r.ratingRepository.AddRating(req)
	if err != nil {
		return domain.RatingEnterprise{}, err
	}

	return rating, err
}

func (r ratingUsecase) GetAllRatingByEnterpriseID(id string) (domain.RatingEnterprises, error) {
	enterprise, err := r.enterpriseRepository.FindByID(id)
	if err != nil {
		return domain.RatingEnterprises{}, err
	}
	enterprises, err := r.ratingRepository.GetAllRatingByEnterpriseID(enterprise.ID.String())
	if err != nil {
		return nil, err
	}

	return enterprises, err

}

func (r ratingUsecase) FindRating(id, userid string) (domain.RatingEnterprise, error) {
	user, err := r.userRepository.FindUserById(userid)
	if err != nil {
		return domain.RatingEnterprise{}, err
	}
	enterprise, err := r.enterpriseRepository.FindByID(id)
	if err != nil {
		return domain.RatingEnterprise{}, err
	}

	rating, err := r.ratingRepository.FindRatingByIDUserAndEnterprise(enterprise.ID.String(), user.ID.String())
	if err != nil {
		return domain.RatingEnterprise{}, err
	}

	return rating, err
}