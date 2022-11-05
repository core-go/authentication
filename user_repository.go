package auth

import (
	"context"
	"time"
)

type UserRepository interface {
	GetUser(ctx context.Context, auth AuthInfo) (*UserInfo, error)
	Pass(ctx context.Context, id string, deactivated *bool) error
	Fail(ctx context.Context, id string, failCount *int, lockedUntilTime *time.Time) error
}
