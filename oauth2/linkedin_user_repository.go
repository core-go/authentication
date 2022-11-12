package oauth2

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type LinkedInUserRepository struct {
}

func NewLinkedInUserRepository() *LinkedInUserRepository {
	return &LinkedInUserRepository{}
}

func (l *LinkedInUserRepository) GetUserFromOAuth2(ctx context.Context, urlRedirect string, clientId string, clientSecret string, code string) (*User, string, error) {
	url := "https://www.linkedin.com/oauth/v2/accessToken?grant_type=authorization_code&code=" + code + "&redirect_uri=" + urlRedirect + "&client_id=" + clientId + "&client_secret=" + clientSecret
	accessToken, er0 := l.getAccessToken(url)
	if er0 != nil {
		return nil, "", er0
	}
	var infoLinkedIn linkedInInfo
	var handel1 linkedInElements

	url1 := "https://api.linkedin.com/v2/me?projection=(id,localizedFirstName,localizedLastName)&oauth2_access_token=" + accessToken
	urlEmail := "https://api.linkedin.com/v2/emailAddress?q=members&projection=(elements*(handle~))&oauth2_access_token=" + accessToken
	resp, er1 := http.Get(url1)
	if er1 != nil {
		return nil, accessToken, er1
	}
	resp1, er2 := http.Get(urlEmail)
	if er2 != nil {
		return nil, accessToken, er2
	}
	contents, er3 := ioutil.ReadAll(resp.Body)
	if er3 != nil {
		return nil, accessToken, er3
	}
	contentsE, er4 := ioutil.ReadAll(resp1.Body)
	if er4 != nil {
		return nil, accessToken, er4
	}
	er5 := json.Unmarshal(contents, &infoLinkedIn)
	if er5 != nil {
		return nil, accessToken, er5
	}
	er6 := json.Unmarshal(contentsE, &handel1)
	if er6 != nil {
		return nil, accessToken, er6
	}
	infoLinkedIn.Elements = handel1.Elements
	if len(infoLinkedIn.Id) == 0 {
		return nil, accessToken, fmt.Errorf("LinkedIn Id cannot be empty")
	}
	var user User
	user.Account = infoLinkedIn.Id
	user.GivenName = infoLinkedIn.FirstName
	user.FamilyName = infoLinkedIn.LastName
	user.DisplayName = infoLinkedIn.LastName + " " + infoLinkedIn.FirstName
	user.Email = infoLinkedIn.Elements[0].Email.EmailAddress
	// user.Gender = GenderUnknown
	return &user, accessToken, nil
}

func (l *LinkedInUserRepository) GetRequestTokenOAuth(ctx context.Context, key string, secret string) (string, error) {
	return key, nil
}

func (l *LinkedInUserRepository) getAccessToken(url string) (string, error) {
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
