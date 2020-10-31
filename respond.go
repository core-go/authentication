package auth

import (
	"context"
	"encoding/json"
	"net/http"
)
type AuthActivityLogWriter interface {
	Write(ctx context.Context, resource string, action string, success bool, desc string) error
}
func RespondString(w http.ResponseWriter, r *http.Request, code int, result string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write([]byte(result))
}
func Respond(w http.ResponseWriter, r *http.Request, code int, result interface{}, logService AuthActivityLogWriter, resource string, action string, success bool, desc string) {
	response, _ := json.Marshal(result)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
	if logService != nil {
		newCtx := context.WithValue(r.Context(), "request", r)
		logService.Write(newCtx, resource, action, success, desc)
	}
}
