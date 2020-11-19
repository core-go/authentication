package auth

type TokenWhitelistService interface {
	Add(token string, id string) error
	Check(id string, token string) bool
}
