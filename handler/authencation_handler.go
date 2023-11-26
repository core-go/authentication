package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"

	. "github.com/core-go/auth"
)

type AuthenticationHandler struct {
	Auth          func(ctx context.Context, user AuthInfo) (AuthResult, error)
	SystemError   int
	Timeout       int
	Error         func(context.Context, string, ...map[string]interface{})
	Ip            string
	UserId        string
	Whitelist     func(id string, token string) error
	IpFromRequest bool
	Log           func(ctx context.Context, resource string, action string, success bool, desc string) error
	Resource      string
	Action        string
	Cookie        bool
	CookieName    string
	Host          string
	Decrypt       func(string) (string, error)
}

func NewAuthenticationHandlerWithDecrypter(authenticate func(context.Context, AuthInfo) (AuthResult, error), systemError int, timeout int, logError func(context.Context, string, ...map[string]interface{}), addTokenIntoWhitelist func(id string, token string) error, cookie bool, ipFromRequest bool, decrypt func(string) (string, error), writeLog func(context.Context, string, string, bool, string) error, options ...string) *AuthenticationHandler {
	var ip, userId, cookieName, resource, action string
	if len(options) > 0 {
		ip = options[0]
	} else {
		ip = "ip"
	}
	if len(options) > 1 {
		userId = options[1]
	} else {
		userId = "userId"
	}
	if len(options) > 2 {
		cookieName = options[2]
	} else {
		cookieName = "id"
	}
	if len(options) > 3 {
		resource = options[3]
	} else {
		resource = "authentication"
	}
	if len(options) > 4 {
		action = options[4]
	} else {
		action = "authenticate"
	}
	return &AuthenticationHandler{Auth: authenticate, SystemError: systemError, Timeout: timeout, Cookie: cookie, CookieName: cookieName, Resource: resource, Action: action, Error: logError, Ip: ip, UserId: userId, Whitelist: addTokenIntoWhitelist, Log: writeLog, Decrypt: decrypt, IpFromRequest: ipFromRequest}
}
func NewAuthenticationHandler(authenticate func(context.Context, AuthInfo) (AuthResult, error), systemError int, timeout int, logError func(context.Context, string, ...map[string]interface{}), options ...func(context.Context, string, string, bool, string) error) *AuthenticationHandler {
	var writeLog func(context.Context, string, string, bool, string) error
	if len(options) >= 1 {
		writeLog = options[0]
	}
	return NewAuthenticationHandlerWithDecrypter(authenticate, systemError, timeout, logError, nil, true, true, nil, writeLog, "ip", "userId", "authentication", "authenticate")
}
func NewAuthenticationHandlerWithWhitelist(authenticate func(context.Context, AuthInfo) (AuthResult, error), systemError int, timeout int, logError func(context.Context, string, ...map[string]interface{}), addTokenIntoWhitelist func(id string, token string) error, cookie bool, ipFromRequest bool, options ...func(context.Context, string, string, bool, string) error) *AuthenticationHandler {
	var writeLog func(context.Context, string, string, bool, string) error
	if len(options) >= 1 {
		writeLog = options[0]
	}
	return NewAuthenticationHandlerWithDecrypter(authenticate, systemError, timeout, logError, addTokenIntoWhitelist, cookie, ipFromRequest, nil, writeLog, "ip", "userId", "authentication", "authenticate")
}

func (h *AuthenticationHandler) Authenticate(w http.ResponseWriter, r *http.Request) {
	var user AuthInfo
	if strings.Contains(r.Header.Get("Content-Type"), "multipart/form-data") {
		if err := r.ParseMultipartForm(1073741824); err != nil {
			http.Error(w, "cannot parse form: "+err.Error(), http.StatusBadRequest)
			return
		}
		modelType := reflect.TypeOf(user)
		mapIndexModels, err := getIndexesByTagJson(modelType)
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
					valueField, err := parseIntWithType(v[0], idType)
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
			ip = getRemoteIp(r)
		}
		ctx = context.WithValue(ctx, h.Ip, ip)
		r = r.WithContext(ctx)
	}

	if h.Decrypt != nil {
		if decodedPassword, er2 := h.Decrypt(user.Password); er2 != nil {
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
		if h.Error != nil {
			h.Error(r.Context(), er3.Error())
		}
		if result.Status == h.Timeout {
			respond(w, r, http.StatusGatewayTimeout, "timeout", h.Log, h.Resource, h.Action, false, er3.Error())
		} else {
			result.Status = h.SystemError
			respond(w, r, http.StatusInternalServerError, result, h.Log, h.Resource, h.Action, false, er3.Error())
		}
	} else {
		if h.Whitelist != nil {
			h.Whitelist(result.User.Id, result.User.Token)
		}
		if len(h.UserId) > 0 && result.User != nil && len(result.User.Id) > 0 {
			ctx = context.WithValue(ctx, h.UserId, result.User.Id)
			r = r.WithContext(ctx)
		}
		if h.Cookie {
			var token string
			expired := time.Now()
			if result.User != nil {
				token = result.User.Token
				if result.User.TokenExpiredTime != nil {
					expired = *result.User.TokenExpiredTime
					result.User.TokenExpiredTime = nil
				}
			}
			host := r.Header.Get("Origin")
			if strings.Contains(host, h.Host) || strings.Contains(host, "localhost") {
				u, err := url.Parse(host)
				if err != nil {
					respond(w, r, http.StatusInternalServerError, nil, h.Log, h.Resource, h.Action, false, err.Error())
					return
				}
				host = strings.TrimPrefix(u.Hostname(), "www.")
			}
			if token == "" {
				http.Error(w, "cannot get token", http.StatusUnauthorized)
				return
			}
			http.SetCookie(w, &http.Cookie{
				Name: h.CookieName,
				Domain: host,
				Value: token,
				HttpOnly: true,
				Path: "/",
				MaxAge: 0,
				Expires: expired,
				SameSite: http.SameSiteStrictMode,
				Secure: true,
			})
			result.User.Token = ""
		}
		respond(w, r, http.StatusOK, result, h.Log, h.Resource, h.Action, true, "")
	}
}
func getRemoteIp(r *http.Request) string {
	remoteIP, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		remoteIP = r.RemoteAddr
	}
	return remoteIP
}
func getIndexesByTagJson(modelType reflect.Type) (map[string]int, error) {
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
func parseIntWithType(value string, idType string) (v interface{}, err error) {
	switch idType {
	case "int64", "*int64":
		return strconv.ParseInt(value, 10, 64)
	case "int", "int32", "*int32":
		return strconv.Atoi(value)
	default:
	}
	return value, nil
}
