package auth

import (
	"context"
	"net/http"
)

type PrivilegesHandler struct {
	Privileges func(ctx context.Context) ([]Privilege, error)
	Resource   string
	Action     string
	WriteLog   func(ctx context.Context, resource string, action string, success bool, desc string) error
}

func NewPrivilegesHandler(reader func(context.Context) ([]Privilege, error)) *PrivilegesHandler {
	return NewDefaultPrivilegesHandler(reader, "", "", nil)
}
func NewDefaultPrivilegesHandler(reader func(context.Context) ([]Privilege, error), resource string, action string, writeLog func(context.Context, string, string, bool, string) error) *PrivilegesHandler {
	if len(resource) == 0 {
		resource = "privilege"
	}
	if len(action) == 0 {
		action = "all"
	}
	h := PrivilegesHandler{Privileges: reader, Resource: resource, Action: action, WriteLog: writeLog}
	return &h
}
func (c *PrivilegesHandler) Handle(w http.ResponseWriter, r *http.Request) {
	privileges, err := c.Privileges(r.Context())
	if err != nil {
		respond(w, r, http.StatusInternalServerError, internalServerError, c.WriteLog, c.Resource, c.Action, false, err.Error())
	} else {
		respond(w, r, http.StatusOK, privileges, c.WriteLog, c.Resource, c.Action, true, "")
	}
}
