package oauth2

type OAuth2SchemaConfig struct {
	UserId   string `mapstructure:"user_id"`
	UserName string `mapstructure:"user_name"`
	Email    string `mapstructure:"email"`
	Status   string `mapstructure:"status"`

	OAuth2Email string `mapstructure:"oauth2_email"`
	Account     string `mapstructure:"account"`
	Active      string `mapstructure:"active"`

	DisplayName string `mapstructure:"display_name"`
	Picture     string `mapstructure:"picture"`
	Locale      string `mapstructure:"locale"`
	Gender      string `mapstructure:"gender"`

	DateOfBirth string `mapstructure:"date_of_birth"`
	GivenName   string `mapstructure:"given_name"`
	MiddleName  string `mapstructure:"middle_name"`
	FamilyName  string `mapstructure:"family_name"`

	CreatedTime string `mapstructure:"created_time"`
	CreatedBy   string `mapstructure:"created_by"`
	UpdatedTime string `mapstructure:"updated_time"`
	UpdatedBy   string `mapstructure:"updated_by"`
	Version     string `mapstructure:"version"`
}
