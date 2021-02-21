package auth

type StatusConfig struct {
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

func InitStatus(c StatusConfig) Status {
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
	return s
}

/*
const (
	StatusSuccess               = Status(0)
	StatusSuccessAndReactivated = Status(1)
	StatusTwoFactorRequired     = Status(2)
	StatusFail                  = Status(3)
	StatusWrongPassword         = Status(4)
	StatusPasswordExpired       = Status(5)
	StatusAccessTimeLocked      = Status(6)
	StatusLocked                = Status(7)
	StatusSuspended             = Status(8)
	StatusDisabled              = Status(9)
	StatusSystemError           = Status(10)
)
*/
