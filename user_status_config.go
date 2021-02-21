package auth

type UserStatusConfig struct {
	Deactivated string `mapstructure:"deactivated"`
	Disable     string `mapstructure:"disable"`
	Suspended   string `mapstructure:"suspended"`
}
