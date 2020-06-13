package auth

import "context"

type TwoFactorAuthenticationRepository interface {
	Require(ctx context.Context, id string) (bool, error)
}
