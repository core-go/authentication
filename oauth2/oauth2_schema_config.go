package oauth2

type OAuth2SchemaConfig struct {
	Id       string `yaml:"id" mapstructure:"id"`
	Username string `yaml:"username" mapstructure:"username"`
	Email    string `yaml:"email" mapstructure:"email"`
	Status   string `yaml:"status" mapstructure:"status"`

	OAuth2Email string `mapstructure:"oauth2_email"`
	Account     string `mapstructure:"account"`
	Active      string `mapstructure:"active"`

	DisplayName string `yaml:"display_name" mapstructure:"display_name"`
	Picture     string `mapstructure:"picture" mapstructure:"picture"`
	Locale      string `mapstructure:"locale" mapstructure:"locale"`
	Gender      string `mapstructure:"gender" mapstructure:"gender"`

	DateOfBirth string `yaml:"date_of_birth" mapstructure:"date_of_birth"`
	GivenName   string `yaml:"given_name" mapstructure:"given_name"`
	MiddleName  string `yaml:"middle_name" mapstructure:"middle_name"`
	FamilyName  string `yaml:"family_name" mapstructure:"family_name"`

	CreatedTime string `yaml:"created_time" mapstructure:"created_time"`
	CreatedBy   string `yaml:"created_by" mapstructure:"created_by"`
	UpdatedTime string `yaml:"updated_time" mapstructure:"updated_time"`
	UpdatedBy   string `yaml:"updated_by" mapstructure:"updated_by"`
	Version     string `yaml:"version" mapstructure:"version"`
}
