package auth

type PayloadConfig struct {
	Ip         string `mapstructure:"ip"`
	UserId     string `mapstructure:"user_id"`
	Username   string `mapstructure:"username"`
	Contact    string `mapstructure:"contact"`
	UserType   string `mapstructure:"user_type"`
	Roles      string `mapstructure:"roles"`
	Privileges string `mapstructure:"privileges"`
	Tokens     string `mapstructure:"tokens"`
}
