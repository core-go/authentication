package azure

import (
	"context"
	"crypto/rsa"
	"errors"
	"fmt"
	"strings"

	auth "github.com/core-go/authentication"
	"github.com/golang-jwt/jwt"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwk"
)

type Config struct {
	TenantId     string   `yaml:"tenant_id" mapstructure:"tenant_id"`
	ClientId     string   `yaml:"client_id" mapstructure:"client_id"`
	Scopes       []string `yaml:"scopes" mapstructure:"scopes"`
	ClientSecret string   `yaml:"client_secret" mapstructure:"client_secret"`
}

type UserRepository interface {
	Exist(ctx context.Context, id string) (bool, string, error)
	Insert(ctx context.Context, id string, user *AzureUser) (bool, error)
}

type Authenticator struct {
	GetUserByToken func(ctx context.Context, azureToken string) (*AzureUser, error)
	UserRepository UserRepository
	Privileges     func(ctx context.Context, id string) ([]auth.Privilege, error)
	GenerateToken  func(payload interface{}, secret string, expiresIn int64) (string, error)
	TokenConfig    auth.TokenConfig
	Config         Config
	Id             string
}

func NewAzureAuthenticator(
	getUserByToken func(ctx context.Context, azureToken string) (*AzureUser, error),
	userPort UserRepository,
	generateToken func(payload interface{}, secret string, expiresIn int64) (string, error),
	config Config,
	tokenConfig auth.TokenConfig,
	privileges func(ctx context.Context, id string) ([]auth.Privilege, error),
	id string,
) *Authenticator {
	if len(id) == 0 {
		id = "id"
	}
	return &Authenticator{getUserByToken, userPort, privileges, generateToken, tokenConfig, config, id}
}

const expired = "Token is expired"

// Authenticate authorization jwt here doesn't contain prefix bearer
func (a Authenticator) Authenticate(ctx context.Context, authorization string) (*auth.UserAccount, bool, error) {
	if len(authorization) == 0 {
		return nil, false, errors.New("invalid authorization")
	}
	azureToken, er1 := VerifyAzureADJWT(ctx, authorization)
	if er1 != nil {
		if strings.Contains(er1.Error(), expired) {
			return nil, true, nil
		}
		return nil, false, er1
	}

	azureID, er2 := VerifyAzureADJWTClaims(azureToken, a.Config.TenantId, a.Config.ClientId)
	if er2 != nil {
		if strings.Contains(er2.Error(), expired) {
			return nil, true, nil
		}
		return nil, false, er2
	}

	var displayName, userId string
	userId = azureID
	exist, displayName, er3 := a.UserRepository.Exist(ctx, azureID)
	if er3 != nil {
		return nil, false, er3
	}

	if !exist {
		azureUser, er4 := a.GetUserByToken(ctx, authorization)
		if er4 != nil {
			if strings.Contains(er4.Error(), expired) {
				return nil, true, nil
			}
			return nil, false, er4
		}
		displayName = azureUser.DisplayName
		userId = azureUser.Id
		ok, er5 := a.UserRepository.Insert(ctx, userId, azureUser)
		if er5 != nil {
			return nil, false, er5
		}
		if !ok {
			return nil, false, errors.New("cannot create user")
		}
	}
	account := &auth.UserAccount{
		Id:          userId,
		DisplayName: &displayName,
	}
	if a.Privileges != nil {
		privileges, er6 := a.Privileges(ctx, azureID)
		if er6 != nil {
			return nil, false, er6
		}
		account.Privileges = privileges
	}
	payload := map[string]interface{}{a.Id: azureID}
	token, er7 := a.GenerateToken(payload, a.TokenConfig.Secret, a.TokenConfig.Expires)
	if er7 != nil {
		return nil, false, er7
	}
	account.Token = token
	return account, false, nil
}

// VerifyAzureADJWTClaims verify if the claims information carried by jwt are valid or not.
func VerifyAzureADJWTClaims(token *jwt.Token, tenantId string, clientID string) (string, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("token claims are invalid")
	}

	tenantID, ok := claims["tid"].(string)
	if !ok || tenantID != tenantId {
		return "", errors.New("tid is invalid")
	}

	oid, ok := claims["oid"].(string) // user id on azure are called object id
	if !ok {
		return "", errors.New("oid is invalid")
	}

	if !claims.VerifyAudience(clientID, true) { // client id or app id
		return "", errors.New("aud is invalid")
	}

	return oid, nil
}

// VerifyAzureADJWT verify if an jwt is issued by Azure AD.
func VerifyAzureADJWT(ctx context.Context, tokenString string) (*jwt.Token, error) {
	keySet, err := jwk.Fetch(ctx, "https://login.microsoftonline.com/common/discovery/v2.0/keys")
	if err != nil {
		return nil, fmt.Errorf("cannot fetch public keys: %w", err)
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if token.Method.Alg() != jwa.RS256.String() {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		kid, ok := token.Header["kid"].(string)
		if !ok {
			return nil, fmt.Errorf("kid header not found")
		}

		keys, ok := keySet.LookupKeyID(kid)
		if !ok {
			return nil, fmt.Errorf("key %v not found", kid)
		}

		publicKey := &rsa.PublicKey{}
		err = keys.Raw(publicKey)
		if err != nil {
			return nil, fmt.Errorf("could not parse pubkey %w", err)
		}

		return publicKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("cannot parse token string: %w", err)
	}
	return token, nil
}
