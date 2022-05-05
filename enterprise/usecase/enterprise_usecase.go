package usecase

import (
	"errors"
	"github.com/nrmadi02/mini-project/domain"
	request2 "github.com/nrmadi02/mini-project/web/request"
	uuid "github.com/satori/go.uuid"
)

type enterpriseUsecase struct {
	enterpriseRepository domain.EnterpriseRepository
	tagRepository        domain.TagRepository
	userRepository       domain.UserRepository
}

func NewEnterpriseUsecase(er domain.EnterpriseRepository, tr domain.TagRepository, ur domain.UserRepository) domain.EnterpriseUsecase {
	return enterpriseUsecase{
		enterpriseRepository: er,
		tagRepository:        tr,
		userRepository:       ur,
	}
}

func (e enterpriseUsecase) CreateNewEnterprise(request request2.CreateEnterpriseRequest, userid string) (domain.Enterprise, error) {
	tagsList, err := e.tagRepository.FindByIDs(request.Tags)
	if err != nil {
		return domain.Enterprise{}, err
	}

	user, err := e.userRepository.FindUserById(userid)
	if err != nil {
		return domain.Enterprise{}, err
	}

	reqBody := domain.Enterprise{
		ID:          uuid.NewV4(),
		UserID:      user.ID,
		Name:        request.Name,
		Address:     request.Address,
		NumberPhone: request.NumberPhone,
		Postcode:    request.Postcode,
		Description: request.Description,
		Latitude:    request.Latitude,
		Longitude:   request.Longitude,
		Status:      0,
		Tags:        tagsList,
	}

	res, err := e.enterpriseRepository.Save(reqBody)
	if err != nil {
		return domain.Enterprise{}, err
	}
	return res, err
}

func (e enterpriseUsecase) UpdateStatusEnterprise(id string, status int) (domain.Enterprise, error) {
	res, err := e.enterpriseRepository.UpdateStatusByID(id, status)
	if err != nil {
		return domain.Enterprise{}, err
	}
	return res, err
}

func (e enterpriseUsecase) GetDetailEnterpriseByID(id string) (domain.Enterprise, error) {
	enterprise, err := e.enterpriseRepository.FindByID(id)
	if err != nil {
		return domain.Enterprise{}, err
	}
	return enterprise, err
}

func (e enterpriseUsecase) GetListEnterpriseByStatus(status int) (domain.Enterprises, error) {
	var enterprises domain.Enterprises
	var err error
	if status == 0 {
		enterprises, err = e.enterpriseRepository.FindByStatusDraft()
	} else if status == 1 {
		enterprises, err = e.enterpriseRepository.FindByStatusPublish()
	} else {
		return domain.Enterprises{}, errors.New("something wrong")
	}

	if err != nil {
		return domain.Enterprises{}, err
	}

	return enterprises, err
}

func (e enterpriseUsecase) UpdateEnterpriseByID(id string, userid string, request request2.CreateEnterpriseRequest) (domain.Enterprise, error) {
	tagsList, err := e.tagRepository.FindByIDs(request.Tags)
	if err != nil {
		return domain.Enterprise{}, err
	}

	enterpriseByID, err := e.enterpriseRepository.FindByID(id)
	if err != nil {
		return domain.Enterprise{}, err
	}

	if enterpriseByID.ID.String() != userid {
		if err != nil {
			return domain.Enterprise{}, errors.New("to update enterprise must current user")
		}
	}

	req := domain.Enterprise{
		ID:          enterpriseByID.ID,
		UserID:      enterpriseByID.UserID,
		Name:        request.Name,
		Address:     request.Address,
		NumberPhone: request.NumberPhone,
		Postcode:    request.Postcode,
		Description: request.Description,
		Latitude:    request.Latitude,
		Longitude:   request.Longitude,
		Status:      enterpriseByID.Status,
		Tags:        tagsList,
	}
	res, err := e.enterpriseRepository.Update(enterpriseByID, req)
	return res, err
}

func (e enterpriseUsecase) GetListAllEnterprise(search string) (domain.Enterprises, error) {
	enterprises, err := e.enterpriseRepository.FindAll(search)
	if err != nil {
		return domain.Enterprises{}, err
	}
	return enterprises, err
}

func (e enterpriseUsecase) DeleteEnterpriseByID(id string) error {
	enterprise, err := e.enterpriseRepository.FindByID(id)
	if err != nil {
		return err
	}

	err = e.enterpriseRepository.Delete(enterprise)
	if err != nil {
		return err
	}

	return err
}
