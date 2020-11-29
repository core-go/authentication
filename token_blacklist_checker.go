package auth

import "time"

type TokenBlacklistChecker interface {
	Revoke(token string, reason string, expires time.Time) error
	RevokeAllTokens(id string, reason string) error
	Check(id string, token string, createAt time.Time) string
}
