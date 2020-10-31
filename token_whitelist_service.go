package auth

type TokenWhitelistService interface {
	Add(token, secret, reason string) error
	Check(id string, token string) bool
}
