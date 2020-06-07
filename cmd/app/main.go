package main

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx"
	"github.com/labstack/echo/v4"
	"github.com/poofik33/db-technopark/database"
	forum_delivery "github.com/poofik33/db-technopark/internal/forum/delivery"
	forum_repository "github.com/poofik33/db-technopark/internal/forum/repository"
	forum_usecase "github.com/poofik33/db-technopark/internal/forum/usecase"
	"github.com/poofik33/db-technopark/internal/middlewares"
	post_delivery "github.com/poofik33/db-technopark/internal/post/delivery"
	post_repository "github.com/poofik33/db-technopark/internal/post/repository"
	post_usecase "github.com/poofik33/db-technopark/internal/post/usecase"
	service_delivery "github.com/poofik33/db-technopark/internal/service/delivery"
	service_repository "github.com/poofik33/db-technopark/internal/service/repostitory"
	service_usecase "github.com/poofik33/db-technopark/internal/service/usecase"
	thread_delivery "github.com/poofik33/db-technopark/internal/thread/delivery"
	thread_repository "github.com/poofik33/db-technopark/internal/thread/repository"
	thread_usecase "github.com/poofik33/db-technopark/internal/thread/usecase"
	user_delivery "github.com/poofik33/db-technopark/internal/user/delivery"
	user_repository "github.com/poofik33/db-technopark/internal/user/repository"
	user_usecase "github.com/poofik33/db-technopark/internal/user/usecase"
	vote_repository "github.com/poofik33/db-technopark/internal/vote/repository"
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
	dbConn, err := pgx.NewConnPool(pgx.ConnPoolConfig{
		ConnConfig: pgx.ConnConfig{
			Host:     "localhost",
			Port:     5432,
			Database: "forums",
			User:     "forums_user",
			Password: "difficult_password",
		},
		MaxConnections: 10000,
	})

	if err = database.InitDB(dbConn); err != nil {
		logrus.Fatal(fmt.Errorf("database init err %s", err))
		return
	}

	e := echo.New()

	e.Use(middlewares.PanicMiddleware)
	_ = middlewares.NewMetricsController(e)

	e.Validator = &CustomValidator{validator: validator.New()}

	ur := user_repository.NewUserRepository(dbConn)
	fr := forum_repository.NewForumRepository(dbConn)
	tr := thread_repository.NewThreadRepository(dbConn)
	sr := service_repository.NewServiceRepository(dbConn)
	vr := vote_repository.NewVoteRepository(dbConn)
	pr := post_repository.NewPostRepository(dbConn)

	uUC := user_usecase.NewUserUsecase(ur)
	fUC := forum_usecase.NewForumUsecase(fr, ur, pr, tr)
	tUC := thread_usecase.NewThreadUsecase(tr, ur, fr, pr, vr)
	sUC := service_usecase.NewServiceUsecase(sr)
	pUC := post_usecase.NewPostUsecase(pr, fr, vr, tr, ur)

	_ = user_delivery.NewUserHandler(e, uUC)
	_ = forum_delivery.NewForumHandler(e, fUC, tUC)
	_ = thread_delivery.NewThreadHandler(e, tUC)
	_ = service_delivery.NewServiceHandler(e, sUC)
	_ = post_delivery.NewPostHandler(e, pUC)

	e.Start(":5000")
}
