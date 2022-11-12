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

type MicrosoftUserRepository struct {
	CallbackURL string
}

func NewMicrosoftUserRepository(callbackURL string) *MicrosoftUserRepository {
	return &MicrosoftUserRepository{callbackURL}
}

func (g *MicrosoftUserRepository) GetUserFromOAuth2(ctx context.Context, urlRedirect string, clientId string, clientSecret string, code string) (*User, string, error) {
	url := "https://login.microsoftonline.com/common/oauth2/v2.0/token"
	reqBody := u.Values{}

	reqBody.Set("grant_type", "authorization_code")
	reqBody.Set("scope", "user.read,mail.read")
	reqBody.Set("code", code)
	reqBody.Set("client_id", clientId)
	reqBody.Set("client_secret", clientSecret)
	reqBody.Set("redirect_uri", g.CallbackURL)
	accessToken, er0 := g.getAccessToken(url, reqBody)
	if er0 != nil {
		return nil, "", er0
	}

	var infoMicrosoft MicrosoftInfo

	url1 := "https://graph.microsoft.com/v1.0/me"

	bearer := "Bearer " + accessToken

	req, _ := http.NewRequest("GET", url1, nil)

	req.Header.Add("Authorization", bearer)

	client := &http.Client{}
	resp, er1 := client.Do(req)

	if resp.StatusCode != 200 || er1 != nil {
		return nil, accessToken, er1
	}
	contents, er2 := ioutil.ReadAll(resp.Body)
	if er2 != nil {
		return nil, accessToken, er2
	}
	er3 := json.Unmarshal(contents, &infoMicrosoft)
	if er3 != nil {
		return nil, accessToken, er3
	}
	var user User
	user.Account = infoMicrosoft.Id
	user.DisplayName = infoMicrosoft.DisplayName
	user.GivenName = infoMicrosoft.GivenName
	user.FamilyName = infoMicrosoft.Surname
	user.Email = infoMicrosoft.Email
	// user.Gender = GenderUnknown
	return &user, accessToken, nil
}

func (g *MicrosoftUserRepository) GetRequestTokenOAuth(ctx context.Context, key string, secret string) (string, error) {
	return key, nil
}

func (g *MicrosoftUserRepository) getAccessToken(url string, body u.Values) (string, error) {
	res, er0 := http.NewRequest("POST", url, strings.NewReader(body.Encode()))

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
