package auth

import "context"

type UserInfoService interface {
	GetUserInfo(ctx context.Context, auth AuthInfo) (*UserInfo, error)
	Pass(ctx context.Context, user UserInfo) error
	Fail(ctx context.Context, user UserInfo) error
}
