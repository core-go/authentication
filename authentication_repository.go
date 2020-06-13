package auth

import (
	"context"
	"time"
)

type AuthenticationRepository interface {
	PassAuthenticationAndActivate(ctx context.Context, userId string) (int64, error)
	PassAuthentication(ctx context.Context, userId string) (int64, error)
	WrongPassword(ctx context.Context, userId string, failCount int, lockedUntil *time.Time) error

	GetUserInfo(ctx context.Context, username string) (*UserInfo, error)
}
