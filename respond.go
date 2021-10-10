package auth

import (
	"context"
	"encoding/json"
	"net/http"
)

const internalServerError = "Internal Server Error"

func respond(w http.ResponseWriter, r *http.Request, code int, result interface{}, writeLog func(context.Context, string, string, bool, string) error, resource string, action string, success bool, desc string) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	err := json.NewEncoder(w).Encode(result)
	if writeLog != nil {
		writeLog(r.Context(), resource, action, success, desc)
	}
	return err
}
