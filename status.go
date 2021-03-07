package auth

type StatusConfig struct {
	Timeout               *int `mapstructure:"timeout"`
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
	Timeout               int `mapstructure:"timeout"`
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

func InitStatus(c *StatusConfig) Status {
	var x StatusConfig
	if c != nil {
		x = *c
	}
	var s Status
	if x.Timeout != nil {
		s.Timeout = *x.Timeout
	} else {
		s.Timeout = -1
	}
	if x.Fail != nil {
		s.Fail = *x.Fail
	} else {
		s.Fail = 0
	}
	if x.Success != nil {
		s.Success = *x.Success
	} else {
		s.Success = 1
	}
	if x.SuccessAndReactivated != nil {
		s.SuccessAndReactivated = *x.SuccessAndReactivated
	} else {
		s.SuccessAndReactivated = s.Success
	}
	if x.TwoFactorRequired != nil {
		s.TwoFactorRequired = *x.TwoFactorRequired
	} else {
		s.TwoFactorRequired = 2
	}
	if x.WrongPassword != nil {
		s.WrongPassword = *x.WrongPassword
	} else {
		s.WrongPassword = s.Fail
	}
	if x.PasswordExpired != nil {
		s.PasswordExpired = *x.PasswordExpired
	} else {
		s.PasswordExpired = s.Fail
	}
	if x.AccessTimeLocked != nil {
		s.AccessTimeLocked = *x.AccessTimeLocked
	} else {
		s.AccessTimeLocked = s.Fail
	}
	if x.Locked != nil {
		s.Locked = *x.Locked
	} else {
		s.Locked = s.Fail
	}
	if x.Suspended != nil {
		s.Suspended = *x.Suspended
	} else {
		s.Suspended = s.Fail
	}
	if x.Disabled != nil {
		s.Disabled = *x.Disabled
	} else {
		s.Disabled = s.Fail
	}
	if x.Error != nil {
		s.Error = *x.Error
	} else {
		s.Error = s.Fail
	}
	if x.NotFound != nil {
		s.NotFound = *x.NotFound
	} else {
		s.NotFound = s.Fail
	}
	return s
}
