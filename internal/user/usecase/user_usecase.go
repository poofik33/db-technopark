package usecase

import (
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

func (uu *UserUsecase) AddUser(nickname string, user *models.User) ([]*models.User, error) {
	u1, err := uu.ur.GetByNickname(nickname)
	if err != nil && err != tools.ErrDoesntExists {
		return nil, err
	}
	u2, err := uu.ur.GetByEmail(user.Email)
	if err != nil && err != tools.ErrDoesntExists {
		return nil, err
	}

	if u1 != nil || u2 != nil {
		returnUsers := []*models.User{}
		if u1 != nil {
			returnUsers = append(returnUsers, u1)
			if u2 != nil && u1.Nickname != u2.Nickname {
				returnUsers = append(returnUsers, u2)
			}
		} else if u2 != nil {
			returnUsers = append(returnUsers, u2)
		}
		return returnUsers, tools.ErrUserExistWith
	}

	user.SetNickname(nickname)
	if err = uu.ur.InsertInto(user); err != nil {
		return nil, err
	}

	return []*models.User{user}, nil
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
		if err == tools.ErrDoesntExists {
			return tools.ErrUserDoesntExists
		}
		return err
	}

	newEmailCheckUser, err := uu.ur.GetByEmail(user.Email)
	if err != nil && err != tools.ErrDoesntExists {
		return err
	}
	if err != tools.ErrDoesntExists && newEmailCheckUser.Nickname != u.Nickname {
		return tools.ErrUserExistWith
	}

	user.SetNickname(nickname)
	if user.Email == "" {
		user.Email = u.Email
	}
	if user.Fullname == "" {
		user.Fullname = u.Fullname
	}
	if user.About == "" {
		user.About = u.About
	}

	if err = uu.ur.Update(user); err != nil {
		return err
	}

	return nil
}
