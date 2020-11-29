package auth

import (
	"net/http"
	"strings"
	"time"
)

type SignOutHandler struct {
	TokenVerifier         TokenVerifier
	Resource              string
	Action                string
	Secret                string
	TokenBlacklistChecker TokenBlacklistChecker
	LogWriter             AuthActivityLogWriter
}
func NewDefaultSignOutHandler(tokenVerifier TokenVerifier, secret string, tokenBlacklistService TokenBlacklistChecker, logWriter AuthActivityLogWriter) *SignOutHandler {
	return NewSignOutHandler(tokenVerifier, "", "", secret, tokenBlacklistService, logWriter)
}
func NewSignOutHandler(tokenVerifier TokenVerifier, resource string, action string, secret string, tokenBlacklistService TokenBlacklistChecker, logWriter AuthActivityLogWriter) *SignOutHandler {
	if len(resource) == 0 {
		resource = "authentication"
	}
	if len(action) == 0 {
		action = "signout"
	}
	return &SignOutHandler{TokenVerifier: tokenVerifier, Resource: resource, Action: action, Secret: secret, TokenBlacklistChecker: tokenBlacklistService, LogWriter: logWriter}
}
func (h *SignOutHandler) SignOut(w http.ResponseWriter, r *http.Request) {
	data := r.Header["Authorization"]
	token := parseToken(data)

	if len(token) == 0 {
		http.Error(w, "Invalid token", http.StatusBadRequest)
		return
	}

	_, _, expiresAt, er1 := h.TokenVerifier.VerifyToken(token, h.Secret)

	if er1 != nil {
		http.Error(w, "Invalid token", http.StatusBadRequest)
		return
	}

	if h.TokenBlacklistChecker == nil {
		http.Error(w, "No service to sign out", http.StatusServiceUnavailable)
		return
	}

	expiresTime := time.Unix(expiresAt, 0)

	er2 := h.TokenBlacklistChecker.Revoke(token, "The token has signed out.", expiresTime)
	if er2 != nil {
		if h.LogWriter != nil {
			h.LogWriter.Write(r.Context(), h.Resource, h.Action, false, er2.Error())
		}
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}
	if h.LogWriter != nil {
		h.LogWriter.Write(r.Context(), h.Resource, h.Action, true, "")
	}
	respond(w, r, http.StatusOK, true, h.LogWriter, h.Resource, h.Action, true, "")
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
