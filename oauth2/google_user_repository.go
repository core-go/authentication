package oauth2

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type GoogleUserRepository struct {
}

func NewGoogleUserRepository() *GoogleUserRepository {
	return &GoogleUserRepository{}
}

func (g *GoogleUserRepository) GetUserFromOAuth2(ctx context.Context, urlRedirect string, clientId string, clientSecret string, code string) (*User, string, error) {
	url := "https://www.googleapis.com/oauth2/v4/token?redirect_uri=" + urlRedirect +
		"&client_id=" + clientId + "&client_secret=" + clientSecret + "&scope=&grant_type=authorization_code&code=" + code
	accessToken, er0 := g.getAccessToken(url)
	if er0 != nil {
		return nil, "", er0
	}

	var infoGoogle GoogleInfo

	url1 := "https://www.googleapis.com/oauth2/v1/userinfo?access_token=" + accessToken
	resp, er1 := http.Get(url1)
	if resp.StatusCode != 200 || er1 != nil {
		return nil, accessToken, er1
	}
	contents, er2 := ioutil.ReadAll(resp.Body)
	if er2 != nil {
		return nil, accessToken, er2
	}
	er3 := json.Unmarshal(contents, &infoGoogle)
	if er3 != nil {
		return nil, accessToken, er3
	}
	var user User
	user.Account = infoGoogle.Id
	user.GivenName = infoGoogle.FirstName
	user.FamilyName = infoGoogle.LastName
	user.DisplayName = infoGoogle.Name
	user.Email = infoGoogle.Email
	user.Picture = infoGoogle.Picture
	// user.Gender = GenderUnknown
	return &user, accessToken, nil
}

func (g *GoogleUserRepository) GetRequestTokenOAuth(ctx context.Context, key string, secret string) (string, error) {
	return key, nil
}

func (g *GoogleUserRepository) getAccessToken(url string) (string, error) {
	res, er0 := http.NewRequest("POST", url, nil)
	if er0 != nil {
		return "", er0
	}

	var tok BearerToken
	accessTokClient := &http.Client{}

	resp, er1 := accessTokClient.Do(res)
	if er1 != nil {
		return "", er1
	}

	bearer, er2 := ioutil.ReadAll(resp.Body)
	if er2 != nil {
		return "", er2
	}

	er3 := json.Unmarshal(bearer, &tok)
	// fmt.Printf("Raw Response Body:\n%v\n", string(bearer))
	if er3 != nil {
		return "", er3
	}
	return tok.AccessToken, nil
}
