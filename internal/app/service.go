package service

import (
	"github.com/flypay/bootcamp-zoe_flower-users-api/internal/httpservice"
	"github.com/flypay/go-kit/v4/pkg/log"
	"github.com/labstack/echo/v4"
)

func RunService(
	http *echo.Echo,
) error {
	// logger := log.DefaultLogger
	// logger.Debugf("Listening on port %v", http.Server)
	httpservice.RegisterHandlers(http, httpservice.HTTPHandler{
		Logger: log.DefaultLogger,
	})
	return nil
}
