package auth

type StatusConfig struct {
	NotFound              *int `mapstructure:"not_found"`
	Fail                  *int `mapstructure:"fail"`
	Success               *int `mapstructure:"success"`
	SuccessAndReactivated *int `mapstructure:"success_and_reactivated"`
	TwoFactorRequired     *int `mapstructure:"two_factor_required"`
	WrongPassword         *int `mapstructure:"wrong_password"`
	PasswordExpired       *int `mapstructure:"password_expired"`
	AccessTimeLocked      *int `mapstructure:"access_time_locked"`
	Locked                *int `mapstructure:"locked"`
	Suspended             *int `mapstructure:"suspended"`
	Disabled              *int `mapstructure:"disabled"`
	Error                 *int `mapstructure:"Error"`
}
type Status struct {
	NotFound              int `mapstructure:"not_found"`
	Fail                  int `mapstructure:"fail"`
	Success               int `mapstructure:"success"`
	SuccessAndReactivated int `mapstructure:"success_and_reactivated"`
	TwoFactorRequired     int `mapstructure:"two_factor_required"`
	WrongPassword         int `mapstructure:"wrong_password"`
	PasswordExpired       int `mapstructure:"password_expired"`
	AccessTimeLocked      int `mapstructure:"access_time_locked"`
	Locked                int `mapstructure:"locked"`
	Suspended             int `mapstructure:"suspended"`
	Disabled              int `mapstructure:"disabled"`
	Error                 int `mapstructure:"error"`
}

func InitStatus(c0 *StatusConfig) Status {
	var c StatusConfig
	if c0 != nil {
		c = *c0
	}
	var s Status
	if c.Error != nil {
		s.Error = *c.Error
	} else {
		s.Error = 4
	}
	if c.Fail != nil {
		s.Fail = *c.Fail
	} else {
		s.Fail = 0
	}
	if c.Success != nil {
		s.Success = *c.Success
	} else {
		s.Success = 1
	}
	if c.SuccessAndReactivated != nil {
		s.SuccessAndReactivated = *c.SuccessAndReactivated
	} else {
		s.SuccessAndReactivated = s.Success
	}
	if c.TwoFactorRequired != nil {
		s.TwoFactorRequired = *c.TwoFactorRequired
	} else {
		s.TwoFactorRequired = 2
	}
	if c.WrongPassword != nil {
		s.WrongPassword = *c.WrongPassword
	} else {
		s.WrongPassword = s.Fail
	}
	if c.PasswordExpired != nil {
		s.PasswordExpired = *c.PasswordExpired
	} else {
		s.PasswordExpired = s.Fail
	}
	if c.AccessTimeLocked != nil {
		s.AccessTimeLocked = *c.AccessTimeLocked
	} else {
		s.AccessTimeLocked = s.Fail
	}
	if c.Locked != nil {
		s.Locked = *c.Locked
	} else {
		s.Locked = s.Fail
	}
	if c.Suspended != nil {
		s.Suspended = *c.Suspended
	} else {
		s.Suspended = s.Fail
	}
	if c.Disabled != nil {
		s.Disabled = *c.Disabled
	} else {
		s.Disabled = s.Fail
	}
	if c.Error == nil && s.Error != 4 {
		s.Error = s.Fail
	}
	if c.NotFound != nil {
		s.NotFound = *c.NotFound
	} else {
		s.NotFound = s.Fail
	}
	return s
}
