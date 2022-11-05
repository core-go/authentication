package ldap

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/core-go/auth"
	"gopkg.in/ldap.v3"
	"net"
	"strconv"
	"strings"
	"time"
)

type LDAPAuthenticator struct {
	Config LDAPConfig
	Domain string
	Status auth.Status
}

func GetDomain(baseDN string) (string, error) {
	d := ""
	x := strings.Split(strings.ToLower(baseDN), ",")
	for _, s := range x {
		y := strings.TrimSpace(s)
		if strings.HasPrefix(y, "dc=") {
			d = d + "." + y[3:]
		}
	}
	if len(d) <= 1 {
		return "", errors.New("invalid base DN")
	}
	return d[1:], nil
}
func NewLDAPAuthenticator(ldapConfig LDAPConfig, status auth.Status) (*LDAPAuthenticator, error) {
	domain := strings.TrimSpace(ldapConfig.Domain)
	var err error
	if len(domain) <= 0 {
		domain, err = GetDomain(ldapConfig.BaseDN)
		if err != nil {
			return nil, err
		}
	}
	return &LDAPAuthenticator{Config: ldapConfig, Domain: domain, Status: status}, nil
}
func NewConn(c LDAPConfig) (*ldap.Conn, error) {
	var l *ldap.Conn
	var err error
	if c.Timeout > 0 {
		ldap.DefaultTimeout = time.Duration(c.Timeout) * time.Millisecond
	}
	if c.TLS != nil && *c.TLS {
		if c.InsecureSkipVerify != nil && *c.InsecureSkipVerify {
			l, err = ldap.DialTLS("tcp", c.Server, &tls.Config{ServerName: c.Server, InsecureSkipVerify: true})
		} else {
			l, err = ldap.DialTLS("tcp", c.Server, &tls.Config{ServerName: c.Server})
		}
	} else {
		l, err = ldap.Dial("tcp", c.Server)
		if err == nil {
			if c.StartTLS != nil && *c.StartTLS {
				if c.InsecureSkipVerify != nil && *c.InsecureSkipVerify {
					err = l.StartTLS(&tls.Config{ServerName: c.Server, InsecureSkipVerify: true})
				} else {
					err = l.StartTLS(&tls.Config{ServerName: c.Server})
				}
			}
		}
	}
	return l, err
}
func (s *LDAPAuthenticator) Authenticate(ctx context.Context, info auth.AuthInfo) (auth.AuthResult, error) {
	result := auth.AuthResult{}
	account := auth.UserAccount{}
	result.Status = s.Status.Fail
	l, er1 := NewConn(s.Config)
	if er1 != nil {
		if e, ok0 := er1.(*ldap.Error); ok0 {
			e2 := e.Err
			if e3, ok2 := e2.(*net.OpError); ok2 {
				e4 := e3.Err
				if e4 != nil && e4.Error() == "i/o timeout" {
					result.Status = s.Status.Timeout
					return result, e4
				}
			}
			return result, e2
		}
		return result, er1
	}
	defer l.Close()
	username := info.Username
	if len(s.Domain) > 0 && strings.Index(username, "@") < 0 {
		username = info.Username + "@" + s.Domain
	}
	er2 := l.Bind(username, info.Password)
	if er2 != nil {
		if e, ok := er2.(*ldap.Error); ok {
			if e.ResultCode == ldap.LDAPResultInvalidCredentials {
				return result, nil
			}
		}
		return result, er2
	}
	result.Status = s.Status.Success
	if len(s.Config.Filter) == 0 || (len(s.Config.Id) == 0 && len(s.Config.DisplayName) == 0 && len(s.Config.Contact) == 0) {
		account.Id = info.Username
		result.User = &account
		return result, nil
	}
	filters := make([]string, 0)
	if len(s.Config.Id) > 0 {
		filters = append(filters, s.Config.Id)
	}
	if len(s.Config.DisplayName) > 0 {
		filters = append(filters, s.Config.DisplayName)
	}
	if len(s.Config.Contact) > 0 {
		filters = append(filters, s.Config.Contact)
	}
	if len(filters) > 0 {
		x := fmt.Sprintf("(&(%s=%s))", s.Config.Filter, info.Username)
		searchRequest := ldap.NewSearchRequest(
			s.Config.BaseDN,
			ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 1, 0, false,
			x,
			filters,
			nil,
		)
		sr, er3 := l.Search(searchRequest)
		if er3 != nil {
			account.Id = info.Username
			result.User = &account
			return result, er3
		}
		if len(sr.Entries) >= 1 {
			entry := sr.Entries[0]
			if len(s.Config.Id) > 0 {
				account.Id = entry.GetAttributeValue(s.Config.Id)
			}
			if len(s.Config.DisplayName) > 0 {
				v := entry.GetAttributeValue(s.Config.DisplayName)
				account.DisplayName = &v
			}
			if len(s.Config.Contact) > 0 {
				v := entry.GetAttributeValue(s.Config.Contact)
				account.Contact = &v
			}
			if len(s.Config.Email) > 0 {
				v := entry.GetAttributeValue(s.Config.Email)
				account.Email = &v
			}
			if len(s.Config.Phone) > 0 {
				v := entry.GetAttributeValue(s.Config.Phone)
				account.Phone = &v
			}
		}
		result.User = &account
	}
	return result, nil
}

const u = 11644473600

func ToDate(ldap string) *time.Time {
	if ldap == "9223372036854775807" {
		return nil
	}
	i, er := strconv.ParseInt(ldap, 10, 64)
	if er != nil {
		return nil
	}
	l := i / 10000000
	x := time.Unix(l-u, 0)
	return &x
}
