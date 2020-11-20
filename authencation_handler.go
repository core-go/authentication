package auth

import (
	"context"
	"encoding/json"
	"net"
	"net/http"
)

type AuthenticationHandler struct {
	Authenticator         Authenticator
	Resource              string
	Action                string
	LogError              func(context.Context, string)
	Ip                    string
	TokenWhitelistService TokenWhitelistService
	LogWriter             AuthActivityLogWriter
	Decrypter             PasswordDecrypter
	EncryptionKey         string
}

func NewAuthenticationHandlerWithDecrypter(authenticator Authenticator, resource string, action string, logError func(context.Context, string), ip string, tokenWhitelistService TokenWhitelistService,logWriter AuthActivityLogWriter, decrypter PasswordDecrypter, encryptionKey string) *AuthenticationHandler {
	if len(ip) == 0 {
		ip = "ip"
	}
	if len(resource) == 0 {
		resource = "authentication"
	}
	if len(action) == 0 {
		action = "authenticate"
	}
	return &AuthenticationHandler{Authenticator: authenticator, Resource: resource, Action: action, LogError: logError, Ip: ip,TokenWhitelistService: tokenWhitelistService, LogWriter: logWriter, Decrypter: decrypter, EncryptionKey: encryptionKey}
}

func NewAuthenticationHandler(authenticator Authenticator, logError func(context.Context, string), logService AuthActivityLogWriter, tokenWhitelistService TokenWhitelistService) *AuthenticationHandler {
	return NewAuthenticationHandlerWithDecrypter(authenticator, "", "", logError, "", tokenWhitelistService, logService, nil, "")
}

func (h *AuthenticationHandler) Authenticate(w http.ResponseWriter, r *http.Request) {
	ip := GetRemoteIp(r)
	var ctx context.Context
	ctx = r.Context()
	if len(h.Ip) > 0 {
		ctx = context.WithValue(ctx, h.Ip, ip)
		r = r.WithContext(ctx)
	}

	var user AuthInfo
	er1 := json.NewDecoder(r.Body).Decode(&user)
	if er1 != nil {
		if h.LogError != nil {
			msg := "cannot decode authentication info: " + er1.Error()
			h.LogError(r.Context(), msg)
		}
		respondString(w, r, http.StatusBadRequest, "cannot decode authentication info")
		return
	}

	if h.Decrypter != nil && len(h.EncryptionKey) > 0 {
		if decodedPassword, er2 := h.Decrypter.Decrypt(user.Password, h.EncryptionKey); er2 != nil {
			if h.LogError != nil {
				msg := "cannot decrypt password: " + er2.Error()
				h.LogError(r.Context(), msg)
			}
			respondString(w, r, http.StatusBadRequest, "cannot decrypt password")
			return
		} else {
			user.Password = decodedPassword
		}
	}

	result, er3 := h.Authenticator.Authenticate(r.Context(), user)
	if er3 != nil {
		result.Status = StatusSystemError
		if h.LogError != nil {
			h.LogError(r.Context(), er3.Error())
		}
		respond(w, r, http.StatusOK, result, h.LogWriter, h.Resource, h.Action, false, er3.Error())
	} else {
		if h.TokenWhitelistService != nil {
			h.TokenWhitelistService.Add(result.User.UserId, result.User.Token)
		}
		respond(w, r, http.StatusOK, result, h.LogWriter, h.Resource, h.Action, true, "")
	}
}
func GetRemoteIp(r *http.Request) string {
	remoteIP, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		remoteIP = r.RemoteAddr
	}
	return remoteIP
}
