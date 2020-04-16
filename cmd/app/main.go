package main

import (
	"database/sql"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/poofik33/db-technopark/database"
	"github.com/poofik33/db-technopark/internal/middlewares"
	"github.com/poofik33/db-technopark/internal/user/delivery"
	"github.com/poofik33/db-technopark/internal/user/repository"
	"github.com/poofik33/db-technopark/internal/user/usecase"
	"github.com/sirupsen/logrus"

	_ "github.com/lib/pq"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func main() {
	connString := fmt.Sprintf("host=%s port=%d dbname=%s user=%s sslmode=disable password=%s",
		"localhost",
		5432,
		"forums",
		"forums_user",
		"difficult_password")


	dbConn, err := sql.Open("postgres", connString)
	if err != nil {
		logrus.Fatal(fmt.Errorf("database open connection err %s", err))
		return
	}

	if err = database.InitDB(dbConn); err != nil {
		logrus.Fatal(fmt.Errorf("database init err %s", err))
		return
	}

	e := echo.New()

	e.Use(middlewares.PanicMiddleware)
	e.Validator = &CustomValidator{validator: validator.New()}

	ur := repository.NewUserRepository(dbConn)
	uc := usecase.NewUserUsecase(ur)
	_ = delivery.NewUserHandler(e, uc)

	e.Start(":5000")
}
