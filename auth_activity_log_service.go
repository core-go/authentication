package auth

import "context"

type AuthActivityLogService interface {
	Write(ctx context.Context, resource string, action string, success bool, desc string) error
}
