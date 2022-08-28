package handler

import (
	"context"
	"net/http"
	"strings"

	. "github.com/core-go/auth"
)

type PrivilegesByEntityHandler struct {
	Load     func(ctx context.Context, id string) ([]Privilege, error)
	Error    func(context.Context, string, ...map[string]interface{})
	Offset   int
	Log      func(ctx context.Context, resource string, action string, success bool, desc string) error
	Resource string
	Action   string
}

func NewPrivilegesByEntityHandler(load func(ctx context.Context, id string) ([]Privilege, error), options ...func(context.Context, string, ...map[string]interface{})) *PrivilegesByEntityHandler {
	var logError func(context.Context, string, ...map[string]interface{})
	if len(options) >= 1 {
		logError = options[0]
	}
	return NewPrivilegesByEntityHandlerWithLog(load, logError, 0, nil)
}
func NewPrivilegesByEntityHandlerWithLog(load func(ctx context.Context, id string) ([]Privilege, error), logError func(context.Context, string, ...map[string]interface{}), offset int, writeLog func(context.Context, string, string, bool, string) error, options ...string) *PrivilegesByEntityHandler {
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
	h := PrivilegesByEntityHandler{Load: load, Error: logError, Resource: resource, Action: action, Offset: offset, Log: writeLog}
	return &h
}
func (c *PrivilegesByEntityHandler) Privileges(w http.ResponseWriter, r *http.Request) {
	id := ""
	url := r.URL.Path
	if c.Offset <= 0 {
		i := strings.LastIndex(url, "/")
		if i >= 0 {
			id = url[i+1:]
		}
	} else {
		s := strings.Split(url, "/")
		if len(s)-c.Offset-1 >= 0 {
			id = s[len(s)-c.Offset-1]
		} else {
			http.Error(w, "URL is not valid", http.StatusBadRequest)
			return
		}
	}
	privileges, err := c.Load(r.Context(), id)
	if err != nil {
		if c.Error != nil {
			c.Error(r.Context(), err.Error())
		}
		respond(w, r, http.StatusInternalServerError, internalServerError, c.Log, c.Resource, c.Action, false, err.Error())
	} else {
		respond(w, r, http.StatusOK, privileges, c.Log, c.Resource, c.Action, true, "")
	}
}
