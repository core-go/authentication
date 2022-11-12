package azure

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/confidential"
)

type AzureUser struct {
	Id                string `json:"id"`
	DisplayName       string `json:"displayName"`
	MobilePhone       string `json:"mobilePhone"`
	UserPrincipalName string `json:"userPrincipalName"`

	GivenName         string `json:"givenName"`
	JobTitle          string `json:"jobTitle"`
	Mail              string `json:"mail"`
	OfficeLocation    string `json:"officeLocation"`
	PreferredLanguage string `json:"preferredLanguage"`
	Surname           string `json:"surname"`
}

type AzureClient struct {
	httpClient *http.Client
	credential *confidential.Credential
	client     *confidential.Client
	config     Config
}

func NewAzureClient(httpClient *http.Client, cfg Config) (*AzureClient, error) {
	// confidential clients have a credential, such as a secret or a certificate
	cred, err := confidential.NewCredFromSecret(cfg.ClientSecret)
	if err != nil {
		return nil, fmt.Errorf("could not create a confidential from the secret: %w", err)
	}
	client, err := confidential.New(fmt.Sprintf("https://login.microsoftonline.com/%s", cfg.TenantId), cfg.ClientId, cred)
	return &AzureClient{httpClient, &cred, &client, cfg}, err
}

func (a AzureClient) GetUserByToken(ctx context.Context, azureToken string) (*AzureUser, error) {
	r, err := a.client.AcquireTokenOnBehalfOf(ctx, azureToken, a.config.Scopes)
	if err != nil {
		return nil, err
	}
	response, err := MakeRequest(ctx, a.httpClient, http.MethodGet, "https://graph.microsoft.com/v1.0/me", nil, map[string]string{
		"Authorization": "Bearer " + r.AccessToken,
	})

	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	// Unmarshal response body into User struct
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	var user AzureUser
	if err := json.Unmarshal(body, &user); err != nil {
		return &user, err
	}
	return &user, nil
}

func MakeRequest(ctx context.Context, client *http.Client, method string, url string, body []byte, headers map[string]string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	req.Header.Set("Content-Type", "application/json")
	return client.Do(req)
}
