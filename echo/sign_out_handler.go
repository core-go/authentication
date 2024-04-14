package echo

import (
	"context"
	"github.com/labstack/echo/v4"
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
	Cookie      bool
	CookieName  string
	CookieDomain string
}

func NewSignOutHandler(verifyToken func(tokenString string, secret string) (map[string]interface{}, int64, int64, error), secret string, revokeToken func(ctx context.Context, token string, reason string, expires time.Time) error, logError func(context.Context, string, ...map[string]interface{}), options...func(context.Context, string, string, bool, string) error) *SignOutHandler {
	var writeLog func(context.Context, string, string, bool, string) error
	if len(options) >= 1 {
		writeLog = options[0]
	}
	return NewSignOutHandlerWithLog(verifyToken, secret, revokeToken, logError, writeLog, false, "", "id", "authentication", "signout")
}
func NewSignOutHandlerWithLog(verifyToken func(tokenString string, secret string) (map[string]interface{}, int64, int64, error), secret string, revokeToken func(ctx context.Context,token string, reason string, expires time.Time) error, logError func(context.Context, string, ...map[string]interface{}), writeLog func(context.Context, string, string, bool, string) error, cookie bool, options ...string) *SignOutHandler {
	var cookieName, cookieDomain, resource, action string
	if len(options) > 0 {
		cookieDomain = options[0]
	}
	if len(options) > 1 {
		cookieName = options[1]
	} else {
		cookieName = "id"
	}
	if len(options) > 2 {
		resource = options[2]
	} else {
		resource = "authentication"
	}
	if len(options) > 3 {
		action = options[3]
	} else {
		action = "signout"
	}
	return &SignOutHandler{VerifyToken: verifyToken, Secret: secret, Cookie: cookie, CookieName: cookieName, CookieDomain: cookieDomain, RevokeToken: revokeToken, Error: logError, Log: writeLog, Resource: resource, Action: action}
}
func (h *SignOutHandler) SignOutCookie(ctx echo.Context) error {
	if _, err := getCookieValueByID(ctx.Request(), h.CookieName); err == nil {
		cookie := &http.Cookie{
			Name: h.CookieName,
			Domain: h.CookieDomain,
			Value: "",
			HttpOnly: true,
			Path: "/",
			MaxAge: -1,
			SameSite: http.SameSiteStrictMode,
			Secure: true,
		}
		ctx.SetCookie(cookie)
	}
	return respond(ctx, http.StatusOK, true, h.Log, h.Resource, h.Action, true, "")
}
func (h *SignOutHandler) SignOut(ctx echo.Context) error {
	if h.Cookie {
		return h.SignOutCookie(ctx)
	}
	data := ctx.Request().Header["Authorization"]
	token := parseToken(data)

	if len(token) == 0 {
		return ctx.String(http.StatusBadRequest, "Invalid token")
	}

	_, _, expiresAt, er1 := h.VerifyToken(token, h.Secret)

	if er1 != nil {
		return ctx.String(http.StatusBadRequest, "Invalid token")
	}

	if h.RevokeToken == nil {
		return ctx.String(http.StatusServiceUnavailable, "No service to sign out")
	}

	expiresTime := time.Unix(expiresAt, 0)

	er2 := h.RevokeToken(ctx.Request().Context(), token, "The token has signed out.", expiresTime)
	if er2 != nil {
		if h.Error != nil {
			h.Error(ctx.Request().Context(), er2.Error())
		}
		if h.Log != nil {
			h.Log(ctx.Request().Context(), h.Resource, h.Action, false, er2.Error())
		}
		return ctx.String(http.StatusInternalServerError, internalServerError)
	}
	if h.Log != nil {
		h.Log(ctx.Request().Context(), h.Resource, h.Action, true, "")
	}
	return respond(ctx, http.StatusOK, true, h.Log, h.Resource, h.Action, true, "")
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

func getCookieValueByID(r *http.Request, cookieName string) (string, error) {
	cookie, err := r.Cookie(cookieName)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}
