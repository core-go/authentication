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
	Auth          func(ctx context.Context, user AuthInfo) (AuthResult, error)
	SystemError   int
	Error         func(context.Context, string)
	Ip            string
	UserId        string
	Whitelist     func(id string, token string) error
	IpFromRequest bool
	Log           func(ctx context.Context, resource string, action string, success bool, desc string) error
	Resource      string
	Action        string
	Decrypt       func(cipherText string, secretKey string) (string, error)
	EncryptionKey string
}

func NewAuthenticationHandlerWithDecrypter(authenticate func(context.Context, AuthInfo) (AuthResult, error), systemError int, logError func(context.Context, string), addTokenIntoWhitelist func(id string, token string) error, ipFromRequest bool, decrypt func(cipherText string, secretKey string) (string, error), encryptionKey string, writeLog func(context.Context, string, string, bool, string) error, options ...string) *AuthenticationHandler {
	var ip, userId, resource, action string
	if len(options) >= 1 {
		ip = options[0]
	} else {
		ip = "ip"
	}
	if len(options) >= 2 {
		userId = options[1]
	} else {
		userId = "userId"
	}
	if len(options) >= 3 {
		resource = options[2]
	} else {
		resource = "authentication"
	}
	if len(options) >= 4 {
		action = options[3]
	} else {
		action = "authenticate"
	}
	return &AuthenticationHandler{Auth: authenticate, SystemError: systemError, Resource: resource, Action: action, Error: logError, Ip: ip, UserId: userId, Whitelist: addTokenIntoWhitelist, Log: writeLog, Decrypt: decrypt, EncryptionKey: encryptionKey, IpFromRequest: ipFromRequest}
}

func NewAuthenticationHandler(authenticate func(context.Context, AuthInfo) (AuthResult, error), systemError int, logError func(context.Context, string), addTokenIntoWhitelist func(id string, token string) error, ipFromRequest bool, options ...func(context.Context, string, string, bool, string) error) *AuthenticationHandler {
	var writeLog func(context.Context, string, string, bool, string) error
	if len(options) >= 1 {
		writeLog = options[0]
	}
	return NewAuthenticationHandlerWithDecrypter(authenticate, systemError, logError, addTokenIntoWhitelist, ipFromRequest, nil, "", writeLog, "ip", "userId", "authentication", "authenticate")
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
			if h.Error != nil {
				h.Error(r.Context(), "cannot decode authentication info: "+err.Error())
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
			if h.Error != nil {
				msg := "cannot decode authentication info: " + er1.Error()
				h.Error(r.Context(), msg)
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

	if h.Decrypt != nil && len(h.EncryptionKey) > 0 {
		if decodedPassword, er2 := h.Decrypt(user.Password, h.EncryptionKey); er2 != nil {
			if h.Error != nil {
				h.Error(r.Context(), "cannot decrypt password: "+er2.Error())
			}
			http.Error(w, "cannot decrypt password", http.StatusBadRequest)
			return
		} else {
			user.Password = decodedPassword
		}
	}

	result, er3 := h.Auth(r.Context(), user)
	if er3 != nil {
		result.Status = h.SystemError
		if h.Error != nil {
			h.Error(r.Context(), er3.Error())
		}
		respond(w, r, http.StatusOK, result, h.Log, h.Resource, h.Action, false, er3.Error())
	} else {
		if h.Whitelist != nil {
			h.Whitelist(result.User.Id, result.User.Token)
		}
		if len(h.UserId) > 0 && result.User != nil && len(result.User.Id) > 0 {
			ctx = context.WithValue(ctx, h.UserId, result.User.Id)
			r = r.WithContext(ctx)
		}
		respond(w, r, http.StatusOK, result, h.Log, h.Resource, h.Action, true, "")
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