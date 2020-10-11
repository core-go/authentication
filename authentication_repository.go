package auth

import (
	"context"
	"time"
)

type AuthenticationRepository interface {
	PassAndActivate(ctx context.Context, userId string) (int64, error)
	Pass(ctx context.Context, userId string) (int64, error)
	Fail(ctx context.Context, userId string, failCount int, lockedUntil *time.Time) error

	GetUserInfo(ctx context.Context, username string) (*UserInfo, error)
}
