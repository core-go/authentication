package auth

type TokenWhitelistChecker interface {
	Add(id string, token string) error
	Check(id string, token string) bool
}
