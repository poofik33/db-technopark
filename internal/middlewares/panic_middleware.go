package middlewares

import (
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"net/http"
)

func PanicMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		defer func() {
			if err := recover(); err != nil {
				logrus.Error("Panic error: ", err)
				c.NoContent(http.StatusInternalServerError)
			}
		}()
		return next(c)
	}
}
