package auth

type SchemaConfig struct {
	Id         string `mapstructure:"id"`
	UserId     string `mapstructure:"user_id"`
	UserName   string `mapstructure:"user_name"`
	Password   string `mapstructure:"password"`
	TwoFactors string `mapstructure:"two_factors"`

	SuccessTime         string `mapstructure:"success_time"`
	FailTime            string `mapstructure:"fail_time"`
	FailCount           string `mapstructure:"fail_count"`
	LockedUntilTime     string `mapstructure:"locked_until_time"`
	PasswordChangedTime string `mapstructure:"password_changed_time"`
	Status              string `mapstructure:"status"`

	Contact        string `mapstructure:"contact"`
	DisplayName    string `mapstructure:"display_name"`
	MaxPasswordAge string `mapstructure:"max_password_age"`
	UserType       string `mapstructure:"user_type"`
	Roles          string `mapstructure:"roles"`
	AccessDateFrom string `mapstructure:"access_date_from"`
	AccessDateTo   string `mapstructure:"access_date_to"`
	AccessTimeFrom string `mapstructure:"access_time_from"`
	AccessTimeTo   string `mapstructure:"access_time_to"`

	Language   string `mapstructure:"language"`
	Gender     string `mapstructure:"gender"`
	DateFormat string `mapstructure:"date_format"`
	TimeFormat string `mapstructure:"time_format"`
	ImageUrl   string `mapstructure:"image_url"`
}
