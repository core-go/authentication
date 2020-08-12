package auth

import "time"

type TokenBlacklistService interface {
	Revoke(token string, reason string, expires time.Time) error
	RevokeAllTokens(id string, reason string) error
	Check(id string, token string, createAt time.Time) string
}
