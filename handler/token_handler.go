package handler

import (
	"context"
	"net/http"
	"net/url"
	"strings"
	"time"

	a "github.com/core-go/authentication"
)

type TokenHandler struct {
	Host                string
	SameSite            http.SameSite
	CookieName          string
	RememberCookieName  string
	TokenConfig         a.TokenConfig
	RememberTokenConfig a.TokenConfig
	Error               func(context.Context, string, ...map[string]interface{})
	GenerateToken       func(payload interface{}, secret string, expiresIn int64) (string, error)
	GetAndVerifyToken   func(authorization string, secret string) (bool, string, map[string]interface{}, int64, int64, error)
	Resource            string
	Log                 func(ctx context.Context, resource string, action string, success bool, desc string) error
}

func (h *TokenHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	if h.GetAndVerifyToken == nil {
		http.Error(w, "refresh token is not supported", http.StatusNotImplemented)
		return
	}

	rememberCookie, err := r.Cookie(h.RememberCookieName)
	if err != nil || rememberCookie == nil || len(rememberCookie.Value) == 0 {
		http.Error(w, h.RememberCookieName+" is required in cookies", http.StatusUnauthorized)
		return
	}

	isToken, _, data, _, _, err := h.GetAndVerifyToken(
		rememberCookie.Value,
		h.RememberTokenConfig.Secret,
	)
	if !isToken || err != nil {
		http.Error(w, "invalid remember token", http.StatusUnauthorized)
		return
	}

	newToken, err := h.GenerateToken(data, h.TokenConfig.Secret, h.TokenConfig.Expires)
	if err != nil {
		if h.Error != nil {
			h.Error(r.Context(), err.Error())
		}
		http.Error(w, "failed to generate token", http.StatusInternalServerError)
		return
	}

	host := r.Header.Get("Origin")
	if strings.Contains(host, h.Host) || strings.Contains(host, "localhost") {
		u, parseErr := url.Parse(host)
		if parseErr == nil {
			host = strings.TrimPrefix(u.Hostname(), "www.")
		}
	}

	http.SetCookie(w, &http.Cookie{
		Name:     h.CookieName,
		Domain:   host,
		Value:    newToken,
		HttpOnly: true,
		Path:     "/",
		MaxAge:   0,
		Expires:  time.Now().Add(time.Duration(h.TokenConfig.Expires) * time.Millisecond),
		SameSite: h.SameSite,
		Secure:   true,
	})

	respond(w, r, http.StatusOK, 1, h.Log, h.Resource, "refresh_token", true, "")
}
