package auth

type TokenConfig struct {
	Secret  string `mapstructure:"secret"`
	Expires int64  `mapstructure:"expires"`
}
