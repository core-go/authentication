package auth

import "context"

type PrivilegeService interface {
	GetPrivileges(ctx context.Context, id string) ([]Privilege, error)
}
