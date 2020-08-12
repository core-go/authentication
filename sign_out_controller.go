package auth

import (
	"net/http"
	"strings"
	"time"
)

type SignOutController struct {
	TokenVerifier         TokenVerifier
	Secret                string
	TokenBlacklistService TokenBlacklistService
	LogService            AuthActivityLogService
}

func NewSignOutController(tokenVerifier TokenVerifier, secret string, tokenBlacklistService TokenBlacklistService, logService AuthActivityLogService) *SignOutController {
	return &SignOutController{tokenVerifier, secret, tokenBlacklistService, logService}
}
func (c *SignOutController) SignOut(w http.ResponseWriter, r *http.Request) {
	data := r.Header["Authorization"]
	token := parseToken(data)

	if len(token) == 0 {
		RespondString(w, r, http.StatusBadRequest, "Invalid token")
		return
	}

	_, _, expiresAt, er1 := c.TokenVerifier.VerifyToken(token, c.Secret)

	if er1 != nil {
		RespondString(w, r, http.StatusBadRequest, "Invalid token")
		return
	}

	if c.TokenBlacklistService == nil {
		RespondString(w, r, http.StatusInternalServerError, "No service to sign out")
		return
	}

	expiresTime := time.Unix(expiresAt, 0)

	er2 := c.TokenBlacklistService.Revoke(token, "The token has signed out.", expiresTime)
	if er2 != nil {
		if c.LogService != nil {
			c.LogService.SaveLog(r.Context(), "Authentication", "Sign Out", false, er2.Error())
		}
		RespondString(w, r, http.StatusInternalServerError, er2.Error())
		return
	}
	if c.LogService != nil {
		c.LogService.SaveLog(r.Context(), "Authentication", "Sign Out", true, "")
	}
	RespondString(w, r, http.StatusOK, "true")
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
