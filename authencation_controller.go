package auth

import (
	"encoding/json"
	"net/http"
)

type AuthenticationController struct {
	Authenticator Authenticator
	Decrypter     PasswordDecrypter
	EncryptionKey string
	LogService    AuthActivityLogService
}

func NewAuthenticationController(authenticationService Authenticator, decrypter PasswordDecrypter, encryptionKey string, logService AuthActivityLogService) *AuthenticationController {
	return &AuthenticationController{authenticationService, decrypter, encryptionKey, logService}
}

func NewDefaultAuthenticationController(authenticationService Authenticator) *AuthenticationController {
	return &AuthenticationController{authenticationService, nil, "", nil}
}

func (c *AuthenticationController) Authenticate(w http.ResponseWriter, r *http.Request) {
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

	result, er3 := c.Authenticator.Authenticate(r.Context(), user)
	if er3 != nil {
		result.Status = StatusSystemError
		Respond(w, r, http.StatusOK, result, c.LogService, "Authentication", "Sign in", false, er3.Error())
	} else {
		Respond(w, r, http.StatusOK, result, c.LogService, "Authentication", "Sign in", true, "")
	}
}
