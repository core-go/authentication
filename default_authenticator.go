package auth

import (
	"context"
	"errors"
	"log"
	"strings"
	"time"
)

type DefaultAuthenticator struct {
	Status             Status
	PayloadConfig      PayloadConfig
	UserInfoService    UserInfoService
	Check              func(ctx context.Context, user AuthInfo) (AuthResult, error)
	PasswordComparator ValueComparator
	Privileges         func(ctx context.Context, id string) ([]Privilege, error)
	GenerateToken      func(payload interface{}, secret string, expiresIn int64) (string, error)
	TokenConfig        TokenConfig
	CodeExpires        int
	CodeService        CodeService
	SendCode           func(ctx context.Context, to string, code string, expireAt time.Time, params interface{}) error
	GenerateCode       func() string
}

func NewBasicAuthenticator(status Status, check func(context.Context, AuthInfo) (AuthResult, error), userInfoService UserInfoService, generateToken func(interface{}, string, int64) (string, error), tokenConfig TokenConfig, payloadConfig PayloadConfig, options ...func(context.Context, string) ([]Privilege, error)) *DefaultAuthenticator {
	var loadPrivileges func(context.Context, string) ([]Privilege, error)
	if len(options) >= 1 {
		loadPrivileges = options[0]
	}
	return NewBasicAuthenticatorWithTwoFactors(status, check, userInfoService, generateToken, tokenConfig, payloadConfig, loadPrivileges, nil, nil, 0)
}
func NewBasicAuthenticatorWithTwoFactors(status Status, check func(context.Context, AuthInfo) (AuthResult, error), userInfoService UserInfoService, generateToken func(interface{}, string, int64) (string, error), tokenConfig TokenConfig, payloadConfig PayloadConfig, loadPrivileges func(context.Context, string) ([]Privilege, error), sendCode func(context.Context, string, string, time.Time, interface{}) error, codeService CodeService, codeExpires int, options ...func() string) *DefaultAuthenticator {
	if check == nil {
		panic(errors.New("basic check cannot be nil"))
	}
	if sendCode != nil && (codeService == nil || codeExpires <= 0) {
		panic(errors.New("when using two-factor, codeService and sendCode must not be nil, and codeExpires must be greater than 0"))
	}
	var generate func() string
	if len(options) >= 1 {
		generate = options[0]
	}
	service := &DefaultAuthenticator{
		Status:          status,
		PayloadConfig:   payloadConfig,
		Check:           check,
		UserInfoService: userInfoService,
		Privileges:      loadPrivileges,
		GenerateToken:   generateToken,
		TokenConfig:     tokenConfig,
		CodeExpires:     codeExpires,
		CodeService:     codeService,
		SendCode:        sendCode,
		GenerateCode:    generate,
	}
	return service
}
func NewAuthenticator(status Status, userInfoService UserInfoService, passwordComparator ValueComparator, generateToken func(interface{}, string, int64) (string, error), tokenConfig TokenConfig, payloadConfig PayloadConfig, options ...func(context.Context, string) ([]Privilege, error)) *DefaultAuthenticator {
	var loadPrivileges func(context.Context, string) ([]Privilege, error)
	if len(options) >= 1 {
		loadPrivileges = options[0]
	}
	return NewAuthenticatorWithTwoFactors(status, userInfoService, passwordComparator, generateToken, tokenConfig, payloadConfig, loadPrivileges, nil, nil, 0, nil)
}
func NewAuthenticatorWithTwoFactors(status Status, userInfoService UserInfoService, passwordComparator ValueComparator, generateToken func(interface{}, string, int64) (string, error), tokenConfig TokenConfig, payloadConfig PayloadConfig, loadPrivileges func(context.Context, string) ([]Privilege, error), sendCode func(context.Context, string, string, time.Time, interface{}) error, codeService CodeService, codeExpires int, options ...func() string) *DefaultAuthenticator {
	if passwordComparator == nil {
		panic(errors.New("password comparator cannot be nil"))
	}
	if sendCode != nil && (codeService == nil || codeExpires <= 0) {
		panic(errors.New("when using two-factor, codeService and sendCode must not be nil, and codeExpires must be greater than 0"))
	}
	var generate func() string
	if len(options) >= 1 {
		generate = options[0]
	}
	service := &DefaultAuthenticator{
		Status:             status,
		PayloadConfig:      payloadConfig,
		Check:              nil,
		UserInfoService:    userInfoService,
		PasswordComparator: passwordComparator,
		Privileges:         loadPrivileges,
		GenerateToken:      generateToken,
		TokenConfig:        tokenConfig,
		CodeExpires:        codeExpires,
		CodeService:        codeService,
		SendCode:           sendCode,
		GenerateCode:       generate,
	}
	return service
}

func (s *DefaultAuthenticator) Authenticate(ctx context.Context, info AuthInfo) (AuthResult, error) {
	result := AuthResult{Status: s.Status.Fail}

	username := info.Username
	password := info.Password

	if len(strings.TrimSpace(username)) == 0 && len(strings.TrimSpace(password)) == 0 || (info.Step > 0 && len(info.Passcode) == 0) {
		return result, nil
	}

	if s.Check != nil && info.Step <= 0 {
		var er0 error
		result, er0 = s.Check(ctx, info)
		if er0 != nil || result.Status != s.Status.Success && result.Status != s.Status.SuccessAndReactivated {
			return result, er0
		}
		if s.UserInfoService == nil {
			var tokenExpiredTime = time.Now().Add(time.Second * time.Duration(int(s.TokenConfig.Expires/1000)))
			var payload map[string]interface{}
			if result.User == nil {
				payload = make(map[string]interface{})
				if len(s.PayloadConfig.UserId) > 0 {
					payload[s.PayloadConfig.UserId] = info.Username
				}
				if len(s.PayloadConfig.Username) > 0 {
					payload[s.PayloadConfig.Username] = info.Username
				}
			} else {
				u := result.User
				payload = UserAccountToPayload(ctx, u, s.PayloadConfig)
			}
			token, er4 := s.GenerateToken(payload, s.TokenConfig.Secret, s.TokenConfig.Expires)
			if er4 != nil {
				return result, er4
			}
			account := UserAccount{}
			account.Token = token
			result.Status = s.Status.Success
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
		result.Status = s.Status.NotFound
		return result, er1
	}

	if s.Check == nil && info.Step <= 0 {
		validPassword, er2 := s.PasswordComparator.Compare(password, user.Password)
		if er2 != nil {
			return result, er2
		}
		if !validPassword {
			er3 := s.UserInfoService.Fail(ctx, *user)
			if er3 != nil {
				return result, er3
			}
			result.Status = s.Status.WrongPassword
			return result, nil
		}
		account := UserAccount{}
		result.User = &account
	}
	if user.Disable {
		result.Status = s.Status.Disabled
		return result, nil
	}

	if user.Suspended {
		result.Status = s.Status.Suspended
		return result, nil
	}

	locked := user.LockedUntilTime != nil && (compareDate(time.Now(), *user.LockedUntilTime) < 0)
	if locked {
		result.Status = s.Status.Locked
		return result, nil
	}

	var passwordExpiredTime *time.Time = nil // date.addDays(time.Now(), 10)
	if user.PasswordChangedTime != nil && user.MaxPasswordAge != 0 {
		t := addDays(*user.PasswordChangedTime, user.MaxPasswordAge)
		passwordExpiredTime = &t
	}
	if passwordExpiredTime != nil && compareDate(time.Now(), *passwordExpiredTime) > 0 {
		result.Status = s.Status.PasswordExpired
		return result, nil
	}

	if !IsAccessDateValid(user.AccessDateFrom, user.AccessDateTo) {
		result.Status = s.Status.Disabled
		return result, nil
	}
	if !IsAccessTimeValid(user.AccessTimeFrom, user.AccessTimeTo) {
		result.Status = s.Status.AccessTimeLocked
		return result, nil
	}

	if user.TwoFactors {
		userId := user.Id
		if info.Step <= 0 {
			var codeSend string
			if s.GenerateCode != nil {
				codeSend = s.GenerateCode()
			} else {
				codeSend = Generate(6)
			}

			codeSave, er0 := s.PasswordComparator.Hash(codeSend)
			if er0 != nil {
				return result, er0
			}
			expiredAt := addSeconds(time.Now(), s.CodeExpires)
			count, er1 := s.CodeService.Save(ctx, userId, codeSave, expiredAt)
			if count > 0 && er1 == nil {
				er3 := s.SendCode(ctx, username, codeSend, expiredAt, user.Contact)
				if er3 != nil {
					return result, er3
				}
				result.Status = s.Status.TwoFactorRequired
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
	payload := ToPayload(ctx, user, s.PayloadConfig)
	//payload := StoredUser{Id: user.Id, Username: user.Username, Contact: user.Contact, Type: user.Type, Roles: user.Roles, Privileges: user.Privileges}
	token, er4 := s.GenerateToken(payload, s.TokenConfig.Secret, jwtTokenExpires)
	if er4 != nil {
		return result, er4
	}
	if user.Deactivated == true {
		result.Status = s.Status.SuccessAndReactivated
	} else {
		result.Status = s.Status.Success
	}

	account := mapUserInfoToUserAccount(*user)
	account.Token = token
	account.TokenExpiredTime = &tokenExpiredTime
	if s.Privileges != nil {
		privileges, er5 := s.Privileges(ctx, user.Id)
		if er5 != nil {
			return result, er5
		}
		if privileges != nil && len(privileges) > 0 {
			account.Privileges = privileges
		}
	}
	result.User = &account
	er6 := s.UserInfoService.Pass(ctx, *user)
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
	account.Id = user.Id
	account.Username = user.Username
	account.Type = user.UserType
	if user.Roles != nil {
		account.Roles = user.Roles
	}
	if len(user.Id) > 0 {
		account.Id = user.Id
	}
	if len(user.DisplayName) > 0 {
		account.DisplayName = user.DisplayName
	}
	if len(user.Contact) > 0 {
		account.Contact = user.Contact
	}
	if len(user.Email) > 0 {
		account.Email = user.Email
	}
	if len(user.Phone) > 0 {
		account.Phone = user.Phone
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
	if len(user.ImageURL) > 0 {
		account.ImageURL = user.ImageURL
	}
	if len(user.Gender) > 0 {
		account.Gender = user.Gender
	}
	return account
}
func FromContext(ctx context.Context, key string) string {
	u := ctx.Value(key)
	if u == nil {
		return ""
	}
	v, ok := u.(string)
	if !ok {
		return ""
	}
	return v
}
func UserAccountToPayload(ctx context.Context, u *UserAccount, s PayloadConfig) map[string]interface{} {
	payload := make(map[string]interface{})
	if len(s.Ip) > 0 {
		ip := FromContext(ctx, s.Ip)
		if len(ip) > 0 {
			payload[s.Ip] = ip
		}
	}
	if u == nil {
		return payload
	}
	if s.UserId != "" {
		payload[s.UserId] = u.Id
	}
	if s.Username != "" {
		payload[s.Username] = u.Username
	}
	if s.Contact != "" && len(u.Contact) > 0 {
		payload[s.Contact] = u.Contact
		u.Contact = ""
	}
	if s.UserType != "" && len(u.Type) > 0 {
		payload[s.UserType] = u.Type
		u.Type = ""
	}
	if s.Roles != "" && u.Roles != nil && len(u.Roles) > 0 {
		payload[s.Roles] = u.Roles
		u.Roles = nil
	}
	return payload
}
func ToPayload(ctx context.Context, user *UserInfo, s PayloadConfig) map[string]interface{} {
	payload := make(map[string]interface{})
	if len(s.Ip) > 0 {
		ip := FromContext(ctx, s.Ip)
		if len(ip) > 0 {
			payload[s.Ip] = ip
		}
	}
	if user == nil {
		return payload
	}
	if len(s.UserId) > 0 && len(user.Id) > 0 {
		payload[s.UserId] = user.Id
	}
	if len(s.Username) > 0 && len(user.Username) > 0 {
		payload[s.Username] = user.Username
	}
	if len(s.Contact) > 0 && len(user.Contact) > 0 {
		payload[s.Contact] = user.Contact
		user.Contact = ""
	}
	if len(s.Email) > 0 && len(user.Email) > 0 {
		payload[s.Email] = user.Email
		user.Email = ""
	}
	if len(s.Phone) > 0 && len(user.Phone) > 0 {
		payload[s.Phone] = user.Phone
		user.Phone = ""
	}
	if len(s.UserType) > 0 && len(user.UserType) > 0 {
		payload[s.UserType] = user.UserType
		user.UserType = ""
	}
	if len(s.Roles) > 0 && user.Roles != nil && len(user.Roles) > 0 {
		payload[s.Roles] = user.Roles
		user.Roles = nil
	}
	if len(s.Privileges) > 0 && user.Privileges != nil && len(user.Privileges) > 0 {
		payload[s.Roles] = user.Privileges
		user.Privileges = nil
	}
	return payload
}
