package auth

import (
	"context"
	"encoding/json"
	"net"
	"net/http"
)

type AuthenticationHandler struct {
	Authenticator Authenticator
	Resource      string
	Action        string
	LogError      func(context.Context, string)
	Ip            string
	LogWriter     AuthActivityLogWriter
	Decrypter     PasswordDecrypter
	EncryptionKey string
}

func NewAuthenticationHandlerWithDecrypter(authenticator Authenticator, resource string, action string, logError func(context.Context, string), ip string, logService AuthActivityLogWriter, decrypter PasswordDecrypter, encryptionKey string) *AuthenticationHandler {
	if len(ip) == 0 {
		ip = "ip"
	}
	if len(resource) == 0 {
		resource = "authentication"
	}
	if len(action) == 0 {
		action = "signin"
	}
	return &AuthenticationHandler{Authenticator: authenticator, Resource: resource, Action: action, LogError: logError, Ip: ip, LogWriter: logService, Decrypter: decrypter, EncryptionKey: encryptionKey}
}

func NewAuthenticationHandler(authenticator Authenticator, logError func(context.Context, string), logService AuthActivityLogWriter) *AuthenticationHandler {
	return NewAuthenticationHandlerWithDecrypter(authenticator, "", "", logError, "", logService, nil, "")
}

func (h *AuthenticationHandler) Authenticate(w http.ResponseWriter, r *http.Request) {
	ip := GetRemoteIp(r)
	var ctx context.Context
	ctx = r.Context()
	if len(h.Ip) > 0 {
		ctx = context.WithValue(ctx, h.Ip, ip)
	}

	var user AuthInfo
	er1 := json.NewDecoder(r.Body).Decode(&user)
	if er1 != nil {
		RespondString(w, r, http.StatusBadRequest, "cannot decode authentication info: "+er1.Error())
		return
	}

	if h.Decrypter != nil && len(h.EncryptionKey) > 0 {
		if decodedPassword, er2 := h.Decrypter.Decrypt(user.Password, h.EncryptionKey); er2 != nil {
			RespondString(w, r, http.StatusBadRequest, "cannot decrypt password: "+er2.Error())
			return
		} else {
			user.Password = decodedPassword
		}
	}

	result, er3 := h.Authenticator.Authenticate(ctx, user)
	if er3 != nil {
		result.Status = StatusSystemError
		if h.LogError != nil {
			h.LogError(r.Context(), er3.Error())
		}
		Respond(w, r, http.StatusOK, result, h.LogWriter, h.Resource, h.Action, false, er3.Error())
	} else {
		Respond(w, r, http.StatusOK, result, h.LogWriter, h.Resource, h.Action, true, "")
	}
}
func GetRemoteIp(r *http.Request) string {
	remoteIP, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		remoteIP = r.RemoteAddr
	}
	return remoteIP
}
