package echo

import (
	"context"
	a "github.com/core-go/authentication"
	"github.com/labstack/echo"
	"net/http"
)

type PrivilegesHandler struct {
	Load     func(ctx context.Context) ([]a.Privilege, error)
	Error    func(context.Context, string, ...map[string]interface{})
	Log      func(ctx context.Context, resource string, action string, success bool, desc string) error
	Resource string
	Action   string
}

func NewPrivilegesHandler(load func(context.Context) ([]a.Privilege, error), options ...func(context.Context, string, ...map[string]interface{})) *PrivilegesHandler {
	var logError func(context.Context, string, ...map[string]interface{})
	if len(options) >= 1 {
		logError = options[0]
	}
	return NewPrivilegesHandlerWithLog(load, logError, nil)
}
func NewPrivilegesHandlerWithLog(load func(context.Context) ([]a.Privilege, error), logError func(context.Context, string, ...map[string]interface{}), writeLog func(context.Context, string, string, bool, string) error, options ...string) *PrivilegesHandler {
	var resource, action string
	if len(options) >= 1 {
		resource = options[0]
	} else {
		resource = "privilege"
	}
	if len(options) >= 2 {
		action = options[1]
	} else {
		action = "all"
	}
	h := PrivilegesHandler{Load: load, Error: logError, Resource: resource, Action: action, Log: writeLog}
	return &h
}
func (c *PrivilegesHandler) All(ctx echo.Context) error {
	r := ctx.Request()
	privileges, err := c.Load(r.Context())
	if err != nil {
		if c.Error != nil {
			c.Error(r.Context(), err.Error())
		}
		return respond(ctx, http.StatusInternalServerError, internalServerError, c.Log, c.Resource, c.Action, false, err.Error())
	} else {
		return respond(ctx, http.StatusOK, privileges, c.Log, c.Resource, c.Action, true, "")
	}
}
