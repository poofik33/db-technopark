package usecase

import (
	"errors"
	"github.com/poofik33/db-technopark/internal/models"
	"github.com/poofik33/db-technopark/internal/tools"
	"github.com/poofik33/db-technopark/internal/user"
)

type UserUsecase struct {
	ur user.Repository
}

func NewUserUsecase(ur user.Repository) user.Usecase {
	return &UserUsecase{
		ur: ur,
	}
}

func (uu *UserUsecase) AddUser(nickname string, user *models.User) (*models.User, error) {
	u, err := uu.ur.GetByNickname(nickname)
	if err != nil && err != tools.ErrDoesntExists {
		return nil, err
	}
	u, err = uu.ur.GetByEmail(user.Email)
	if err != nil  && err != tools.ErrDoesntExists {
		return nil, err
	}
	if u != nil {
		return u, errors.New("Exists")
	}

	user.SetNickname(nickname)
	if err = uu.ur.InsertInto(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (uu *UserUsecase) GetByNickname(nickname string) (*models.User, error) {
	u, err := uu.ur.GetByNickname(nickname)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (uu *UserUsecase) Update(nickname string, user *models.User) error {
	u, err := uu.ur.GetByNickname(nickname)
	if err != nil {
		return err
	}
	if u == nil {
		return errors.New("Doesn't Exists")
	}

	newEmailCheckUser, err := uu.ur.GetByEmail(user.Email)
	if err != nil && newEmailCheckUser == nil{
		return err
	}
	if err != tools.ErrDoesntExists && newEmailCheckUser.Nickname != u.Nickname {
		return tools.ErrExistWithEmail
	}

	user.SetNickname(nickname)
	if err = uu.ur.Update(user); err != nil {
		return err
	}

	return nil
}
