package auth

import "context"

type PrivilegesReader interface {
	Privileges(ctx context.Context) ([]Privilege, error)
}
