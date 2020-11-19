package auth

import (
	"errors"
	"time"
	"strconv"
)

type DefaultTokenWhitelistTokenService struct {
	Secret       string
	TokenIp      string
	TokenPrefix  string
	TokenService TokenVerifier
	CacheService CacheService
}

func NewTokenWhitelistTokenService(secret string, tokenIp string, keyPrefix string, tokenService TokenVerifier, cacheService CacheService) *DefaultTokenWhitelistTokenService {
	return &DefaultTokenWhitelistTokenService{secret, tokenIp, keyPrefix, tokenService, cacheService}
}

func (b *DefaultTokenWhitelistTokenService) generateKey(token string) string {
	return b.TokenPrefix + "::token::" + token
}

func (b *DefaultTokenWhitelistTokenService) generateKeyForId(id string) string {
	return b.TokenPrefix + "::token::" + id
}

func (b *DefaultTokenWhitelistTokenService) Add(token string, id string) error {
	_, _, eta, err := b.TokenService.VerifyToken(token, b.Secret)
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

func (b *DefaultTokenWhitelistTokenService) Check(id string, token string) bool {
	key := b.generateKeyForId(id)

	value, err := b.CacheService.Get(key)
	if err != nil {
		return false
	}
	if value != nil {
		if tokenStore, ok := value.(string); ok {
			tokenStore, _ := strconv.Unquote(tokenStore)

			payloadStore, _, _, err1 := b.TokenService.VerifyToken(tokenStore, b.Secret)
			payload, _, _, err2 := b.TokenService.VerifyToken(token, b.Secret)
			if err1 != nil || err2 != nil {
				return false
			}
			ipStore, ok1 := payloadStore[b.TokenIp];
			ip, ok2 := payload[b.TokenIp];
			if ok1 && ok2 {
				if ip == ipStore {
					return true
				}
			}
		}
	}
	return false
}
