package twitter

type TwitterAccessToken struct {
	Token       string `json:"oauth_token"`
	TokenSecret string `json:"oauth_token_secret"`
	UserId      string `json:"user_id"`
	ScreenName  string `json:"screen_name"`
}

type TwitterInfo struct {
	Id         int    `json: "Id"`
	Name       string `json: "Name"`
	ScreenName string `json:"Screen_name"`
	Picture    string `json:"Profile_image_url"`
}
