package twitter

import (
	"context"
	"encoding/json"
	"github.com/dghubble/oauth1"
	"github.com/dghubble/oauth1/twitter"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/core-go/authentication/oauth2"
)

type TwitterUserRepository struct {
	CallbackURL string
}

func NewTwitterUserRepository(callbackURL string) *TwitterUserRepository {
	return &TwitterUserRepository{callbackURL}
}

func (g *TwitterUserRepository) GetUserFromOAuth2(ctx context.Context, urlRedirect string, clientId string, clientSecret string, code string) (*oauth2.User, string, error) {
	oauthToken := code[0:strings.Index(code, ":")]
	oauthVerifier := code[strings.Index(code, ":")+1:]

	accessToken, er0 := g.getAccessToken(oauthToken, oauthVerifier)
	if er0 != nil {
		return nil, "", er0
	}

	config := oauth1.Config{
		ConsumerKey:    clientId,
		ConsumerSecret: clientSecret,
		CallbackURL:    g.CallbackURL,
		Endpoint:       twitter.AuthorizeEndpoint,
	}
	token := oauth1.NewToken(accessToken.Token, accessToken.TokenSecret)
	httpClient := config.Client(oauth1.NoContext, token)

	//path:= "https://api.twitter.com/1.1/account/verify_credentials.json?include_email=true"
	path := "https://api.twitter.com/1.1/users/show.json?user_id=" + accessToken.UserId
	resp, er1 := httpClient.Get(path)
	if er1 != nil {
		return nil, accessToken.Token, er1
	}

	body, er2 := ioutil.ReadAll(resp.Body)
	if er2 != nil {
		return nil, accessToken.Token, er2
	}

	// fmt.Println(string(body))
	t := TwitterInfo{}
	er3 := json.Unmarshal(body, &t)
	if er3 != nil {
		return nil, accessToken.Token, er3
	}

	var user oauth2.User
	user.Account = strconv.Itoa(t.Id)
	user.GivenName = t.Name
	user.FamilyName = t.Name
	user.DisplayName = t.ScreenName
	user.Picture = t.Picture
	user.Email = t.ScreenName
	// user.Gender = oauth2.GenderUnknown
	return &user, accessToken.Token, nil
}

func (g *TwitterUserRepository) getAccessToken(oauthToken string, oauthVerifier string) (TwitterAccessToken, error) {
	url := `https://api.twitter.com/oauth/access_token?oauth_token=` + oauthToken + `&oauth_verifier=` + oauthVerifier
	t := TwitterAccessToken{}
	res, er0 := http.NewRequest("POST", url, nil)
	if er0 != nil {
		return t, er0
	}

	accessTokClient := &http.Client{}
	resp, er1 := accessTokClient.Do(res)
	if er1 != nil {
		return t, er1
	}

	bearer, er2 := ioutil.ReadAll(resp.Body)
	if er2 != nil {
		return t, er2
	}

	params := getParams(string(bearer))
	t.UserId = params["user_id"]
	t.Token = params["oauth_token"]
	t.ScreenName = params["screen_name"]
	t.TokenSecret = params["oauth_token_secret"]
	return t, nil
}

func (g *TwitterUserRepository) GetRequestTokenOAuth(ctx context.Context, key string, secret string) (string, error) {
	config := oauth1.Config{
		ConsumerKey:    key,
		ConsumerSecret: secret,
		CallbackURL:    g.CallbackURL,
		Endpoint:       twitter.AuthorizeEndpoint,
	}
	requestToken, _, err := config.RequestToken()
	if err != nil {
		return key, err
	}
	return requestToken, nil
}

func getParams(body string) map[string]string {
	params := make(map[string]string)
	arr := strings.Split(body, "&")
	for _, v := range arr {
		arr1 := strings.Split(v, "=")
		params[arr1[0]] = arr1[1]
	}
	return params
}
