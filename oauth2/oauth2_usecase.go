package oauth2

import (
	"context"
	auth "github.com/core-go/authentication"
	"strings"
)

type OAuth2UseCase struct {
	Status                  auth.Status
	OAuth2UserRepositories  map[string]OAuth2UserRepository
	UserRepositories        map[string]UserRepository
	ConfigurationRepository ConfigurationRepository
	Generate                func(ctx context.Context) (string, error)
	TokenService            TokenPort
	TokenConfig             auth.TokenConfig
	PayloadConfig           auth.PayloadConfig
	Privileges              func(ctx context.Context, id string) ([]auth.Privilege, error)
	AccessTime              func(ctx context.Context, id string) (*auth.AccessTime, error)
}

func NewOAuth2Service(status auth.Status, oauth2UserRepositories map[string]OAuth2UserRepository, userRepositories map[string]UserRepository, configurationRepository ConfigurationRepository, generate func(context.Context) (string, error), tokenService TokenPort, tokenConfig auth.TokenConfig, privileges func(context.Context, string) ([]auth.Privilege, error), options ...func(context.Context, string) (*auth.AccessTime, error)) *OAuth2UseCase {
	if generate == nil {
		panic("Generate cannot be nil")
	}
	var loadAccessTime func(context.Context, string) (*auth.AccessTime, error)
	if len(options) >= 1 {
		loadAccessTime = options[0]
	}
	return &OAuth2UseCase{
		Status:                  status,
		OAuth2UserRepositories:  oauth2UserRepositories,
		UserRepositories:        userRepositories,
		ConfigurationRepository: configurationRepository,
		Generate:                generate,
		TokenService:            tokenService,
		TokenConfig:             tokenConfig,
		Privileges:              privileges,
		AccessTime:              loadAccessTime,
	}
}
func (s *OAuth2UseCase) Configurations(ctx context.Context) ([]Configuration, error) {
	models, err := s.ConfigurationRepository.GetConfigurations(ctx)
	return models, err
}
func (s *OAuth2UseCase) Configuration(ctx context.Context, id string) (*Configuration, error) {
	model, _, err := s.ConfigurationRepository.GetConfiguration(ctx, id)
	return model, err
}

func (s *OAuth2UseCase) Authenticate(ctx context.Context, info *OAuth2Info, authorization string) (auth.AuthResult, error) {
	result := auth.AuthResult{Status: s.Status.Fail}
	var linkUserId = ""
	if info.Link {
		if len(authorization) == 0 {
			linkUserId = ""
		} else {
			if strings.HasPrefix(authorization, "Bearer ") != true {
				return result, nil
			}
			token := authorization[7:]
			_, _, _, er0 := s.TokenService.VerifyToken(token, s.TokenConfig.Secret)
			if er0 != nil {
				result.Status = s.Status.Error
				return result, er0
			}
			linkUserId = s.getStringValue(token, "userId") // TODO
		}
	}
	integrations, clientId, er1 := s.ConfigurationRepository.GetConfiguration(ctx, info.Id)
	if er1 != nil {
		return result, er1
	}

	if len(integrations.ClientId) > 0 {
		if len(info.Id) == 0 {
			return result, nil
		}
		integrations.ClientId = clientId
		return s.processAccount(ctx, info, *integrations, linkUserId)
	}
	return result, nil
}
func (s *OAuth2UseCase) getStringValue(tokenData interface{}, field string) string {
	if authorizationToken, ok := tokenData.(map[string]interface{}); ok {
		value, _ := authorizationToken[field].(string)
		return value
	}
	return ""
}
func (s *OAuth2UseCase) buildResult(ctx context.Context, id, email, displayName string, sourceType string, accessToken string, newUser bool) (auth.AuthResult, error) {
	user := auth.AccessTime{}
	result := auth.AuthResult{Status: s.Status.Error}
	if s.AccessTime != nil {
		accessTime, er1 := s.AccessTime(ctx, id)
		if er1 != nil {
			return result, er1
		}
		if accessTime != nil {
			user = *accessTime
			if !auth.IsAccessDateValid(accessTime.AccessDateFrom, accessTime.AccessDateTo) {
				result := auth.AuthResult{Status: s.Status.Disabled}
				return result, nil
			}
			if !auth.IsAccessTimeValid(accessTime.AccessTimeFrom, accessTime.AccessTimeTo) {
				result := auth.AuthResult{Status: s.Status.AccessTimeLocked}
				return result, nil
			}
		}
	}

	tokenExpiredTime, jwtTokenExpires := auth.SetTokenExpiredTime(user.AccessTimeFrom, user.AccessTimeTo, s.TokenConfig.Expires)
	payload := BuildPayload(id, email, s.PayloadConfig)
	var tokens map[string]string
	if len(s.PayloadConfig.Tokens) > 0 {
		tokens = make(map[string]string)
		tokens[sourceType] = accessToken
		payload[s.PayloadConfig.Tokens] = tokens
	}
	token, er2 := s.TokenService.GenerateToken(payload, s.TokenConfig.Secret, jwtTokenExpires)

	if er2 != nil {
		return result, er2
	}
	var account auth.UserAccount
	account.Username = email
	account.Id = id
	account.Contact = &email
	account.DisplayName = &displayName
	account.Token = token
	account.TokenExpiredTime = &tokenExpiredTime

	if s.Privileges != nil {
		privileges, er1 := s.Privileges(ctx, id)
		if er1 != nil {
			return result, er1
		}
		account.Privileges = privileges
	}
	result.Status = s.Status.Success
	result.User = &account
	return result, nil
}
func (s *OAuth2UseCase) processAccount(ctx context.Context, data *OAuth2Info, integration Configuration, linkUserId string) (auth.AuthResult, error) {
	code := data.Code
	urlRedirect := data.RedirectUri
	clientSecret := integration.ClientSecret
	clientId := integration.ClientId
	repository := s.OAuth2UserRepositories[data.Id]
	user, accessToken, err := repository.GetUserFromOAuth2(ctx, urlRedirect, clientId, clientSecret, code)
	if err != nil || user == nil {
		result := auth.AuthResult{Status: s.Status.Error}
		return result, err
	}
	return s.checkAccount(ctx, user, accessToken, linkUserId, data.Id)
}

func (s *OAuth2UseCase) checkAccount(ctx context.Context, user *User, accessToken string, linkUserId string, types string) (auth.AuthResult, error) {
	personRepository := s.UserRepositories[types]
	eId, disable, suspended, er0 := personRepository.GetUser(ctx, user.Email) //i
	result := auth.AuthResult{Status: s.Status.Error}
	if er0 != nil {
		return result, er0
	}
	if len(linkUserId) > 0 {
		if eId != linkUserId {
			result := auth.AuthResult{Status: s.Status.Fail}
			return result, nil
		}
		ok1, er2 := personRepository.Update(ctx, linkUserId, user.Email, user.Account)
		if ok1 && er2 == nil {
			return s.buildResult(ctx, eId, user.Email, user.DisplayName, types, accessToken, false)
		}
	}
	if len(eId) != 0 {
		ok1, er2 := personRepository.Update(ctx, linkUserId, user.Email, user.Account)
		if ok1 && er2 == nil {
			return s.buildResult(ctx, eId, user.Email, user.DisplayName, types, accessToken, false)
		}
	}
	if len(eId) == 0 {
		userId, er3 := s.Generate(ctx)
		if er3 != nil {
			return result, er3
		}
		duplicate, er4 := personRepository.Insert(ctx, userId, user)
		if duplicate {
			i := 1
			for duplicate && i <= 5 {
				i++
				userId, er3 = s.Generate(ctx)
				if er3 != nil {
					return result, er3
				}
				duplicate, er4 = personRepository.Insert(ctx, userId, user)
				if er4 != nil {
					return result, er4
				}
			}
			if duplicate {
				return result, nil
			}
		}
		if er4 == nil && !duplicate {
			return s.buildResult(ctx, eId, user.Email, user.DisplayName, types, accessToken, true)
		}
		return result, er4
	}
	if disable {
		result.Status = s.Status.Disabled
		return result, nil
	}
	if suspended {
		result.Status = s.Status.Suspended
		return result, nil
	}

	ok3, er5 := personRepository.Update(ctx, eId, user.Email, user.Account)
	if ok3 && er5 == nil {
		return s.buildResult(ctx, eId, user.Email, user.Account, types, accessToken, false)
	}

	return result, nil
}
func BuildPayload(id, email string, c auth.PayloadConfig) map[string]interface{} {
	m := make(map[string]interface{})
	if len(c.Id) > 0 {
		m[c.Id] = id
	}
	if len(c.Username) > 0 {
		m[c.Username] = email
	}
	if len(c.Contact) > 0 {
		m[c.Contact] = email
	}
	return m
}
