package httpservice

import (
	"net/http"

	"github.com/flypay/go-kit/v4/pkg/eventbus"
	"github.com/flypay/go-kit/v4/pkg/log"
	"github.com/labstack/echo/v4"
)

type HTTPHandler struct {
	Logger   log.Logger
	Producer eventbus.Producer
}

// EXAMPLES
// echo.NewHTTPError(http.StatusInternalServerError)

// Add User is the http Handler

func (h HTTPHandler) AddUser(ctx echo.Context) error {
	var req User
	// var stopErr protoreflect.ProtoMessage
	// var stopErr2 context.Context
	if err := ctx.Bind(&req); err != nil {
		h.Logger.Errorf("failed to bind user create request: %s", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	h.Logger.Infof("received user create request %v", req)
	// h.Logger.Infof("CONTEXT: %v", ctx)
	// h.Producer.Emit(stopErr2, stopErr)
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
