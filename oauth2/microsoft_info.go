package oauth2

type MicrosoftInfo struct {
	Id          string
	Email       string `json:"UserPrincipalName"`
	DisplayName string
	GivenName   string
	Surname     string
}
