package httpservice

import (
	"fmt"
	"net/http"

	"github.com/flypay/go-kit/v4/pkg/log"
	"github.com/labstack/echo/v4"
)

type HTTPHandler struct {
	Logger log.Logger
}

// EXAMPLES
// echo.NewHTTPError(http.StatusInternalServerError)

// Add User is the http Handler

// h.Logger.Debugf("add user called with %s", ctx.Request().Method)
// h.Logger.Printf("hello", ctx.Request().URL.Query().Get("user_id"))
// h.Logger.Printf("hello", ctx.Bind(user))

// dob := user.Dob
// firstName := user.FirstName
// lastName := user.LastName
// slackHandle := user.SlackHandle
// userID := user.UserId

// fmt.Println("DOB:", dob)
// fmt.Println("First Name:", firstName)
// fmt.Println("Last Name:", lastName)
// fmt.Println("Slack Handle:", slackHandle)
// fmt.Println("UserID:", userID)

// logger := log.DefaultLogger
// logger.Debugf("Storing query results: %v", x)
// fmt.Print("USER", ctx.ParamValues())

func (h HTTPHandler) AddUser(ctx echo.Context) error {
	// var user User
	var user User
	err := ctx.Bind(&user)
	if err != nil {
		print(err, "ERROR!!")
		fmt.Println("Error parsing JSON:", err)
		return ctx.String(http.StatusBadRequest, "bad request")
	}
	fmt.Printf("%v", user)
	return ctx.JSON(http.StatusOK, user)
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
