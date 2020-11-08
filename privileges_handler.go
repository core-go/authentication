package auth

import "net/http"

type PrivilegesHandler struct {
	Reader    PrivilegesReader
	Resource  string
	Action    string
	LogWriter AuthActivityLogWriter
}
func NewPrivilegesHandler(reader PrivilegesReader, resource string, action string, logWriter AuthActivityLogWriter) *PrivilegesHandler {
	if len(resource) == 0 {
		resource = "Privileges"
	}
	if len(action) == 0 {
		action = "All"
	}
	h := PrivilegesHandler{Reader: reader, Resource: resource, Action: action, LogWriter: logWriter}
	return &h
}
func (c *PrivilegesHandler) Privileges(w http.ResponseWriter, r *http.Request) {
	privileges, err := c.Reader.Privileges(r.Context())
	if err != nil {
		respond(w, r, http.StatusInternalServerError, InternalServerError, c.LogWriter, c.Resource, c.Action, false, err.Error())
	} else {
		respond(w, r, http.StatusOK, privileges, c.LogWriter, c.Resource, c.Action, true, "")
	}
}
