package auth

import "context"

type Authenticator interface {
	Authenticate(ctx context.Context, user AuthInfo) (AuthResult, error)
}
