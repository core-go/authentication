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

type DropboxUserRepository struct {
}

func NewDropboxUserRepository() *DropboxUserRepository {
	return &DropboxUserRepository{}
}

func (s *DropboxUserRepository) GetUserFromOAuth2(ctx context.Context, urlRedirect string, clientId string, clientSecret string, code string) (*User, string, error) {
	url := "https://api.dropbox.com/oauth2/token"
	reqBody := u.Values{}
	reqBody.Set("grant_type", "authorization_code")
	reqBody.Set("code", code)
	reqBody.Set("client_id", clientId)
	reqBody.Set("client_secret", clientSecret)
	reqBody.Set("redirect_uri", "http://localhost:3001/index.html?redirect=oAuth2")
	accessToken, er0 := s.getAccessToken(url, reqBody)
	if er0 != nil {
		return nil, "", er0
	}

	var infoDropbox dropboxInfo

	url1 := "https://api.dropboxapi.com/2/users/get_current_account"

	bearer := "Bearer " + accessToken

	req, er1 := http.NewRequest("POST", url1, nil)

	req.Header.Add("Authorization", bearer)

	client := &http.Client{}
	resp, _ := client.Do(req)

	if resp.StatusCode != 200 || er1 != nil {
		return nil, accessToken, er1
	}
	contents, er2 := ioutil.ReadAll(resp.Body)

	if er2 != nil {
		return nil, accessToken, er2
	}
	er3 := json.Unmarshal(contents, &infoDropbox)
	if er3 != nil {
		return nil, accessToken, er3
	}
	var user User
	user.Account = infoDropbox.AccountId
	user.DisplayName = infoDropbox.Name.DisplayName
	user.GivenName = infoDropbox.Name.GivenName
	user.FamilyName = infoDropbox.Name.SurName
	user.Email = infoDropbox.Email
	user.Picture = infoDropbox.Picture
	// user.Gender = GenderUnknown
	return &user, accessToken, nil
}

func (s *DropboxUserRepository) GetRequestTokenOAuth(ctx context.Context, key string, secret string) (string, error) {
	return key, nil
}

func (s *DropboxUserRepository) getAccessToken(url string, body u.Values) (string, error) {
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
