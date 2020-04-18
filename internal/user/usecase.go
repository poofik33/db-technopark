package user

import "github.com/poofik33/db-technopark/internal/models"

type Usecase interface {
	AddUser(string, *models.User) (*models.User, error)
	GetByNickname(string) (*models.User, error)
	Update(string, *models.User) error
}
