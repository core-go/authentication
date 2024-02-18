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
)

const expired = "Token is expired"

type StoreService interface {
	Put(ctx context.Context, key string, obj interface{}, timeToLive time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Remove(ctx context.Context, key string) (bool, error)
}

type AuthenticationHandler struct {
	Auth               Authenticate
	SystemError        int
	Timeout            int
	Error              func(context.Context, string, ...map[string]interface{})
	Ip                 string
	UserId             string
	Whitelist          func(id string, token string) error
	IpFromRequest      bool
	Log                func(ctx context.Context, resource string, action string, success bool, desc string) error
	Resource           string
	Action             string
	Cookie             bool
	PrefixSessionIndex string
	CookieName         string
	Host               string
	SameSite           http.SameSite
	Expired            time.Duration
	SingleSession      bool
	Id                 string
	SId                string
	Generate           func(ctx context.Context) (string, error)
	LogoutAction       string
	Store              StoreService

	RefreshExpire   func(w http.ResponseWriter, sessionId string) error
	DecodeSessionID func(value string) (string, error)
	EncodeSessionID func(sid string) string

	Decrypt func(string) (string, error)
}
type LogError func(context.Context, string, ...map[string]interface{})
type Authenticate func(context.Context, AuthInfo) (AuthResult, error)

func NewAuthenticationHandlerWithDecrypter(authenticate Authenticate,
	systemError int,
	timeout int,
	logError func(context.Context, string, ...map[string]interface{}),
	addTokenIntoWhitelist func(id string, token string) error,
	cookie bool,
	ipFromRequest bool,
	sameSite http.SameSite,
	decrypt func(string) (string, error),
	writeLog func(context.Context, string, string, bool, string) error,
	options ...string) *AuthenticationHandler {
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
	return &AuthenticationHandler{Auth: authenticate, SystemError: systemError, Timeout: timeout, SameSite: sameSite, Cookie: cookie, CookieName: cookieName, Resource: resource, Action: action, Error: logError, Ip: ip, UserId: userId, Whitelist: addTokenIntoWhitelist, Log: writeLog, Decrypt: decrypt, IpFromRequest: ipFromRequest}
}
func NewAuthenticationHandlerWithCache(authenticate Authenticate, systemError int, timeout int, logError LogError,
	store StoreService,
	generate func(ctx context.Context) (string, error),
	expired time.Duration,
	host string,
	sameSite http.SameSite,
	enableCookie bool,
	singleSession bool,
	writeLog func(context.Context, string, string, bool, string) error,
	options ...string) *AuthenticationHandler {
	var ip, id, sid, userId, cookieName, prefixSessionIndex, resource, action, logoutAction string
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
		prefixSessionIndex = options[3]
	} else {
		prefixSessionIndex = "index:"
	}
	if len(options) > 4 {
		sid = options[4]
	} else {
		sid = "sid"
	}
	if len(options) > 5 {
		id = options[5]
	} else {
		id = "id"
	}
	if len(options) > 6 {
		resource = options[6]
	} else {
		resource = "authentication"
	}
	if len(options) > 7 {
		action = options[7]
	} else {
		action = "authenticate"
	}
	if len(options) > 8 {
		logoutAction = options[8]
	} else {
		logoutAction = "logout"
	}
	return &AuthenticationHandler{
		Auth:               authenticate,
		Resource:           resource,
		Action:             action,
		Error:              logError,
		SameSite:           sameSite,
		SingleSession:      singleSession,
		Ip:                 ip,
		Id:                 id,
		SId:                sid,
		UserId:             userId,
		Cookie:             enableCookie,
		CookieName:         cookieName,
		PrefixSessionIndex: prefixSessionIndex,
		Log:                writeLog,
		Store:              store,
		Generate:           generate,
		Expired:            expired,
		Host:               host,
		LogoutAction:       logoutAction,
	}
}

func NewAuthenticationHandler(authenticate func(context.Context, AuthInfo) (AuthResult, error), systemError int, timeout int, logError func(context.Context, string, ...map[string]interface{}), options ...func(context.Context, string, string, bool, string) error) *AuthenticationHandler {
	var writeLog func(context.Context, string, string, bool, string) error
	if len(options) >= 1 {
		writeLog = options[0]
	}
	return NewAuthenticationHandlerWithDecrypter(authenticate, systemError, timeout, logError,
		nil, true, true, http.SameSiteStrictMode, nil, writeLog,
		"ip", "userId", "id", "authenticate", "authenticate")
}

func NewAuthenticationHandlerWithWhitelist(authenticate func(context.Context, AuthInfo) (AuthResult, error), systemError int, timeout int, logError func(context.Context, string, ...map[string]interface{}), addTokenIntoWhitelist func(id string, token string) error, cookie bool, ipFromRequest bool, options ...func(context.Context, string, string, bool, string) error) *AuthenticationHandler {
	var writeLog func(context.Context, string, string, bool, string) error
	if len(options) >= 1 {
		writeLog = options[0]
	}
	return NewAuthenticationHandlerWithDecrypter(authenticate, systemError, timeout, logError, addTokenIntoWhitelist, cookie, ipFromRequest, http.SameSiteStrictMode, nil, writeLog, "ip", "userId", "authentication", "authenticate")
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
			host := r.Header.Get("Origin")
			if strings.Contains(host, h.Host) || strings.Contains(host, "localhost") {
				u, err := url.Parse(host)
				if err != nil {
					respond(w, r, http.StatusInternalServerError, nil, h.Log, h.Resource, h.Action, false, err.Error())
					return
				}
				host = strings.TrimPrefix(u.Hostname(), "www.")
			}
			expired := time.Now()
			if result.User != nil {
				token = result.User.Token
				if result.User.TokenExpiredTime != nil {
					expired = *result.User.TokenExpiredTime
					result.User.TokenExpiredTime = nil
				}
			}
			if token == "" {
				http.Error(w, "cannot get token", http.StatusUnauthorized)
				return
			}
			ip := getForwardedRemoteIp(r)
			if len(h.Ip) > 0 {
				ctx = context.WithValue(ctx, h.Ip, ip)
				r = r.WithContext(ctx)
			}
			if h.Store != nil && h.Generate != nil && len(h.Host) > 0 {
				if h.SingleSession {
					indexData := make(map[string]interface{})
					data1, _ := h.Store.Get(r.Context(), h.PrefixSessionIndex+result.User.Id)
					if len(data1) > 0 {
						err := json.Unmarshal([]byte(data1), &indexData)
						if err != nil {
							respond(w, r, http.StatusInternalServerError, nil, h.Log, h.Resource, h.Action, false, err.Error())
							return
						}
						sid := GetString(indexData, h.SId)
						if len(sid) > 0 {
							_, err2 := h.Store.Remove(r.Context(), sid)
							if err2 != nil {
								respond(w, r, http.StatusInternalServerError, nil, h.Log, h.Resource, h.Action, false, err2.Error())
								return
							}
						}
					}
				}
				session := make(map[string]string)
				session["token"] = token
				session[h.Id] = result.User.Id
				host := r.Header.Get("Origin")
				if strings.Contains(host, h.Host) || strings.Contains(host, "localhost") {
					u, err := url.Parse(host)
					if err != nil {
						respond(w, r, http.StatusInternalServerError, nil, h.Log, h.Resource, h.Action, false, err.Error())
						return
					}
					host = strings.TrimPrefix(u.Hostname(), "www.")
				}
				sessionId := ""
				uuid, err := h.Generate(r.Context())
				if err != nil {
					h.Error(r.Context(), err.Error())
					respond(w, r, http.StatusInternalServerError, nil, h.Log, h.Resource, h.Action, false, err.Error())
					return
				}
				sessionId = uuid
				indexData := make(map[string]string)
				indexData[h.SId] = sessionId
				indexData["ip"] = ip
				indexData["userAgent"] = r.UserAgent()
				err2 := h.Store.Put(r.Context(), h.PrefixSessionIndex+result.User.Id, indexData, h.Expired)
				if err2 != nil {
					h.Error(r.Context(), err.Error())
					respond(w, r, http.StatusInternalServerError, nil, h.Log, h.Resource, h.Action, false, err2.Error())
					return
				}
				err2 = h.Store.Put(r.Context(), sessionId, session, h.Expired)
				if err2 != nil {
					h.Error(r.Context(), err.Error())
					respond(w, r, http.StatusInternalServerError, nil, h.Log, h.Resource, h.Action, false, err2.Error())
					return
				}
				if h.EncodeSessionID != nil {
					sessionId = h.EncodeSessionID(sessionId)
				}
				http.SetCookie(w, &http.Cookie{
					Name:     h.CookieName,
					Domain:   host,
					Value:    sessionId,
					HttpOnly: true,
					Path:     "/",
					MaxAge:   0,
					Expires:  time.Now().Add(h.Expired),
					SameSite: h.SameSite,
					Secure:   true,
				})
				result.User.Token = ""
			} else {
				http.SetCookie(w, &http.Cookie{
					Name:     h.CookieName,
					Domain:   host,
					Value:    token,
					HttpOnly: true,
					Path:     "/",
					MaxAge:   0,
					Expires:  expired,
					SameSite: h.SameSite,
					Secure:   true,
				})
				result.User.Token = ""
			}
		}
		respond(w, r, http.StatusOK, result, h.Log, h.Resource, h.Action, true, "")
	}
}
func (h *AuthenticationHandler) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(h.CookieName)
	if err != nil {
		respond(w, r, http.StatusInternalServerError, "", h.Log, h.Resource, h.Action, false, err.Error())
		return
	}
	if cookie == nil || cookie.Value == "" {
		respond(w, r, http.StatusOK, expired, h.Log, h.Resource, h.LogoutAction, true, "")
		return
	}
	valueCookie := cookie.Value
	if h.DecodeSessionID != nil {
		valueCookie, err = h.DecodeSessionID(valueCookie)
		if err != nil {
			respond(w, r, http.StatusInternalServerError, "", h.Log, h.Resource, h.LogoutAction, false, err.Error())
			return
		}
	}
	data, err := GetCookie(r.Context(), valueCookie, h.SId, h.Store.Get)
	if err != nil {
		if err.Error() == "redis: nil" {
			respond(w, r, http.StatusOK, 1, h.Log, h.Resource, h.LogoutAction, true, err.Error())
			return
		}
	}
	sessionId := GetString(data, h.SId)
	if len(sessionId) > 0 {
		_, err = h.Store.Remove(r.Context(), sessionId)
		if err != nil {
			respond(w, r, http.StatusInternalServerError, "", h.Log, h.Resource, h.LogoutAction, false, err.Error())
			return
		}
		http.SetCookie(w, &http.Cookie{
			Name:     h.CookieName,
			Domain:   h.Host,
			Value:    "",
			HttpOnly: true,
			Path:     "/",
			MaxAge:   -1,
			SameSite: h.SameSite,
			Secure:   true,
		})
	}
	userId := GetString(data, h.Id)
	if len(userId) > 0 {
		_, err = h.Store.Remove(r.Context(), h.PrefixSessionIndex+userId)
		if err != nil {
			respond(w, r, http.StatusInternalServerError, "", h.Log, h.Resource, h.LogoutAction, false, err.Error())
			return
		}
	}
	respond(w, r, http.StatusOK, 1, h.Log, h.Resource, h.LogoutAction, true, "")
}
func GetCookie(ctx context.Context, value string, sid string, cache func(context.Context, string) (string, error)) (map[string]interface{}, error) {
	var data map[string]interface{}
	s, err := cache(ctx, value)
	if err != nil {
		return data, err
	}
	if len(s) > 0 {
		err = json.Unmarshal([]byte(s), &data)
		if err != nil {
			return nil, err
		}
	}
	data[sid] = value
	return data, err
}
func getRemoteIp(r *http.Request) string {
	remoteIP, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		remoteIP = r.RemoteAddr
	}
	return remoteIP
}
func getForwardedRemoteIp(r *http.Request) string {
	ips := r.Header.Get("X-FORWARDED-FOR")
	splitIps := strings.Split(ips, ",")
	for _, ip := range splitIps {
		netIP := net.ParseIP(ip)
		if netIP != nil {
			return ip
		}
	}
	return ""
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
func GetString(data map[string]interface{}, key string) string {
	if data == nil {
		return ""
	}
	if value, ok := data[key]; ok {
		return value.(string)
	}
	return ""
}
