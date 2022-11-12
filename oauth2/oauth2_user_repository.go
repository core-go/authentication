package oauth2

import "context"

type OAuth2UserRepository interface {
	GetUserFromOAuth2(ctx context.Context, urlRedirect string, clientId string, clientSecret string, code string) (*User, string, error)
	GetRequestTokenOAuth(ctx context.Context, key string, secret string) (string, error)
}
