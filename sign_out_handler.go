package auth

import (
	"context"
	"net/http"
	"strings"
	"time"
)

type SignOutHandler struct {
	VerifyToken           func(tokenString string, secret string) (map[string]interface{}, int64, int64, error)
	Resource              string
	Action                string
	Secret                string
	RevokeToken           func(token string, reason string, expires time.Time) error
	WriteLog              func(ctx context.Context, resource string, action string, success bool, desc string) error
}

func NewDefaultSignOutHandler(verifyToken func(tokenString string, secret string) (map[string]interface{}, int64, int64, error), secret string, revokeToken func(token string, reason string, expires time.Time) error, writeLog func(context.Context, string, string, bool, string) error) *SignOutHandler {
	return NewSignOutHandler(verifyToken, "authentication", "signout", secret, revokeToken, writeLog)
}
func NewSignOutHandler(verifyToken func(tokenString string, secret string) (map[string]interface{}, int64, int64, error), resource string, action string, secret string, revokeToken func(token string, reason string, expires time.Time) error, writeLog func(context.Context, string, string, bool, string) error) *SignOutHandler {
	if len(resource) == 0 {
		resource = "authentication"
	}
	if len(action) == 0 {
		action = "signout"
	}
	return &SignOutHandler{VerifyToken: verifyToken, Resource: resource, Action: action, Secret: secret, RevokeToken: revokeToken, WriteLog: writeLog}
}
func (h *SignOutHandler) SignOut(w http.ResponseWriter, r *http.Request) {
	data := r.Header["Authorization"]
	token := parseToken(data)

	if len(token) == 0 {
		http.Error(w, "Invalid token", http.StatusBadRequest)
		return
	}

	_, _, expiresAt, er1 := h.VerifyToken(token, h.Secret)

	if er1 != nil {
		http.Error(w, "Invalid token", http.StatusBadRequest)
		return
	}

	if h.RevokeToken == nil {
		http.Error(w, "No service to sign out", http.StatusServiceUnavailable)
		return
	}

	expiresTime := time.Unix(expiresAt, 0)

	er2 := h.RevokeToken(token, "The token has signed out.", expiresTime)
	if er2 != nil {
		if h.WriteLog != nil {
			h.WriteLog(r.Context(), h.Resource, h.Action, false, er2.Error())
		}
		http.Error(w, internalServerError, http.StatusInternalServerError)
		return
	}
	if h.WriteLog != nil {
		h.WriteLog(r.Context(), h.Resource, h.Action, true, "")
	}
	respond(w, r, http.StatusOK, true, h.WriteLog, h.Resource, h.Action, true, "")
}

func parseToken(data []string) string {
	if len(data) == 0 {
		return ""
	}
	authorization := data[0]
	if strings.HasPrefix(authorization, "Bearer ") != true {
		return ""
	}
	return authorization[7:]
}
