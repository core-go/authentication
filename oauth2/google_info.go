package oauth2

type GoogleInfo struct {
	Id            string
	Email         string
	VerifiedEmail bool `json:"Verified_email"`
	Name          string
	FirstName     string `json:"Given_name"`
	LastName      string `json:"Family_name"`
	Picture       string
	Locale        string
}
