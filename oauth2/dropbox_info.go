package oauth2

type dropboxInfo struct {
	AccountId string `json:"Account_id"`
	Email     string
	Name      name
	Picture   string `json:"Profile_photo_url"`
}
type name struct {
	GivenName   string `json:"Given_name"`
	SurName     string `json:"Surname"`
	DisplayName string `json:"Display_name"`
}
