package mock

import (
	"context"
	"github.com/common-go/auth"
	l "github.com/common-go/auth/ldap"
	"strings"
)

type MockLDAPAuthenticator struct {
	Config    l.LDAPConfig
	Service   *l.LDAPAuthenticator
	Usernames []string
	Status    auth.Status
}

func NewDAPAuthenticatorByConfig(conf l.LDAPConfig, status auth.Status) (auth.Authenticator, error) {
	s := conf.Users
	if len(s) > 0 {
		users := strings.Split(conf.Users, ",")
		return NewMockLDAPAuthenticator(conf, users, status)
	} else {
		return l.NewLDAPAuthenticator(conf, status)
	}
}

func NewMockLDAPAuthenticator(ldapConfig l.LDAPConfig, userNames []string, status auth.Status) (*MockLDAPAuthenticator, error) {
	s, err := l.NewLDAPAuthenticator(ldapConfig, status)
	if err != nil {
		return nil, err
	}
	basicAuthenticator := &MockLDAPAuthenticator{
		Config:    ldapConfig,
		Service:   s,
		Usernames: userNames,
		Status:    status,
	}
	return basicAuthenticator, nil
}

func (s *MockLDAPAuthenticator) Authenticate(ctx context.Context, info auth.AuthInfo) (auth.AuthResult, error) {
	username := info.Username

	for _, x := range s.Usernames {
		if username == x {
			result := auth.AuthResult{}
			result.Status = s.Status.Success
			account := auth.UserAccount{}
			account.Id = username
			account.DisplayName = username
			account.Contact = "admin@gmail.com"
			result.User = &account
			return result, nil
		}
	}
	return s.Service.Authenticate(ctx, info)
}
