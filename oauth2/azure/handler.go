package azure

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/core-go/auth"
)

const internalServerError = "Internal Server Error"
const arr = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

type CacheService interface {
	Put(ctx context.Context, key string, obj interface{}, timeToLive time.Duration) error
}

func Random(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = arr[rand.Intn(len(arr))]
	}
	return string(b)
}
func Encode(length int, id string) string {
	if length > 9 {
		length = 9
	}
	rand.Seed(time.Now().UnixNano())
	str := Random(length)
	bytes := []byte(id)
	e := base64.StdEncoding.EncodeToString(bytes)
	return fmt.Sprintf("%d%s%s", length, e, str)
}

type AuthenticationHandler struct {
	Auth     func(ctx context.Context, authorization string) (*auth.UserAccount, bool, error)
	Error    func(context.Context, string, ...map[string]interface{})
	Log      func(ctx context.Context, resource string, action string, success bool, desc string) error
	Ip       string
	UserId   string
	Resource string
	Action   string
	Cache    CacheService
	Expired  time.Duration
	Host     string
	Generate func(ctx context.Context) (string, error)
}

func NewAuthenticationHandlerWithCache(authenticate func(ctx context.Context, authorization string) (*auth.UserAccount, bool, error), logError func(context.Context, string, ...map[string]interface{}), cache CacheService, generate func(ctx context.Context) (string, error), expired time.Duration, host string, writeLog func(context.Context, string, string, bool, string) error, options ...string) *AuthenticationHandler {
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
	return &AuthenticationHandler{Auth: authenticate, Resource: resource, Action: action, Error: logError, Ip: ip, UserId: userId, Log: writeLog, Cache: cache, Generate: generate, Expired: expired, Host: host}
}
func NewAuthenticationHandler(authenticate func(ctx context.Context, authorization string) (*auth.UserAccount, bool, error), logError func(context.Context, string, ...map[string]interface{}), options ...func(context.Context, string, string, bool, string) error) *AuthenticationHandler {
	var writeLog func(context.Context, string, string, bool, string) error
	if len(options) >= 1 {
		writeLog = options[0]
	}
	return NewAuthenticationHandlerWithCache(authenticate, logError, nil, nil, time.Duration(10 * time.Second), "", writeLog, "ip", "userId", "authentication", "authenticate")
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
	session := make(map[string]string)
	session["token"] = user.Token
	session["id"] = user.Id
	session["ip"] = ip
	session["userAgent"] = r.UserAgent()
	err := h.Cache.Put(r.Context(), user.Id, session, h.Expired)
	if err != nil {
		h.Error(r.Context(), err.Error())
		respond(w, r, http.StatusInternalServerError, nil, h.Log, h.Resource, h.Action, false, err.Error())
		return
	}
	if h.Cache != nil && h.Generate != nil && len(h.Host) > 0 {
		host := r.Header.Get("Origin")
		if strings.Contains(host, h.Host) || strings.Contains(host, "localhost") {
			u, err := url.Parse(host)
			if err != nil {
				respond(w, r, http.StatusInternalServerError, nil, h.Log, h.Resource, h.Action, false, err.Error())
				return
			}
			host = strings.TrimPrefix(u.Hostname(), "www.")
		}
		uuid, err := h.Generate(r.Context())
		if err != nil {
			respond(w, r, http.StatusInternalServerError, nil, h.Log, h.Resource, h.Action, false, err.Error())
			return
		}
		http.SetCookie(w, &http.Cookie{
			Name: "id",
			Domain: host,
			Value: Encode(5, uuid),
			HttpOnly: true,
			Path: "/",
			MaxAge: 0,
			Expires: time.Now().Add(h.Expired),
			SameSite: http.SameSiteStrictMode,
			Secure: true,
		})
		user.Token = ""
	}
	respond(w, r, http.StatusOK, user, h.Log, h.Resource, h.Action, true, "")
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

func respond(w http.ResponseWriter, r *http.Request, code int, result interface{}, writeLog func(context.Context, string, string, bool, string) error, resource string, action string, success bool, desc string) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	err := json.NewEncoder(w).Encode(result)
	if writeLog != nil {
		writeLog(r.Context(), resource, action, success, desc)
	}
	return err
}
