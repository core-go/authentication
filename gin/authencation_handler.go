package gin

import (
	"context"
	"encoding/json"
	"errors"
	a "github.com/core-go/auth"
	"github.com/gin-gonic/gin"
	"net"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

const internalServerError = "Internal Server Error"

type AuthenticationHandler struct {
	Auth          func(ctx context.Context, user a.AuthInfo) (a.AuthResult, error)
	SystemError   int
	Timeout       int
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

func NewAuthenticationHandlerWithDecrypter(authenticate func(context.Context, a.AuthInfo) (a.AuthResult, error), systemError int, timeout int, logError func(context.Context, string), addTokenIntoWhitelist func(id string, token string) error, ipFromRequest bool, decrypt func(cipherText string, secretKey string) (string, error), encryptionKey string, writeLog func(context.Context, string, string, bool, string) error, options ...string) *AuthenticationHandler {
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
	return &AuthenticationHandler{Auth: authenticate, SystemError: systemError, Timeout: timeout, Resource: resource, Action: action, Error: logError, Ip: ip, UserId: userId, Whitelist: addTokenIntoWhitelist, Log: writeLog, Decrypt: decrypt, EncryptionKey: encryptionKey, IpFromRequest: ipFromRequest}
}
func NewAuthenticationHandler(authenticate func(context.Context, a.AuthInfo) (a.AuthResult, error), systemError int, timeout int, logError func(context.Context, string), options ...func(context.Context, string, string, bool, string) error) *AuthenticationHandler {
	var writeLog func(context.Context, string, string, bool, string) error
	if len(options) >= 1 {
		writeLog = options[0]
	}
	return NewAuthenticationHandlerWithDecrypter(authenticate, systemError, timeout, logError, nil, true, nil, "", writeLog, "ip", "userId", "authentication", "authenticate")
}
func NewAuthenticationHandlerWithWhitelist(authenticate func(context.Context, a.AuthInfo) (a.AuthResult, error), systemError int, timeout int, logError func(context.Context, string), addTokenIntoWhitelist func(id string, token string) error, ipFromRequest bool, options ...func(context.Context, string, string, bool, string) error) *AuthenticationHandler {
	var writeLog func(context.Context, string, string, bool, string) error
	if len(options) >= 1 {
		writeLog = options[0]
	}
	return NewAuthenticationHandlerWithDecrypter(authenticate, systemError, timeout, logError, addTokenIntoWhitelist, ipFromRequest, nil, "", writeLog, "ip", "userId", "authentication", "authenticate")
}

func (h *AuthenticationHandler) Authenticate(ctx *gin.Context) {
	r := ctx.Request
	var user a.AuthInfo
	if strings.Contains(r.Header.Get("Content-Type"), "multipart/form-data") {
		if err := r.ParseMultipartForm(1073741824); err != nil {
			ctx.String(http.StatusBadRequest, "cannot parse form: "+err.Error())
			return
		}
		modelType := reflect.TypeOf(user)
		mapIndexModels, err := getIndexesByTagJson(modelType)
		if err != nil {
			if h.Error != nil {
				h.Error(r.Context(), "cannot decode authentication info: "+err.Error())
			}
			ctx.String(http.StatusBadRequest, "cannot decode authentication info")
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
						ctx.String(http.StatusBadRequest, "invalid key: "+k)
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
			ctx.String(http.StatusBadRequest, "cannot decode authentication info")
			return
		}
	}

	var ctx2 context.Context
	ctx2 = r.Context()
	if len(h.Ip) > 0 {
		var ip string
		if len(user.Ip) > 0 && h.IpFromRequest {
			ip = user.Ip
		} else {
			ip = getRemoteIp(r)
		}
		ctx2 = context.WithValue(ctx2, h.Ip, ip)
		r = r.WithContext(ctx2)
		ctx.Request = r
	}

	if h.Decrypt != nil && len(h.EncryptionKey) > 0 {
		if decodedPassword, er2 := h.Decrypt(user.Password, h.EncryptionKey); er2 != nil {
			if h.Error != nil {
				h.Error(r.Context(), "cannot decrypt password: "+er2.Error())
			}
			ctx.String(http.StatusBadRequest, "cannot decrypt password")
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
			respond(ctx, http.StatusGatewayTimeout, "timeout", h.Log, h.Resource, h.Action, false, er3.Error())
		} else {
			result.Status = h.SystemError
			respond(ctx, http.StatusInternalServerError, result, h.Log, h.Resource, h.Action, false, er3.Error())
		}
	} else {
		if h.Whitelist != nil {
			h.Whitelist(result.User.Id, result.User.Token)
		}
		if len(h.UserId) > 0 && result.User != nil && len(result.User.Id) > 0 {
			ctx2 = context.WithValue(ctx2, h.UserId, result.User.Id)
			ctx.Request = r.WithContext(ctx2)
		}
		respond(ctx, http.StatusOK, result, h.Log, h.Resource, h.Action, true, "")
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
func respond(ctx *gin.Context, code int, result interface{}, writeLog func(context.Context, string, string, bool, string) error, resource string, action string, success bool, desc string) {
	ctx.JSON(code, result)
	if writeLog != nil {
		writeLog(ctx.Request.Context(), resource, action, success, desc)
	}
}
