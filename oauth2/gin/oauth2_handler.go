package gin

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net"
	"net/http"
	"strings"

	"github.com/core-go/authentication/oauth2"
)

type OAuth2Handler struct {
	OAuth2Service oauth2.OAuth2Service
	SystemError   int
	Error         func(context.Context, string, ...map[string]interface{})
	Config        oauth2.Oauth2ActionConfig
	Log           func(ctx context.Context, resource string, action string, success bool, desc string) error
}

func NewDefaultOAuth2Handler(oauth2Service oauth2.OAuth2Service, systemError int, logError func(context.Context, string, ...map[string]interface{})) *OAuth2Handler {
	return NewOAuth2Handler(oauth2Service, systemError, logError, nil)
}

func NewOAuth2Handler(oauth2Service oauth2.OAuth2Service, systemError int, logError func(context.Context, string, ...map[string]interface{}), writeLog func(context.Context, string, string, bool, string) error, options ...oauth2.Oauth2ActionConfig) *OAuth2Handler {
	var c oauth2.Oauth2ActionConfig
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
func (h *OAuth2Handler) Configuration(ctx *gin.Context) {
	r := ctx.Request
	id := ""
	if r.Method == "GET" {
		i := strings.LastIndex(r.RequestURI, "/")
		if i >= 0 {
			id = r.RequestURI[i+1:]
		}
	} else {
		b, er1 := ioutil.ReadAll(r.Body)
		if er1 != nil {
			ctx.String(http.StatusBadRequest, "body cannot is empty")
			return
		}
		id = strings.Trim(string(b), " ")
	}
	if len(id) == 0 {
		ctx.String(http.StatusBadRequest, "request cannot is empty")
		return
	}
	model, err := h.OAuth2Service.Configuration(r.Context(), id)
	if err != nil {
		if h.Error != nil {
			h.Error(r.Context(), err.Error())
		}
		respond(ctx, http.StatusOK, nil, h.Log, h.Config.Resource, h.Config.Configuration, false, err.Error())
	} else {
		respond(ctx, http.StatusOK, model, h.Log, h.Config.Resource, h.Config.Configuration, true, "")
	}
}
func (h *OAuth2Handler) Configurations(ctx *gin.Context) {
	model, err := h.OAuth2Service.Configurations(ctx.Request.Context())
	if err != nil {
		if h.Error != nil {
			h.Error(ctx.Request.Context(), err.Error())
		}
		respond(ctx, http.StatusOK, nil, h.Log, h.Config.Resource, h.Config.Configurations, false, err.Error())
	} else {
		respond(ctx, http.StatusOK, model, h.Log, h.Config.Resource, h.Config.Configurations, true, "")
	}
}
func (h *OAuth2Handler) Authenticate(ctx *gin.Context) {
	var request oauth2.OAuth2Info
	err := json.NewDecoder(ctx.Request.Body).Decode(&request)
	if err != nil {
		if h.Error != nil {
			h.Error(ctx.Request.Context(), "cannot decode OAuth2Info model: "+err.Error())
		}
		ctx.String(http.StatusBadRequest, "cannot decode OAuth2Info model")
		return
	}
	var authorization string
	if len(ctx.Request.Header["Authorization"]) < 1 {
		authorization = ""
	} else {
		authorization = ctx.Request.Header["Authorization"][0]
	}
	ip := getRemoteIp(ctx.Request)
	var ctx2 context.Context
	ctx2 = ctx.Request.Context()
	if len(h.Config.Ip) > 0 {
		ctx2 = context.WithValue(ctx2, h.Config.Ip, ip)
		ctx.Request = ctx.Request.WithContext(ctx2)
	}
	result, err := h.OAuth2Service.Authenticate(ctx.Request.Context(), &request, authorization)
	if err != nil {
		result.Status = h.SystemError
		if h.Error != nil {
			h.Error(ctx.Request.Context(), err.Error())
		}
		respond(ctx, http.StatusOK, result, h.Log, h.Config.Resource, h.Config.Authenticate, false, err.Error())
	} else {
		respond(ctx, http.StatusOK, result, h.Log, h.Config.Resource, h.Config.Authenticate, true, "")
	}
}

func respond(ctx *gin.Context, code int, result interface{}, writeLog func(context.Context, string, string, bool, string) error, resource string, action string, success bool, desc string) {
	ctx.JSON(code, result)
	if writeLog != nil {
		writeLog(ctx.Request.Context(), resource, action, success, desc)
	}
}
func getRemoteIp(r *http.Request) string {
	remoteIP, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		remoteIP = r.RemoteAddr
	}
	return remoteIP
}
