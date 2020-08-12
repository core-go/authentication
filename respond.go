package auth

import (
	"context"
	"encoding/json"
	"net/http"
)

func RespondString(w http.ResponseWriter, r *http.Request, code int, result string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write([]byte(result))
}
func Respond(w http.ResponseWriter, r *http.Request, code int, result interface{}, logService AuthActivityLogService, resource string, action string, success bool, desc string) {
	response, _ := json.Marshal(result)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
	if logService != nil {
		newCtx := context.WithValue(r.Context(), "request", r)
		logService.SaveLog(newCtx, resource, action, success, desc)
	}
}
