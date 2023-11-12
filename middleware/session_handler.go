package middleware

import (
	"context"
	"encoding/json"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
)

type CacheService interface {
	Get(ctx context.Context, key string) (string, error)
	Remove(ctx context.Context, key string) (bool, error)
	Expire(ctx context.Context, key string, timeToLive time.Duration) (bool, error)
}

type SessionHandler struct {
	EnableSession        bool
	PrefixSessionIndex   string
	AuthProcess          map[string]func(next http.Handler, options ...interface{}) http.Handler
	Whitelist            func(id string, token string) error
	SecretKey            string
	CookieName           string
	VerifyToken          func(tokenString string, secret string) (map[string]interface{}, jwt.StandardClaims, error)
	Cache                CacheService
	InternalAuthenticate http.Handler
	sessionExpiredTime   time.Duration
	LogError             func(ctx context.Context, format string, args ...interface{})
}

func NewSessionHandler(secretKey string, verifyToken func(tokenString string, secret string) (map[string]interface{}, jwt.StandardClaims, error), cache CacheService, sessionExpiredTime time.Duration, enableSession bool, logError func(ctx context.Context, format string, args ...interface{}), opts...string) *SessionHandler {
	var prefixSessionIndex, cookieName string
	if len(opts) > 0 {
		prefixSessionIndex = opts[0]
	} else {
		prefixSessionIndex = "index:"
	}
	if len(opts) > 1 {
		cookieName = opts[1]
	} else {
		cookieName = "id"
	}
	newHandler := &SessionHandler{
		AuthProcess:        make(map[string]func(next http.Handler, options ...interface{}) http.Handler),
		VerifyToken:        verifyToken,
		SecretKey:          secretKey,
		PrefixSessionIndex: prefixSessionIndex,
		CookieName:         cookieName,
		sessionExpiredTime: sessionExpiredTime,
		EnableSession:      enableSession,
		Cache:              cache,
		LogError:           logError,
	}
	return newHandler
}

func (h *SessionHandler) Middleware(next http.Handler, skipRefreshTTL bool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		au := r.Header["Authorization"]
		if h.EnableSession {
			sessionId := ""
			// case if set sessionID in cookie, need get token from cookie
			cookie, err := r.Cookie(h.CookieName)
			if err != nil {
				http.Error(w, "invalid Authorization token", http.StatusUnauthorized)
				return
			}

			if cookie == nil || cookie.Value == "" {
				http.Error(w, "invalid Authorization token", http.StatusUnauthorized)
				return
			}
			sessionId = cookie.Value
			ctx := r.Context()
			if h.Cache != nil {
				var sessionData map[string]string
				s, err := h.Cache.Get(r.Context(), sessionId)
				if err != nil {
					http.Error(w, "Session is expired", http.StatusUnauthorized)
					return
				}
				err2 := json.Unmarshal([]byte(s), &sessionData)
				if err2 != nil {
					if h.LogError != nil {
						h.LogError(r.Context(), "error unmarshal: %s ", err2.Error())
					}
					http.Error(w, "Session is expired", http.StatusUnauthorized)
					return
				}
				if id, ok := sessionData["id"]; ok {
					uData := map[string]interface{}{}
					s, err := h.Cache.Get(r.Context(), h.PrefixSessionIndex+id)
					if err != nil {
						http.Error(w, "Session is expired", http.StatusUnauthorized)
						return
					}
					err2 := json.Unmarshal([]byte(s), &uData)
					if err2 != nil {
						if h.LogError != nil {
							h.LogError(r.Context(), "error unmarshal: %s ", err2.Error())
						}
						http.Error(w, "Session is expired", http.StatusInternalServerError)
						return
					}
					ip := getRemoteIp(r)
					sid, ok := uData["sid"]
					if !ok || sid != sessionId ||
						getValueString(uData, "userAgent") != r.UserAgent() ||
						getValueString(uData, "ip") != ip {
						http.Error(w, "You cannot use multiple devices with a single account", http.StatusUnauthorized)
						return
					}
				} else {
					http.Error(w, "Session is expired", http.StatusUnauthorized)
					return
				}

				azureToken := getValueStringInterface(sessionData, "azure_token")
				ctx = context.WithValue(ctx, "azure_token", azureToken)

				authorizationToken := getValueStringInterface(sessionData, "token")
				ctx = context.WithValue(ctx, "token", authorizationToken)
			}
			if funcProcess, ok := h.AuthProcess["Bearer"]; ok {
				funcProcess(next, skipRefreshTTL, sessionId).ServeHTTP(w, r.WithContext(ctx))
			} else {
				next.ServeHTTP(w, r.WithContext(ctx))
			}
		} else {
			tokens := strings.Split(au[0], " ")
			prefix := tokens[0]
			if len(h.AuthProcess) > 0 {
				if f, ok := h.AuthProcess[prefix]; ok {
					authorizationToken := tokens[1]
					ctx := r.Context()
					ctx = context.WithValue(ctx, "token", authorizationToken)
					f(next, skipRefreshTTL).ServeHTTP(w, r.WithContext(ctx))
				} else {
					http.Error(w, "invalid Authorization token", http.StatusUnauthorized)
					return
				}
			}
		}
	})
}

func (h *SessionHandler) VerifyBearerToken(next http.Handler, options ...interface{}) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, r *http.Request) {
		skipRefreshTTL := false
		sessionId := ""
		if h.EnableSession {
			if len(options) > 0 {
				skipRefreshTTL = options[0].(bool)
			}
			if len(options) > 1 {
				sessionId = options[1].(string)
			}
		}
		ctx := r.Context()
		authorizationToken, exists := ctx.Value("token").(string)
		if !exists || authorizationToken == "" {
			http.Error(writer, "invalid authorization token", http.StatusUnauthorized)
			return
		}
		payload, _, err := h.VerifyToken(authorizationToken, h.SecretKey)
		if err != nil {
			http.Error(writer, "invalid authorization token", http.StatusUnauthorized)
			return
		}
		ip := getRemoteIp(r)
		ctx = context.WithValue(ctx, "ip", ip)
		ctx = context.WithValue(ctx, "token", authorizationToken)
		for k, e := range payload {
			if len(k) > 0 {
				ctx = context.WithValue(ctx, k, e)
			}
		}
		if !skipRefreshTTL && sessionId != "" {
			_, err := h.Cache.Expire(ctx, sessionId, h.sessionExpiredTime)
			if err != nil {
				if h.LogError != nil {
					h.LogError(ctx, err.Error())
				}
				http.Error(writer, "error set expire sessionId", http.StatusInternalServerError)
				return
			}
		}

		next.ServeHTTP(writer, r.WithContext(ctx))
	})
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

func getValueString(data map[string]interface{}, key string) string {
	// Check if the key exists in the map
	if value, ok := data[key]; ok {
		// Key exists, return the corresponding value
		return value.(string)
	}
	// Key does not exist, return an empty string
	return ""
}

func getValueStringInterface(data map[string]string, key string) string {
	if value, ok := data[key]; ok {
		return value
	}
	return ""
}
