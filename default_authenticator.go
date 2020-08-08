package auth

import (
	"context"
	"errors"
	"log"
	"strings"
	"time"
)

type DefaultAuthenticator struct {
	UserInfoService    UserInfoService
	BasicAuthenticator Authenticator
	PasswordComparator ValueComparator
	PrivilegeService   PrivilegeService
	TokenGenerator     TokenGenerator
	TokenConfig        TokenConfig
	CodeExpires        int
	CodeService        CodeService
	CodeSender         CodeSender
	Generator          CodeGenerator
}

func NewBasicAuthenticator(basicAuthenticator Authenticator, userInfoService UserInfoService, privilegeService PrivilegeService, tokenGenerator TokenGenerator, tokenConfig TokenConfig, isUsingTwoFactor bool, codeExpires int, codeService CodeService, codeSender CodeSender, generator CodeGenerator) *DefaultAuthenticator {
	if basicAuthenticator == nil {
		panic(errors.New("basic authenticator cannot be nil"))
	}
	if isUsingTwoFactor && (codeService == nil || codeSender == nil || codeExpires <= 0) {
		panic(errors.New("when using two-factor, codeService and codeSender must not be nil, and codeExpires must be greater than 0"))
	}
	service := &DefaultAuthenticator{
		BasicAuthenticator: basicAuthenticator,
		UserInfoService:    userInfoService,
		PrivilegeService:   privilegeService,
		TokenGenerator:     tokenGenerator,
		TokenConfig:        tokenConfig,
		CodeExpires:        codeExpires,
		CodeService:        codeService,
		CodeSender:         codeSender,
		Generator:          generator,
	}
	return service
}

func NewDefaultAuthenticator(userInfoService UserInfoService, passwordComparator ValueComparator, privilegeService PrivilegeService, tokenGenerator TokenGenerator, tokenConfig TokenConfig, isUsingTwoFactor bool, codeExpires int, codeService CodeService, codeSender CodeSender, generator CodeGenerator) *DefaultAuthenticator {
	if passwordComparator == nil {
		panic(errors.New("password comparator cannot be nil"))
	}
	if isUsingTwoFactor && (codeService == nil || codeSender == nil || codeExpires <= 0) {
		panic(errors.New("when using two-factor, codeService and codeSender must not be nil, and codeExpires must be greater than 0"))
	}
	service := &DefaultAuthenticator{
		BasicAuthenticator: nil,
		UserInfoService:    userInfoService,
		PasswordComparator: passwordComparator,
		PrivilegeService:   privilegeService,
		TokenGenerator:     tokenGenerator,
		TokenConfig:        tokenConfig,
		CodeExpires:        codeExpires,
		CodeService:        codeService,
		CodeSender:         codeSender,
		Generator:          generator,
	}
	return service
}

func (s *DefaultAuthenticator) Authenticate(ctx context.Context, info AuthInfo) (AuthResult, error) {
	result := AuthResult{Status: StatusFail}

	username := info.Username
	password := info.Password

	if len(strings.TrimSpace(username)) == 0 && len(strings.TrimSpace(password)) == 0 || (info.Step > 0 && len(info.Passcode) == 0) {
		return result, nil
	}

	if s.BasicAuthenticator != nil && info.Step <= 0 {
		var er0 error
		result, er0 = s.BasicAuthenticator.Authenticate(ctx, info)
		if er0 != nil || result.Status != StatusSuccess && result.Status != StatusSuccessAndReactivated {
			return result, er0
		}
		if s.UserInfoService == nil {
			var tokenExpiredTime = time.Now().Add(time.Second * time.Duration(int(s.TokenConfig.Expires/1000)))
			var payload StoredUser
			if result.User == nil {
				payload = StoredUser{UserId: info.Username, Username: info.Username, Contact: "", UserType: ""}
			} else {
				u := result.User
				payload = StoredUser{UserId: u.UserId, Username: u.Username, Contact: u.Contact, UserType: u.UserType, Roles: u.Roles}
			}
			token, er4 := s.TokenGenerator.GenerateToken(payload, s.TokenConfig.Secret, s.TokenConfig.Expires)
			if er4 != nil {
				return result, er4
			}
			account := UserAccount{}
			account.Token = token
			result.Status = StatusSuccess
			result.User = &account
			account.TokenExpiredTime = &tokenExpiredTime
			return result, nil
		}
	}

	user, er1 := s.UserInfoService.GetUserInfo(ctx, info)
	if er1 != nil {
		return result, er1
	}
	if user == nil {
		result.Message = "User does not exist"
		return result, er1
	}

	if s.BasicAuthenticator == nil && info.Step <= 0 {
		validPassword, er2 := s.PasswordComparator.Compare(password, user.Password)
		if er2 != nil {
			return result, er2
		}
		if !validPassword {
			er3 := s.UserInfoService.HandleWrongPassword(ctx, *user)
			if er3 != nil {
				return result, er3
			}
			result.Status = StatusWrongPassword
			return result, nil
		}
		account := UserAccount{}
		result.User = &account
	}
	if user.Disable {
		result.Status = StatusDisabled
		return result, nil
	}

	if user.Suspended {
		result.Status = StatusSuspended
		return result, nil
	}

	locked := user.LockedUntilTime != nil && (compareDate(time.Now(), *user.LockedUntilTime) < 0)
	if locked {
		result.Status = StatusLocked
		return result, nil
	}

	var passwordExpiredTime *time.Time = nil // date.addDays(time.Now(), 10)
	if user.PasswordChangedTime != nil && user.MaxPasswordAge != 0 {
		t := addDays(*user.PasswordChangedTime, user.MaxPasswordAge)
		passwordExpiredTime = &t
	}
	if passwordExpiredTime != nil && compareDate(time.Now(), *passwordExpiredTime) > 0 {
		result.Status = StatusPasswordExpired
		return result, nil
	}

	if !IsAccessDateValid(user.AccessDateFrom, user.AccessDateTo) {
		result.Status = StatusDisabled
		return result, nil
	}
	if !IsAccessTimeValid(user.AccessTimeFrom, user.AccessTimeTo) {
		result.Status = StatusAccessTimeLocked
		return result, nil
	}

	if user.TwoFactors {
		userId := user.UserId
		if info.Step <= 0 {
			var codeSend string
			if s.Generator != nil {
				codeSend = s.Generator.Generate()
			} else {
				codeSend = generate(6)
			}

			codeSave, er0 := s.PasswordComparator.Hash(codeSend)
			if er0 != nil {
				return result, er0
			}
			expiredAt := addSeconds(time.Now(), s.CodeExpires)
			count, er1 := s.CodeService.Save(ctx, userId, codeSave, expiredAt)
			if count > 0 && er1 == nil {
				er3 := s.CodeSender.Send(ctx, username, codeSend, expiredAt, user.Contact)
				if er3 != nil {
					return result, er3
				}
				result.Status = StatusTwoFactorRequired
				return result, nil
			}
		}
		code, expiredAt, er4 := s.CodeService.Load(ctx, userId)
		if er4 != nil || len(code) == 0 {
			return result, er4
		}
		if compareDate(expiredAt, time.Now()) < 0 {
			deleteCode(ctx, s.CodeService, userId)
			return result, nil
		}
		valid, er5 := s.PasswordComparator.Compare(info.Passcode, code)
		if er5 == nil {
			deleteCode(ctx, s.CodeService, userId)
		}
		if !valid || er5 != nil {
			return result, er5
		}
	}

	tokenExpiredTime, jwtTokenExpires := SetTokenExpiredTime(user.AccessTimeFrom, user.AccessTimeTo, s.TokenConfig.Expires)
	//tokenExpiredTime, jwtTokenExpires := s.setTokenExpiredTime(*user)
	payload := StoredUser{UserId: user.UserId, Username: user.Username, Contact: user.Contact, UserType: user.UserType, Roles: user.Roles, Privileges: user.Privileges}
	token, er4 := s.TokenGenerator.GenerateToken(payload, s.TokenConfig.Secret, jwtTokenExpires)
	if er4 != nil {
		return result, er4
	}
	if user.Deactivated == true {
		result.Status = StatusSuccessAndReactivated
	} else {
		result.Status = StatusSuccess
	}

	account := mapUserInfoToUserAccount(*user)
	account.Token = token
	account.TokenExpiredTime = &tokenExpiredTime
	if s.PrivilegeService != nil {
		privileges, er5 := s.PrivilegeService.GetPrivileges(ctx, user.UserId)
		if er5 != nil {
			return result, er5
		}
		account.Privileges = &privileges
	}
	result.User = &account
	er6 := s.UserInfoService.PassAuthentication(ctx, *user)
	if er6 != nil {
		return result, er6
	}
	return result, nil
}

func deleteCode(ctx context.Context, codeService CodeService, id string) {
	go func() {
		timeOut := 30 * time.Second
		ctxDelete, cancel := context.WithTimeout(context.Background(), timeOut)
		defer cancel()
		_, err := codeService.Delete(ctxDelete, id)
		if err != nil {
			log.Println(err)
		}
	}()
}

func mapUserInfoToUserAccount(user UserInfo) UserAccount {
	account := UserAccount{}
	account.UserId = user.UserId
	account.Username = user.Username
	account.UserType = user.UserType
	account.Roles = user.Roles
	if len(user.UserId) > 0 {
		account.UserId = user.UserId
	}
	if len(user.DisplayName) > 0 {
		account.DisplayName = user.DisplayName
	}
	if len(user.Contact) > 0 {
		account.Contact = user.Contact
	}

	if len(user.DateFormat) > 0 {
		account.DateFormat = user.DateFormat
	}
	if len(user.TimeFormat) > 0 {
		account.TimeFormat = user.TimeFormat
	}
	if len(user.Language) > 0 {
		account.Language = user.Language
	}
	if len(user.ImageUrl) > 0 {
		account.ImageUrl = user.ImageUrl
	}
	if len(user.Gender) > 0 {
		account.Gender = user.Gender
	}
	return account
}
