package auth

import (
	"context"
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

type AuthenticationHandler struct {
	Authenticator Authenticator
	Resource      string
	Action        string
	LogError      func(context.Context, string)
	Ip            string
	UserId        string
	Whitelist     func(id string, token string) error
	WriteLog      func(ctx context.Context, resource string, action string, success bool, desc string) error
	Decrypter     PasswordDecrypter
	EncryptionKey string
	IpFromRequest bool
}

func NewAuthenticationHandlerWithDecrypter(authenticator Authenticator, resource string, action string, logError func(context.Context, string), ip string, userId string, addTokenIntoWhitelist func(id string, token string) error, writeLog func(context.Context, string, string, bool, string) error, decrypter PasswordDecrypter, encryptionKey string, ipFromRequest bool) *AuthenticationHandler {
	if len(ip) == 0 {
		ip = "ip"
	}
	if len(userId) == 0 {
		userId = "userId"
	}
	if len(resource) == 0 {
		resource = "authentication"
	}
	if len(action) == 0 {
		action = "authenticate"
	}
	return &AuthenticationHandler{Authenticator: authenticator, Resource: resource, Action: action, LogError: logError, Ip: ip, UserId: userId, Whitelist: addTokenIntoWhitelist, WriteLog: writeLog, Decrypter: decrypter, EncryptionKey: encryptionKey, IpFromRequest: ipFromRequest}
}

func NewAuthenticationHandler(authenticator Authenticator, logError func(context.Context, string), writeLog func(context.Context, string, string, bool, string) error, addTokenIntoWhitelist func(id string, token string) error, ipFromRequest bool) *AuthenticationHandler {
	return NewAuthenticationHandlerWithDecrypter(authenticator, "authentication", "authenticate", logError, "ip", "userId", addTokenIntoWhitelist, writeLog, nil, "", ipFromRequest)
}

func (h *AuthenticationHandler) Authenticate(w http.ResponseWriter, r *http.Request) {
	var user AuthInfo
	if strings.Contains(r.Header.Get("Content-Type"), "multipart/form-data") {
		if err := r.ParseMultipartForm(1073741824); err != nil {
			http.Error(w, "cannot parse form: "+err.Error(), http.StatusBadRequest)
			return
		}
		modelType := reflect.TypeOf(user)
		mapIndexModels, err := GetIndexesByTagJson(modelType)
		if err != nil {
			if h.LogError != nil {
				h.LogError(r.Context(), "cannot decode authentication info: "+err.Error())
			}
			http.Error(w, "cannot decode authentication info", http.StatusBadRequest)
			return
		}

		postForm := r.PostForm
		userV := reflect.Indirect(reflect.ValueOf(&user))
		for k, v := range postForm {
			if index, ok := mapIndexModels[k]; ok {
				idType := userV.Field(index).Type().String()
				if strings.Index(idType, "int") >= 0 {
					valueField, err := ParseIntWithType(v[0], idType)
					if err != nil {
						http.Error(w, "invalid key: "+k, http.StatusBadRequest)
						return
					}
					userV.Field(index).Set(reflect.ValueOf(valueField))
				} else {
					userV.Field(index).Set(reflect.ValueOf(v[0]))
				}
			}
		}
	} else {
		er1 := json.NewDecoder(r.Body).Decode(&user)
		if er1 != nil {
			if h.LogError != nil {
				msg := "cannot decode authentication info: " + er1.Error()
				h.LogError(r.Context(), msg)
			}
			http.Error(w, "cannot decode authentication info", http.StatusBadRequest)
			return
		}
	}

	var ctx context.Context
	ctx = r.Context()
	if len(h.Ip) > 0 {
		var ip string
		if len(user.Ip) > 0 && h.IpFromRequest {
			ip = user.Ip
		} else {
			ip = GetRemoteIp(r)
		}
		ctx = context.WithValue(ctx, h.Ip, ip)
		r = r.WithContext(ctx)
	}

	if h.Decrypter != nil && len(h.EncryptionKey) > 0 {
		if decodedPassword, er2 := h.Decrypter.Decrypt(user.Password, h.EncryptionKey); er2 != nil {
			if h.LogError != nil {
				h.LogError(r.Context(), "cannot decrypt password: "+er2.Error())
			}
			http.Error(w, "cannot decrypt password", http.StatusBadRequest)
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
		respond(w, r, http.StatusOK, result, h.WriteLog, h.Resource, h.Action, false, er3.Error())
	} else {
		if h.Whitelist != nil {
			h.Whitelist(result.User.UserId, result.User.Token)
		}
		if len(h.UserId) > 0 && result.User != nil && len(result.User.UserId) > 0 {
			ctx = context.WithValue(ctx, h.UserId, result.User.UserId)
			r = r.WithContext(ctx)
		}
		respond(w, r, http.StatusOK, result, h.WriteLog, h.Resource, h.Action, true, "")
	}
}
func GetRemoteIp(r *http.Request) string {
	remoteIP, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		remoteIP = r.RemoteAddr
	}
	return remoteIP
}
func GetIndexesByTagJson(modelType reflect.Type) (map[string]int, error) {
	mapp := make(map[string]int, 0)
	if modelType.Kind() != reflect.Struct {
		return mapp, errors.New("bad type")
	}
	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)
		tagJson := field.Tag.Get("json")
		tagJson = strings.Split(tagJson, ",")[0]
		if len(tagJson) > 0 {
			mapp[tagJson] = i
		}
	}
	return mapp, nil
}
func ParseIntWithType(value string, idType string) (v interface{}, err error) {
	switch idType {
	case "int64", "*int64":
		return strconv.ParseInt(value, 10, 64)
	case "int", "int32", "*int32":
		return strconv.Atoi(value)
	default:
	}
	return value, nil
}
