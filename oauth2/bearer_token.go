package oauth2

type BearerToken struct {
	TokenType   string `json:"Token_Type"`
	AccessToken string `json:"Access_Token"`
}
