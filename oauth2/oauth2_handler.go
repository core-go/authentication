package oauth2

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
)

type Oauth2ActionConfig struct {
	Ip             string `mapstructure:"ip"`
	Resource       string `mapstructure:"resource"`
	Authenticate   string `mapstructure:"authenticate"`
	Configuration  string `mapstructure:"configuration"`
	Configurations string `mapstructure:"configurations"`
}
type OAuth2Handler struct {
	OAuth2Service OAuth2Service
	SystemError   int
	Error         func(context.Context, string, ...map[string]interface{})
	Config        Oauth2ActionConfig
	Log           func(ctx context.Context, resource string, action string, success bool, desc string) error
}

func NewDefaultOAuth2Handler(oauth2Service OAuth2Service, systemError int, logError func(context.Context, string, ...map[string]interface{})) *OAuth2Handler {
	return NewOAuth2Handler(oauth2Service, systemError, logError, nil)
}

func NewOAuth2Handler(oauth2Service OAuth2Service, systemError int, logError func(context.Context, string, ...map[string]interface{}), writeLog func(context.Context, string, string, bool, string) error, options ...Oauth2ActionConfig) *OAuth2Handler {
	var c Oauth2ActionConfig
	if len(options) >= 1 {
		conf := options[0]
		c.Ip = conf.Ip
		c.Resource = conf.Resource
		c.Authenticate = conf.Authenticate
		c.Configuration = conf.Configuration
		c.Configurations = conf.Configurations
	}
	if len(c.Ip) == 0 {
		c.Ip = "ip"
	}
	if len(c.Resource) == 0 {
		c.Resource = "oauth2"
	}
	if len(c.Authenticate) == 0 {
		c.Authenticate = "authenticate"
	}
	if len(c.Configuration) == 0 {
		c.Configuration = "configuration"
	}
	if len(c.Configurations) == 0 {
		c.Configurations = "configurations"
	}
	return &OAuth2Handler{OAuth2Service: oauth2Service, SystemError: systemError, Config: c, Error: logError, Log: writeLog}
}
func (h *OAuth2Handler) Configuration(w http.ResponseWriter, r *http.Request) {
	id := ""
	if r.Method == "GET" {
		i := strings.LastIndex(r.RequestURI, "/")
		if i >= 0 {
			id = r.RequestURI[i+1:]
		}
	} else {
		b, er1 := ioutil.ReadAll(r.Body)
		if er1 != nil {
			http.Error(w, "body cannot is empty", http.StatusBadRequest)
			return
		}
		id = strings.Trim(string(b), " ")
	}
	if len(id) == 0 {
		http.Error(w, "request cannot is empty", http.StatusBadRequest)
		return
	}
	model, err := h.OAuth2Service.Configuration(r.Context(), id)
	if err != nil {
		if h.Error != nil {
			h.Error(r.Context(), err.Error())
		}
		respond(w, r, http.StatusOK, nil, h.Log, h.Config.Resource, h.Config.Configuration, false, err.Error())
	} else {
		respond(w, r, http.StatusOK, model, h.Log, h.Config.Resource, h.Config.Configuration, true, "")
	}
}
func (h *OAuth2Handler) Configurations(w http.ResponseWriter, r *http.Request) {
	model, err := h.OAuth2Service.Configurations(r.Context())
	if err != nil {
		if h.Error != nil {
			h.Error(r.Context(), err.Error())
		}
		respond(w, r, http.StatusOK, nil, h.Log, h.Config.Resource, h.Config.Configurations, false, err.Error())
	} else {
		respond(w, r, http.StatusOK, model, h.Log, h.Config.Resource, h.Config.Configurations, true, "")
	}
}
func (h *OAuth2Handler) Authenticate(w http.ResponseWriter, r *http.Request) {
	var request OAuth2Info
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		if h.Error != nil {
			h.Error(r.Context(), "cannot decode OAuth2Info model: "+err.Error())
		}
		http.Error(w, "cannot decode OAuth2Info model", http.StatusBadRequest)
		return
	}
	var authorization string
	if len(r.Header["Authorization"]) < 1 {
		authorization = ""
	} else {
		authorization = r.Header["Authorization"][0]
	}
	ip := getRemoteIp(r)
	var ctx context.Context
	ctx = r.Context()
	if len(h.Config.Ip) > 0 {
		ctx = context.WithValue(ctx, h.Config.Ip, ip)
		r = r.WithContext(ctx)
	}
	result, err := h.OAuth2Service.Authenticate(r.Context(), &request, authorization)
	if err != nil {
		result.Status = h.SystemError
		if h.Error != nil {
			h.Error(r.Context(), err.Error())
		}
		respond(w, r, http.StatusOK, result, h.Log, h.Config.Resource, h.Config.Authenticate, false, err.Error())
	} else {
		respond(w, r, http.StatusOK, result, h.Log, h.Config.Resource, h.Config.Authenticate, true, "")
	}
}

func respond(w http.ResponseWriter, r *http.Request, code int, result interface{}, writeLog func(context.Context, string, string, bool, string) error, resource string, action string, success bool, desc string) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	err := json.NewEncoder(w).Encode(result)
	if writeLog != nil {
		newCtx := context.WithValue(r.Context(), "request", r)
		writeLog(newCtx, resource, action, success, desc)
	}
	return err
}
func getRemoteIp(r *http.Request) string {
	remoteIP, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		remoteIP = r.RemoteAddr
	}
	return remoteIP
}
