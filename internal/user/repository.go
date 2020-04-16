package user

import "github.com/poofik33/db-technopark/internal/models"

type Repository interface {
	InsertInto(user *models.User) error
	GetByNickname(nickname string) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	Update(user *models.User) error
}
