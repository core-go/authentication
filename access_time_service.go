package auth

import "context"

type AccessTimeService interface {
	Load(ctx context.Context, id string) (*AccessTime, error)
}
