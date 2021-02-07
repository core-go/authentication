package auth

import (
	"context"
	"net/http"
	"strings"
)

type PrivilegesByEntityHandler struct {
	Privileges func(ctx context.Context, id string) ([]Privilege, error)
	Resource   string
	Action     string
	Offset     int
	WriteLog   func(ctx context.Context, resource string, action string, success bool, desc string) error
}
func NewPrivilegesByEntityHandler(loader func(ctx context.Context, id string) ([]Privilege, error)) *PrivilegesByEntityHandler {
	return NewDefaultPrivilegesByEntityHandler(loader, "", "", 0, nil)
}
func NewDefaultPrivilegesByEntityHandler(loader func(ctx context.Context, id string) ([]Privilege, error), resource string, action string, offset int, writeLog func(context.Context, string, string, bool, string) error) *PrivilegesByEntityHandler {
	if len(resource) == 0 {
		resource = "privilege"
	}
	if len(action) == 0 {
		action = "all"
	}
	h := PrivilegesByEntityHandler{Privileges: loader, Resource: resource, Action: action, Offset: offset, WriteLog: writeLog}
	return &h
}
func (c *PrivilegesByEntityHandler) PrivilegesById(w http.ResponseWriter, r *http.Request) {
	id := ""
	if c.Offset <=0 {
		i := strings.LastIndex(r.RequestURI, "/")
		if i >= 0 {
			id = r.RequestURI[i + 1:]
		}
	} else {
		s := strings.Split(r.RequestURI, "/")
		if len(s) - c.Offset - 1 >= 0 {
			id = s[len(s) - c.Offset - 1]
		} else {
			http.Error(w, "URL is not valid", http.StatusBadRequest)
			return
		}
	}
	privileges, err := c.Privileges(r.Context(), id)
	if err != nil {
		respond(w, r, http.StatusInternalServerError, internalServerError, c.WriteLog, c.Resource, c.Action, false, err.Error())
	} else {
		respond(w, r, http.StatusOK, privileges, c.WriteLog, c.Resource, c.Action, true, "")
	}
}
