package service

import (
	"github.com/flypay/bootcamp-zoe-flower-users-api/internal/httpservice"
	"github.com/flypay/go-kit/v4/pkg/eventbus"
	"github.com/flypay/go-kit/v4/pkg/log"
	"github.com/labstack/echo/v4"
)

// service requires a http endpoint to communicate
// emits an event so needs a producer
// producer and logger is passed to the handler to use
func RunService(
	http *echo.Echo,
	producer eventbus.Producer,
) error {
	httpservice.RegisterHandlers(http,
		httpservice.HTTPHandler{Logger: log.DefaultLogger, Producer: producer},
	)
	return nil
}
