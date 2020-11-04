package auth

import (
	"net/http"
	"strings"
)

type PrivilegesByEntityHandler struct {
	Loader    PrivilegesLoader
	Resource  string
	Action    string
	Offset    int
	LogWriter AuthActivityLogWriter
}

func NewPrivilegesByEntityHandler(loader PrivilegesLoader, resource string, action string, offset int, logWriter AuthActivityLogWriter) *PrivilegesByEntityHandler {
	if len(resource) == 0 {
		resource = "Privileges"
	}
	if len(action) == 0 {
		action = "All"
	}
	h := PrivilegesByEntityHandler{Loader: loader, Resource: resource, Action: action, Offset: offset, LogWriter: logWriter}
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
			RespondString(w, r, http.StatusBadRequest, "URL is not valid")
			return
		}
	}
	privileges, err := c.Loader.Load(r.Context(), id)
	if err != nil {
		Respond(w, r, http.StatusInternalServerError, InternalServerError, c.LogWriter, c.Resource, c.Action, false, err.Error())
	} else {
		Respond(w, r, http.StatusOK, privileges, c.LogWriter, c.Resource, c.Action, true, "")
	}
}
