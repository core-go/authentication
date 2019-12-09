package auth

type PrivilegeService interface {
	GetPrivileges(id string) ([]Privilege, error)
}
