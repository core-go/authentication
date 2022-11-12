package oauth2

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	u "net/url"
	"strings"
)

type AmazonUserRepository struct {
	CallbackURL string
}

func NewAmazonUserRepository(callbackURL string) *AmazonUserRepository {
	return &AmazonUserRepository{callbackURL}
}

func (g *AmazonUserRepository) GetUserFromOAuth2(ctx context.Context, urlRedirect string, clientId string, clientSecret string, code string) (*User, string, error) {
	url := "https://api.amazon.com/auth/o2/token"
	reqBody := u.Values{}
	reqBody.Set("grant_type", "authorization_code")
	reqBody.Set("code", code)
	reqBody.Set("client_id", clientId)
	reqBody.Set("client_secret", clientSecret)
	reqBody.Set("redirect_uri", g.CallbackURL)
	accessToken, er0 := g.getAccessToken(url, reqBody)
	if er0 != nil {
		return nil, "", er0
	}

	var infoAmazon AmazonInfo

	url1 := "https://api.amazon.com/user/profile?access_token=" + accessToken
	resp, er1 := http.Get(url1)
	if resp.StatusCode != 200 || er1 != nil {
		return nil, accessToken, er1
	}
	contents, er2 := ioutil.ReadAll(resp.Body)
	if er2 != nil {
		return nil, accessToken, er2
	}
	er3 := json.Unmarshal(contents, &infoAmazon)
	if er3 != nil {
		return nil, accessToken, er3
	}
	var user User
	user.Account = infoAmazon.UserId
	user.DisplayName = infoAmazon.Name
	user.Email = infoAmazon.Email
	// user.Gender = GenderUnknown
	return &user, accessToken, nil
}

func (g *AmazonUserRepository) GetRequestTokenOAuth(ctx context.Context, key string, secret string) (string, error) {
	return key, nil
}

func (g *AmazonUserRepository) getAccessToken(url string, body u.Values) (string, error) {
	res, er0 := http.NewRequest("POST", url, strings.NewReader(body.Encode()))
	res.Header.Add("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")

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

	fmt.Println("bearer", string(bearer))
	if er2 != nil {
		return "", er2
	}

	er3 := json.Unmarshal(bearer, &tok)
	if er3 != nil {
		return "", er3
	}
	return tok.AccessToken, nil
}
