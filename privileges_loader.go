package auth

import "context"

type PrivilegesLoader interface {
	Load(ctx context.Context, id string) ([]Privilege, error)
}
