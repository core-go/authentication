package auth

import (
	"errors"
	"strconv"
	"time"
)

type DefaultTokenWhitelistChecker struct {
	Secret       string
	TokenIp      string
	TokenPrefix  string
	VerifyToken  func(tokenString string, secret string) (map[string]interface{}, int64, int64, error)
	CacheService CacheService
	Level        int
}

func NewTokenWhitelistChecker(secret string, tokenIp string, keyPrefix string, verifyToken func(tokenString string, secret string) (map[string]interface{}, int64, int64, error), cacheService CacheService, level int) *DefaultTokenWhitelistChecker {
	return &DefaultTokenWhitelistChecker{secret, tokenIp, keyPrefix, verifyToken, cacheService, level}
}

func (b *DefaultTokenWhitelistChecker) generateKey(token string) string {
	return b.TokenPrefix + "::token::" + token
}

func (b *DefaultTokenWhitelistChecker) generateKeyForId(id string) string {
	return b.TokenPrefix + "::token::" + id
}

func (b *DefaultTokenWhitelistChecker) Add(id string, token string) error {
	_, _, eta, err := b.VerifyToken(token, b.Secret)
	if err != nil {
		return err
	}
	now := time.Now()

	if eta <= now.Unix() {
		return errors.New("token expired")
	}
	expire := time.Unix(eta, 0)
	dur := expire.Sub(now)

	key := b.generateKeyForId(id)
	return b.CacheService.Put(key, token, dur)
}

func (b *DefaultTokenWhitelistChecker) Check(id string, token string) bool {
	key := b.generateKeyForId(id)

	value, err := b.CacheService.Get(key)
	if err != nil {
		return false
	}
	if value != nil {
		if tokenStore, ok := value.(string); ok {
			tokenStore, _ := strconv.Unquote(tokenStore)
			if b.Level != 0 {
				if tokenStore != token {
					return false
				}
				return true
			}

			payloadStore, _, _, err1 := b.VerifyToken(tokenStore, b.Secret)
			payload, _, _, err2 := b.VerifyToken(token, b.Secret)
			if err1 != nil || err2 != nil {
				return false
			}
			ipStore, ok1 := payloadStore[b.TokenIp]
			ip, ok2 := payload[b.TokenIp]
			if ok1 && ok2 {
				if ip == ipStore {
					return true
				}
			}
		}
	}
	return false
}
