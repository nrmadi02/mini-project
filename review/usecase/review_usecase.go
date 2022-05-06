package usecase

import (
	"errors"
	"github.com/nrmadi02/mini-project/domain"
	uuid "github.com/satori/go.uuid"
)

type reviewUsecase struct {
	enterpriseRepository domain.EnterpriseRepository
	userRepository       domain.UserRepository
	reviewRepository     domain.ReviewRepository
}

func NewReviewUsecase(er domain.EnterpriseRepository, ur domain.UserRepository, rr domain.ReviewRepository) domain.ReviewUsecase {
	return reviewUsecase{
		enterpriseRepository: er,
		userRepository:       ur,
		reviewRepository:     rr,
	}
}

func (r reviewUsecase) AddReview(enterpriseid, userid string, value string) (domain.Review, error) {
	enterprise, _ := r.enterpriseRepository.FindByID(enterpriseid)
	if enterprise.ID == uuid.FromStringOrNil("") {
		return domain.Review{}, errors.New("enterprise not found")
	}
	user, _ := r.userRepository.FindUserById(userid)
	if user.ID == uuid.FromStringOrNil("") {
		return domain.Review{}, errors.New("user not found")
	}

	req := domain.Review{
		ID:           uuid.NewV4(),
		UserID:       user.ID,
		EnterpriseID: enterprise.ID,
		Review:       value,
	}

	add, err := r.reviewRepository.Add(req)
	if err != nil {
		return domain.Review{}, err
	}

	return add, nil
}

func (r reviewUsecase) UpdateReview(enterpriseid, userid string, value string) (domain.Review, error) {
	review, _ := r.reviewRepository.FindByUserIDAndEnterpriseID(enterpriseid, userid)
	if review.ID == uuid.FromStringOrNil("") {
		return domain.Review{}, errors.New("request enterprise and user")
	}

	update, err := r.reviewRepository.Update(review.EnterpriseID.String(), review.UserID.String(), value)
	if err != nil {
		return domain.Review{}, err
	}
	return update, nil
}

func (r reviewUsecase) DeleteReview(enterpriseid, userid string) error {
	review, _ := r.reviewRepository.FindByUserIDAndEnterpriseID(enterpriseid, userid)
	if review.ID == uuid.FromStringOrNil("") {
		return errors.New("request enterprise and user")
	}

	err := r.reviewRepository.Delete(review)
	return err
}

func (r reviewUsecase) GetListReviewsByEnterpriseID(id string) (domain.Reviews, error) {
	reviews, err := r.reviewRepository.FindByEnterpriseID(id)
	if err != nil {
		return nil, err
	}

	return reviews, nil
}

func (r reviewUsecase) GetReviewByUserIDAndEnterpriseID(enterpriseid, userid string) (domain.Review, error) {
	review, _ := r.reviewRepository.FindByUserIDAndEnterpriseID(enterpriseid, userid)
	if review.ID == uuid.FromStringOrNil("") {
		return domain.Review{}, errors.New("request enterprise and user")
	}

	return review, nil
}

func (r reviewUsecase) GetDetailReviewByID(id string) (domain.Review, error) {
	review, err := r.reviewRepository.FindByID(id)
	if err != nil {
		return domain.Review{}, err
	}
	return review, err
}
