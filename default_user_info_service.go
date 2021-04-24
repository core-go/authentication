package auth

import "context"

type DefaultUserInfoService struct {
	MaxPasswordAge int64
	*BaseUserInfoService
}

func NewUserInfoService(authenticationRepository AuthenticationRepository, maxPasswordAge int64, maxPasswordFailed int, lockedMinutes int) *DefaultUserInfoService {
	b := NewBaseUserInfoService(authenticationRepository, maxPasswordFailed, lockedMinutes)
	return &DefaultUserInfoService{maxPasswordAge, b}
}

func (s *DefaultUserInfoService) GetUserInfo(ctx context.Context, info AuthInfo) (*UserInfo, error) {
	userInfo, err := s.AuthenticationRepository.GetUserInfo(ctx, info.Username)
	if err != nil {
		return nil, err
	}

	if s.MaxPasswordAge > 0 && userInfo.MaxPasswordAge <= 0 {
		userInfo.MaxPasswordAge = s.MaxPasswordAge
	}
	return userInfo, nil
}
