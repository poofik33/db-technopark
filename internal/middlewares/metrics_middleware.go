package middlewares

import (
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"time"
)

type MetricsController struct {
	Duration *prometheus.SummaryVec
}

var id int

func NewMetricsController(router *echo.Echo) *MetricsController {
	duration := prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Name: "duration",
		Help: "The latency of the HTTP requests.",
	}, []string{"path", "method"})

	m := &MetricsController{
		Duration: duration,
	}

	prometheus.MustRegister(duration)

	router.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

	router.Use(m.GetTime)

	return m
}

func (mC *MetricsController) GetTime(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		start := time.Now()
		err := next(c)
		end := time.Since(start).Seconds()
		id++
		mC.Duration.WithLabelValues(c.Path(), c.Request().Method).Observe(end)
		return err
	}
}
