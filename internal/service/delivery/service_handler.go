package delivery

import (
	"github.com/labstack/echo/v4"
	"github.com/poofik33/db-technopark/internal/service"
	"github.com/poofik33/db-technopark/internal/tools"
	"net/http"
)

type ServiceHandler struct {
	serviceUsecase service.Usecase
}

func NewServiceHandler(e *echo.Echo, sUC service.Usecase) *ServiceHandler {
	sh := &ServiceHandler{
		serviceUsecase: sUC,
	}

	e.GET("/service/status", sh.GetStatus())
	e.POST("/service/clear", sh.Clear())

	return sh
}

func (sh *ServiceHandler) GetStatus() echo.HandlerFunc {
	return func(c echo.Context) error {
		stat, err := sh.serviceUsecase.GetStatus()
		if err != nil {
			return c.JSON(http.StatusBadRequest, tools.ErrorResponce{
				Message: err.Error(),
			})
		}
		return c.JSON(http.StatusOK, stat)
	}
}

func (sh *ServiceHandler) Clear() echo.HandlerFunc {
	return func(c echo.Context) error {
		err := sh.serviceUsecase.DeleteAll()
		if err != nil {
			return c.JSON(http.StatusBadRequest, tools.ErrorResponce{
				Message: err.Error(),
			})
		}
		return c.NoContent(http.StatusOK)
	}
}
