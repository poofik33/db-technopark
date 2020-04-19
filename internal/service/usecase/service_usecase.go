package usecase

import (
	"github.com/poofik33/db-technopark/internal/models"
	"github.com/poofik33/db-technopark/internal/service"
)

type ServiceUsecase struct {
	serviceRepo service.Repository
}

func NewServiceUsecase(sr service.Repository) service.Usecase {
	return &ServiceUsecase{
		serviceRepo: sr,
	}
}

func (su *ServiceUsecase) GetStatus() (*models.Status, error) {
	countForum, err := su.serviceRepo.GetCountForum()
	countPost, err := su.serviceRepo.GetCountPost()
	countThread, err := su.serviceRepo.GetCountThread()
	countUser, err := su.serviceRepo.GetCountUser()

	if err != nil {
		return nil, err
	}

	status := &models.Status{
		ForumsCount:  countForum,
		PostsCount:   countPost,
		ThreadsCount: countThread,
		UsersCount:   countUser,
	}
	return status, nil
}

func (su *ServiceUsecase) DeleteAll() error {
	err := su.serviceRepo.DeleteAllVotes()
	err = su.serviceRepo.DeleteAllPost()
	err = su.serviceRepo.DeleteAllThread()
	err = su.serviceRepo.DeleteAllForum()
	err = su.serviceRepo.DeleteAllUser()

	if err != nil {
		return err
	}

	return nil
}
