package auth

import (
	"context"
	"time"
)

type BaseUserInfoService struct {
	AuthenticationRepository AuthenticationRepository
	MaxPasswordFailed        int
	LockedMinutes            int
}

func NewBaseUserInfoService(authenticationRepository AuthenticationRepository, maxPasswordFailed int, lockedMinutes int) *BaseUserInfoService{
	return &BaseUserInfoService{authenticationRepository, maxPasswordFailed, lockedMinutes}
}

func (s *BaseUserInfoService) Pass(ctx context.Context, user UserInfo) error {
	if s.AuthenticationRepository == nil {
		return nil
	}
	if user.Deactivated == true {
		_, er1 := s.AuthenticationRepository.PassAndActivate(ctx, user.UserId)
		return er1
	}
	_, er2 := s.AuthenticationRepository.Pass(ctx, user.UserId)
	return er2
}

func (s *BaseUserInfoService) Fail(ctx context.Context, user UserInfo) error {
	if s.AuthenticationRepository == nil {
		return nil
	}
	if s.LockedMinutes > 0 && s.MaxPasswordFailed > 0 && user.FailCount >= s.MaxPasswordFailed {
		lockedUntilTime := time.Now().Add(time.Minute * time.Duration(s.LockedMinutes))
		return s.AuthenticationRepository.Fail(ctx, user.UserId, 0, &lockedUntilTime)
	}
	count := user.FailCount + 1
	return s.AuthenticationRepository.Fail(ctx, user.UserId, count, nil)
}
