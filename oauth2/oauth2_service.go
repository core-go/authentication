package oauth2

import (
	"context"
	auth "github.com/core-go/authentication"
)

type OAuth2Service interface {
	Configurations(ctx context.Context) ([]Configuration, error)
	Configuration(ctx context.Context, id string) (*Configuration, error)
	Authenticate(ctx context.Context, auth *OAuth2Info, authorization string) (auth.AuthResult, error)
}
