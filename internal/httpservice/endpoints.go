package httpservice

import (
	"net/http"

	"github.com/flypay/go-kit/v4/pkg/log"
	"github.com/labstack/echo/v4"
)

type HTTPHandler struct {
	Logger log.Logger
}

// EXAMPLES

// func (h HTTPHandler) UserGet(ctx echo.Context, userId string) error {
// 	return echo.NewHTTPError(http.StatusNotFound)
// }

// func (h HTTPHandler) UserAdd(ctx echo.Context) error {
// 	return echo.NewHTTPError(http.StatusNotFound)
// }

func (h HTTPHandler) AddUser(c_tx echo.Context) error {
	return echo.NewHTTPError(http.StatusInternalServerError)
}

// do something here
type UserRecord struct {
	Identifier string // Needs to be sentence case for DynamoDB

	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`

	DateCreated  string `json:"date_created,omitempty"`
	DateModified string `json:"date_modified,omitempty"`
}
