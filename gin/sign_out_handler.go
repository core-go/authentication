package gin

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

type SignOutHandler struct {
	VerifyToken func(tokenString string, secret string) (map[string]interface{}, int64, int64, error)
	Secret      string
	RevokeToken func(ctx context.Context, token string, reason string, expires time.Time) error
	Error       func(context.Context, string, ...map[string]interface{})
	Log         func(ctx context.Context, resource string, action string, success bool, desc string) error
	Resource    string
	Action      string
}

func NewSignOutHandler(verifyToken func(tokenString string, secret string) (map[string]interface{}, int64, int64, error), secret string, revokeToken func(ctx context.Context, token string, reason string, expires time.Time) error, logError func(context.Context, string, ...map[string]interface{}), options...func(context.Context, string, string, bool, string) error) *SignOutHandler {
	var writeLog func(context.Context, string, string, bool, string) error
	if len(options) >= 1 {
		writeLog = options[0]
	}
	return NewSignOutHandlerWithLog(verifyToken, secret, revokeToken, logError, writeLog, "authentication", "signout")
}
func NewSignOutHandlerWithLog(verifyToken func(tokenString string, secret string) (map[string]interface{}, int64, int64, error), secret string, revokeToken func(ctx context.Context, token string, reason string, expires time.Time) error, logError func(context.Context, string, ...map[string]interface{}), writeLog func(context.Context, string, string, bool, string) error, options ...string) *SignOutHandler {
	var resource, action string
	if len(options) >= 1 {
		resource = options[0]
	} else {
		resource = "authentication"
	}
	if len(options) >= 2 {
		action = options[1]
	} else {
		action = "signout"
	}
	return &SignOutHandler{VerifyToken: verifyToken, Secret: secret, RevokeToken: revokeToken, Error: logError, Log: writeLog, Resource: resource, Action: action}
}
func (h *SignOutHandler) SignOut(ctx *gin.Context) {
	r := ctx.Request
	data := r.Header["Authorization"]
	token := parseToken(data)

	if len(token) == 0 {
		ctx.String(http.StatusBadRequest, "Invalid token")
		return
	}

	_, _, expiresAt, er1 := h.VerifyToken(token, h.Secret)

	if er1 != nil {
		ctx.String(http.StatusBadRequest, "Invalid token")
		return
	}

	if h.RevokeToken == nil {
		ctx.String(http.StatusServiceUnavailable, "No service to sign out")
		return
	}

	expiresTime := time.Unix(expiresAt, 0)

	er2 := h.RevokeToken(r.Context(), token, "The token has signed out.", expiresTime)
	if er2 != nil {
		if h.Error != nil {
			h.Error(r.Context(), er2.Error())
		}
		if h.Log != nil {
			h.Log(r.Context(), h.Resource, h.Action, false, er2.Error())
		}
		ctx.String(http.StatusInternalServerError, internalServerError)
		return
	}
	if h.Log != nil {
		h.Log(r.Context(), h.Resource, h.Action, true, "")
	}
	respond(ctx, http.StatusOK, true, h.Log, h.Resource, h.Action, true, "")
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
