package user

import "github.com/poofik33/db-technopark/internal/models"

type Usecase interface {
	AddUser(nickname string, user *models.User) (*models.User, error)
	GetByNickname(nickname string) (*models.User, error)
	Update(nickname string, user *models.User) error
}
