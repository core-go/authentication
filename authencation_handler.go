package auth

import (
	"context"
	"encoding/json"
	"net"
	"net/http"
)

type AuthenticationHandler struct {
	Authenticator Authenticator
	Decrypter     PasswordDecrypter
	EncryptionKey string
	LogWriter    AuthActivityLogWriter
	Ip            string
}

func NewAuthenticationHandler(authenticationService Authenticator, decrypter PasswordDecrypter, encryptionKey string, logService AuthActivityLogWriter, ip string) *AuthenticationHandler {
	return &AuthenticationHandler{authenticationService, decrypter, encryptionKey, logService, ip}
}

func NewDefaultAuthenticationHandler(authenticationService Authenticator) *AuthenticationHandler {
	return &AuthenticationHandler{authenticationService, nil, "", nil, ""}
}

func (c *AuthenticationHandler) Authenticate(w http.ResponseWriter, r *http.Request) {
	ip := GetRemoteIp(r)
	var ctx context.Context
	ctx = r.Context()
	if len(c.Ip) > 0 {
		ctx = context.WithValue(ctx, c.Ip, ip)
	}

	var user AuthInfo
	er1 := json.NewDecoder(r.Body).Decode(&user)
	if er1 != nil {
		RespondString(w, r, http.StatusBadRequest, "cannot decode authentication info: "+er1.Error())
		return
	}

	if c.Decrypter != nil && len(c.EncryptionKey) > 0 {
		if decodedPassword, er2 := c.Decrypter.Decrypt(user.Password, c.EncryptionKey); er2 != nil {
			RespondString(w, r, http.StatusBadRequest, "cannot decrypt password: "+er2.Error())
			return
		} else {
			user.Password = decodedPassword
		}
	}

	result, er3 := c.Authenticator.Authenticate(ctx, user)
	if er3 != nil {
		result.Status = StatusSystemError
		Respond(w, r, http.StatusOK, result, c.LogWriter, "Authentication", "Sign in", false, er3.Error())
	} else {
		Respond(w, r, http.StatusOK, result, c.LogWriter, "Authentication", "Sign in", true, "")
	}
}
func GetRemoteIp(r *http.Request) string {
	remoteIP, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		remoteIP = r.RemoteAddr
	}
	return remoteIP
}
