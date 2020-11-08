package auth

import (
	"context"
	"encoding/json"
	"net/http"
)
const InternalServerError = "Internal Server Error"
type AuthActivityLogWriter interface {
	Write(ctx context.Context, resource string, action string, success bool, desc string) error
}
func respondString(w http.ResponseWriter, r *http.Request, code int, result string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write([]byte(result))
}
func respond(w http.ResponseWriter, r *http.Request, code int, result interface{}, logWriter AuthActivityLogWriter, resource string, action string, success bool, desc string) {
	response, _ := json.Marshal(result)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
	if logWriter != nil {
		logWriter.Write(r.Context(), resource, action, success, desc)
	}
}
