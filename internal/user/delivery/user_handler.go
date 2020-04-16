package delivery

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/poofik33/db-technopark/internal/models"
	"github.com/poofik33/db-technopark/internal/tools"
	"github.com/poofik33/db-technopark/internal/user"
	"github.com/sirupsen/logrus"
	"net/http"
)

type UserHandler struct {
	userUC user.Usecase
}

func NewUserHandler(e *echo.Echo, uuc user.Usecase) *UserHandler {
	uh := &UserHandler{
		userUC:uuc,
	}

	e.POST("/user/:nickname/create", uh.CreateUser())
	e.GET("/user/:nickname/profile", uh.GetProfile())
	e.POST("/user/:nickname/profile", uh.UpdateProfile())

	return uh
}

func (uh *UserHandler) CreateUser() echo.HandlerFunc {
	type createUserRequset struct {
		Email string `json:"email" binding:"required" validate:"email"`
		Fullname string `json:"fullname" binding:"required"`
		About string `json:"about"`
	}

	return func(c echo.Context) error {
		req := &createUserRequset{}
		if err := c.Bind(req); err != nil {
			logrus.Error(fmt.Errorf("Binding error %s", err))
			return c.JSON(http.StatusBadRequest, tools.ErrorResponce{err.Error()})
		}

		if err := c.Validate(req); err != nil {
			logrus.Error(fmt.Errorf("Validate error %s", err))
			return c.JSON(http.StatusBadRequest, tools.ErrorResponce{err.Error()})
		}

		nickname := c.Param("nickname")

		u :=  &models.User{
			Email:    req.Email,
			Fullname: req.Fullname,
			About:    req.About,
		}

		returnUser, err := uh.userUC.AddUser(nickname, u)
		if err != nil && returnUser != nil {
			return c.JSON(http.StatusConflict, []*models.User{returnUser,})
		}

		if err != nil {
			logrus.Error(fmt.Errorf("Request error %s", err))
			return c.JSON(http.StatusBadRequest,
				tools.ErrorResponce{err.Error()})
		}

		return c.JSON(http.StatusCreated, returnUser)
	}
}

func (uh *UserHandler) GetProfile() echo.HandlerFunc {
	return func(c echo.Context) error {
		nickname := c.Param("nickname")

		returnUser, err := uh.userUC.GetByNickname(nickname)
		if err != nil && err != tools.ErrDoesntExists {
			logrus.Error(fmt.Errorf("Request error %s", err))
			return c.JSON(http.StatusBadRequest, tools.ErrorResponce{err.Error()})
		}

		if err == tools.ErrDoesntExists {
			return c.JSON(http.StatusNotFound, tools.ErrorResponce{err.Error()})
		}

		return c.JSON(http.StatusOK, returnUser)
	}
}

func (uh *UserHandler) UpdateProfile() echo.HandlerFunc {
	type updateUserRequset struct {
		Email string `json:"email" binding:"required" validate:"email"`
		Fullname string `json:"fullname" binding:"required"`
		About string `json:"about"`
	}

	return func(c echo.Context) error {
		req := &updateUserRequset{}
		if err := c.Bind(req); err != nil {
			logrus.Error(fmt.Errorf("Binding error %s", err))
			return c.JSON(http.StatusBadRequest, tools.ErrorResponce{err.Error()})
		}

		if err := c.Validate(req); err != nil {
			logrus.Error(fmt.Errorf("Validate error %s", err))
			return c.JSON(http.StatusBadRequest, tools.ErrorResponce{err.Error()})
		}

		nickname := c.Param("nickname")

		u :=  &models.User{
			Email:    req.Email,
			Fullname: req.Fullname,
			About:    req.About,
		}

		err := uh.userUC.Update(nickname, u)
		if err != nil {
			if err == tools.ErrExistWithEmail {
				return c.JSON(http.StatusConflict, tools.ErrorResponce{err.Error()})
			}
			if err == tools.ErrDoesntExists {
				return c.JSON(http.StatusConflict, tools.ErrorResponce{err.Error()})
			}
			logrus.Error(fmt.Errorf("Request error %s", err))
			return c.JSON(http.StatusBadRequest, tools.ErrorResponce{err.Error()})
		}



		return c.JSON(http.StatusOK, u)
	}
}
