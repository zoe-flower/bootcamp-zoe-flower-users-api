package httpservice

import (
	"net/http"

	"github.com/flypay/events/pkg/bootcamp"
	"github.com/flypay/go-kit/v4/pkg/eventbus"
	"github.com/flypay/go-kit/v4/pkg/log"
	"github.com/labstack/echo/v4"
)

type HTTPHandler struct {
	Logger   log.Logger
	Producer eventbus.Producer
}

// the binding takes the queryparams from curl command and binds to context? (CHECK)
// we create a user object to emit data in the correct format.
// why is there both User and UserCreated? (CHECK).
func (h HTTPHandler) AddUser(ctx echo.Context) error {
	var u User
	if err := ctx.Bind(&u); err != nil {
		h.Logger.Errorf("failed to bind user create request: %s", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	h.Logger.Infof("received user create request %v", u)
	if u.FirstName == "" {
		return echo.NewHTTPError(http.StatusBadRequest)
	}
	user := bootcamp.UserCreated{
		Id:          u.UserId,
		FirstName:   u.FirstName,
		LastName:    u.LastName,
		DateOfBirth: u.Dob,
		SlackHandle: u.SlackHandle,
	}
	err := h.Producer.Emit(ctx.Request().Context(), &user)
	if err != nil {
		h.Logger.Errorf("failed to bind user create request: %s", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return echo.NewHTTPError(http.StatusOK)
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
