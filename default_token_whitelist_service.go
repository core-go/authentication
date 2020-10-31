package auth

import (
	"errors"
	"strconv"
	"time"
)

type DefaultTokenWhitelistTokenService struct {
	TokenPrefix  string
	TokenService TokenVerifier
	CacheService CacheService
}

func NewTokenWhitelistTokenService(keyPrefix string, tokenService TokenVerifier, cacheService CacheService) *DefaultTokenWhitelistTokenService {
	return &DefaultTokenWhitelistTokenService{keyPrefix, tokenService, cacheService}
}
func (b *DefaultTokenWhitelistTokenService) generateKey(token string) string {
	return b.TokenPrefix + "::token::" + token
}

func (b *DefaultTokenWhitelistTokenService) Add(token, secret, reason string) error {
	_, _, eta, err := b.TokenService.VerifyToken(token, secret)
	if err != nil {
		return err
	}
	now := time.Now()

	if eta <= now.Unix() {
		return errors.New("token expired")
	}
	expire := time.Unix(eta, 0)
	dur := expire.Sub(now)

	key := b.generateKey(token)
	value := reason + JoinChar + strconv.Itoa(int(now.Unix()))
	return b.CacheService.Put(key, value, dur)
}

func (b *DefaultTokenWhitelistTokenService) Check(token string) bool {
	tokenKey := b.generateKey(token)

	keys := []string{tokenKey}
	value, _, err := b.CacheService.GetManyStrings(keys)
	if err != nil {
		return false
	}
	if len(value[tokenKey]) > 0 {
		//index := strings.Index(value[tokenKey], JoinChar)
		//reason := value[tokenKey][0:index]
		//strDate := value[tokenKey][index+1:]
		//
		//i, err := strconv.ParseInt(strDate, 10, 64)
		//if err == nil {
		//	tmDate := time.Unix(i, 0)
		//	if tmDate.Sub(createAt) > 0 {
		//		return reason
		//	}
		return true
	}
	return false
}
