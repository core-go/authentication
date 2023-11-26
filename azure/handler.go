package azure

import (
	"context"
	"encoding/json"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/core-go/auth"
)

const internalServerError = "Internal Server Error"

type CacheService interface {
	Put(ctx context.Context, key string, obj interface{}, timeToLive time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Remove(ctx context.Context, key string) (bool, error)
}

type AuthenticationHandler struct {
	Auth     func(ctx context.Context, authorization string) (*auth.UserAccount, bool, error)
	Error    func(context.Context, string, ...map[string]interface{})
	Log      func(ctx context.Context, resource string, action string, success bool, desc string) error
	Ip       string
	UserId   string
	CookieName string
	PrefixSessionIndex string
	Resource string
	Action   string
	LogoutAction string
	Cache    CacheService
	Expired  time.Duration
	Host     string
	Generate func(ctx context.Context) (string, error)
	SingleSession bool
	SId      string
	Id       string
	SameSite http.SameSite
}

func NewAuthenticationHandlerWithCache(authenticate func(ctx context.Context, authorization string) (*auth.UserAccount, bool, error), logError func(context.Context, string, ...map[string]interface{}), cache CacheService, generate func(ctx context.Context) (string, error), expired time.Duration, host string, sameSite http.SameSite, singleSession bool, writeLog func(context.Context, string, string, bool, string) error, options ...string) *AuthenticationHandler {
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
	return &AuthenticationHandler{Auth: authenticate, Resource: resource, Action: action, Error: logError, SameSite: sameSite, SingleSession: singleSession, Ip: ip, Id: id, SId: sid, UserId: userId, CookieName: cookieName, PrefixSessionIndex: prefixSessionIndex, Log: writeLog, Cache: cache, Generate: generate, Expired: expired, Host: host, LogoutAction: logoutAction}
}
func NewAuthenticationHandler(authenticate func(ctx context.Context, authorization string) (*auth.UserAccount, bool, error), logError func(context.Context, string, ...map[string]interface{}), options ...func(context.Context, string, string, bool, string) error) *AuthenticationHandler {
	var writeLog func(context.Context, string, string, bool, string) error
	if len(options) >= 1 {
		writeLog = options[0]
	}
	return NewAuthenticationHandlerWithCache(authenticate, logError, nil, nil, time.Duration(10 * time.Second), "", http.SameSiteStrictMode, true, writeLog, "ip", "userId", "authentication", "authenticate", "logout")
}

func (h *AuthenticationHandler) Authenticate(w http.ResponseWriter, r *http.Request) {
	var authorization string
	er1 := json.NewDecoder(r.Body).Decode(&authorization)
	if er1 != nil {
		if h.Error != nil {
			msg := "cannot decode authentication info: " + er1.Error()
			h.Error(r.Context(), msg)
		}
		http.Error(w, "cannot decode authentication info", http.StatusBadRequest)
		return
	}

	var ctx context.Context
	ctx = r.Context()
	ip := getRemoteIp(r)
	if len(h.Ip) > 0 {
		ctx = context.WithValue(ctx, h.Ip, ip)
		r = r.WithContext(ctx)
	}

	user, isExpired, er2 := h.Auth(r.Context(), authorization)
	if er2 != nil {
		if h.Error != nil {
			h.Error(r.Context(), er2.Error())
		}
		respond(w, r, http.StatusInternalServerError, internalServerError, h.Log, h.Resource, h.Action, false, er2.Error())
		return
	}
	if isExpired {
		respond(w, r, http.StatusUnauthorized, expired, h.Log, h.Resource, h.Action, false, "")
		return
	}
	if len(h.UserId) > 0 && len(user.Id) > 0 {
		ctx = context.WithValue(ctx, h.UserId, user.Id)
		r = r.WithContext(ctx)
	}
	if h.Cache != nil && h.Generate != nil && len(h.Host) > 0 {
		if h.SingleSession {
			indexData := make(map[string]interface{})
			data1, _ := h.Cache.Get(r.Context(), h.PrefixSessionIndex+user.Id)
			if len(data1) > 0 {
				err := json.Unmarshal([]byte(data1), &indexData)
				if err != nil {
					respond(w, r, http.StatusInternalServerError, nil, h.Log, h.Resource, h.Action, false, err.Error())
					return
				}
				sid := GetString(indexData, h.SId)
				if len(sid) > 0 {
					_, err2 := h.Cache.Remove(r.Context(), sid)
					if err2 != nil {
						respond(w, r, http.StatusInternalServerError, nil, h.Log, h.Resource, h.Action, false, err2.Error())
						return
					}
				}
			}
		}
		session := make(map[string]string)
		session["token"] = user.Token
		session["azure_token"] = authorization
		session[h.Id] = user.Id
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
		err2 := h.Cache.Put(r.Context(), h.PrefixSessionIndex + user.Id, indexData, h.Expired)
		if err2 != nil {
			h.Error(r.Context(), err.Error())
			respond(w, r, http.StatusInternalServerError, nil, h.Log, h.Resource, h.Action, false, err2.Error())
			return
		}
		err2 = h.Cache.Put(r.Context(), sessionId, session, h.Expired)
		if err2 != nil {
			h.Error(r.Context(), err.Error())
			respond(w, r, http.StatusInternalServerError, nil, h.Log, h.Resource, h.Action, false, err2.Error())
			return
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
		user.Token = ""
	}
	respond(w, r, http.StatusOK, user, h.Log, h.Resource, h.Action, true, "")
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
	data, err := GetCookie(r.Context(), cookie.Value, h.SId, h.Cache.Get)
	if err != nil {
		if err.Error() == "redis: nil" {
			respond(w, r, http.StatusOK, 1, h.Log, h.Resource, h.LogoutAction, true, err.Error())
			return
		}
	}
	sessionId := GetString(data, h.SId)
	if len(sessionId) > 0 {
		_, err = h.Cache.Remove(r.Context(), sessionId)
		if err != nil {
			respond(w, r, http.StatusInternalServerError, "", h.Log, h.Resource, h.LogoutAction, false, err.Error())
			return
		}
	}
	userId := GetString(data, h.Id)
	if len(userId) > 0 {
		_, err = h.Cache.Remove(r.Context(), h.PrefixSessionIndex + userId)
		if err != nil {
			respond(w, r, http.StatusInternalServerError, "", h.Log, h.Resource, h.LogoutAction, false, err.Error())
			return
		}
	}
	respond(w, r, http.StatusOK, 1, h.Log, h.Resource, h.LogoutAction, true, "")
}
func getRemoteIp(r *http.Request) string {
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
func GetString(data map[string]interface{}, key string) string {
	if data == nil {
		return ""
	}
	if value, ok := data[key]; ok {
		return value.(string)
	}
	return ""
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
func respond(w http.ResponseWriter, r *http.Request, code int, result interface{}, writeLog func(context.Context, string, string, bool, string) error, resource string, action string, success bool, desc string) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	err := json.NewEncoder(w).Encode(result)
	if writeLog != nil {
		writeLog(r.Context(), resource, action, success, desc)
	}
	return err
}
