package auth

import (
	"context"
	"time"
)

type CodeSender interface {
	Send(ctx context.Context, to string, code string, expireAt time.Time, params interface{}) error
}
