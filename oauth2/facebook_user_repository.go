package oauth2

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type FacebookUserRepository struct {
}

type FACEBOOK string

const (
	FacebookApiVersion FACEBOOK = "v2.5/"
	FacebookApiUrl     FACEBOOK = "https://graph.facebook.com/" + FacebookApiVersion
)

func NewFacebookUserRepository() *FacebookUserRepository {
	return &FacebookUserRepository{}
}

func (f *FacebookUserRepository) GetUserFromOAuth2(ctx context.Context, urlRedirect string, clientId string, clientSecret string, code string) (*User, string, error) {
	url := string(FacebookApiUrl) + "oauth/access_token?client_id=" + clientId + "&redirect_uri=" + urlRedirect + "&client_secret=" + clientSecret + "&code=" + code
	accessToken, er0 := f.getAccessToken(url)
	if er0 != nil {
		return nil, "", er0
	}

	var infoFacebook facebookInfo
	url1 := string(FacebookApiUrl) + "me?fields=id,name,email,first_name,gender,last_name,picture,timezone&access_token=" + accessToken
	resp, er1 := http.Get(url1)
	if resp.StatusCode != 200 || er1 != nil {
		return nil, accessToken, er1
	}
	contents, er2 := ioutil.ReadAll(resp.Body)
	if er2 != nil {
		return nil, accessToken, er2
	}
	er3 := json.Unmarshal(contents, &infoFacebook)
	if er3 != nil {
		return nil, accessToken, er3
	}
	var user User
	user.Account = infoFacebook.Id
	user.GivenName = infoFacebook.FirstName
	user.FamilyName = infoFacebook.LastName
	user.DisplayName = infoFacebook.Name
	user.Email = infoFacebook.Email

	if infoFacebook.Gender == "male" {
		g := "M"
		user.Gender = &g
	} else if infoFacebook.Gender == "female" {
		g := "F"
		user.Gender = &g
	}
	user.Picture = infoFacebook.Picture.Data.Url
	return &user, accessToken, nil
}

func (f *FacebookUserRepository) GetRequestTokenOAuth(ctx context.Context, key string, secret string) (string, error) {
	return key, nil
}

func (f *FacebookUserRepository) getAccessToken(url string) (string, error) {
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

	if er3 != nil {
		return "", er3
	}
	return tok.AccessToken, nil
}
