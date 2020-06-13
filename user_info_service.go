package auth

import "context"

type UserInfoService interface {
	GetUserInfo(ctx context.Context, auth AuthInfo) (*UserInfo, error)
	PassAuthentication(ctx context.Context, user UserInfo) error
	HandleWrongPassword(ctx context.Context, user UserInfo) error
}
