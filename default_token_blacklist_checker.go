package auth

import (
	"strconv"
	"strings"
	"time"
)

const joinChar = "-"

type DefaultTokenBlacklistChecker struct {
	TokenPrefix  string
	TokenExpires int64
	CacheService CacheService
}

func NewTokenBlacklistChecker(keyPrefix string, tokenExpires int64, cacheService CacheService) *DefaultTokenBlacklistChecker {
	return &DefaultTokenBlacklistChecker{keyPrefix, tokenExpires, cacheService}
}

func (b *DefaultTokenBlacklistChecker) generateKey(token string) string {
	return b.TokenPrefix + "::token::" + token
}

func (b *DefaultTokenBlacklistChecker) generateKeyForId(id string) string {
	return b.TokenPrefix + "::token::" + id
}

func (b *DefaultTokenBlacklistChecker) Revoke(token string, reason string, expiredDate time.Time) error {
	key := b.generateKey(token)
	var value string
	if len(reason) > 0 {
		value = reason
	} else {
		value = ""
	}

	today := time.Now()
	expiresInSecond := expiredDate.Sub(today)
	if expiresInSecond <= 0 {
		return nil // Token already expires, don't need add to cache
	} else {
		return b.CacheService.Put(key, value, expiresInSecond*time.Second)
	}
}

func (b *DefaultTokenBlacklistChecker) RevokeAllTokens(id string, reason string) error {
	key := b.generateKeyForId(id)
	today := time.Now()
	value := reason + joinChar + strconv.Itoa(int(today.Unix()))
	return b.CacheService.Put(key, value, time.Duration(b.TokenExpires)*time.Second)
}

func (b *DefaultTokenBlacklistChecker) Check(id string, token string, createAt time.Time) string {
	idKey := b.generateKeyForId(id)
	tokenKey := b.generateKey(token)

	keys := []string{idKey, tokenKey}
	value, _, err := b.CacheService.GetManyStrings(keys)
	if err != nil {
		return ""
	}
	if len(value[idKey]) > 0 {
		index := strings.Index(value[idKey], joinChar)
		reason := value[idKey][0:index]
		strDate := value[idKey][index+1:]
		i, err := strconv.ParseInt(strDate, 10, 64)
		if err == nil {
			tmDate := time.Unix(i, 0)
			if tmDate.Sub(createAt) > 0 {
				return reason
			}
		}
	}
	if len(value[tokenKey]) > 0 {
		return value[tokenKey]
	}
	return ""
}
