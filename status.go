package auth

type StatusConfig struct {
	Timeout               *int `mapstructure:"timeout" json:"timeout,omitempty" gorm:"column:timeout" bson:"timeout,omitempty" dynamodbav:"timeout,omitempty" firestore:"timeout,omitempty"`
	NotFound              *int `mapstructure:"not_found" json:"notFound,omitempty" gorm:"column:notfound" bson:"notFound,omitempty" dynamodbav:"notFound,omitempty" firestore:"notFound,omitempty"`
	Fail                  *int `mapstructure:"fail" json:"fail,omitempty" gorm:"column:fail" bson:"fail,omitempty" dynamodbav:"fail,omitempty" firestore:"fail,omitempty"`
	Success               *int `mapstructure:"success" json:"success,omitempty" gorm:"column:success" bson:"success,omitempty" dynamodbav:"success,omitempty" firestore:"success,omitempty"`
	SuccessAndReactivated *int `mapstructure:"success_and_reactivated" json:"successAndReactivated,omitempty" gorm:"column:successandreactivated" bson:"successAndReactivated,omitempty" dynamodbav:"successAndReactivated,omitempty" firestore:"successAndReactivated,omitempty"`
	TwoFactorRequired     *int `mapstructure:"two_factor_required" json:"twoFactorRequired,omitempty" gorm:"column:twofactorrequired" bson:"twoFactorRequired,omitempty" dynamodbav:"twoFactorRequired,omitempty" firestore:"twoFactorRequired,omitempty"`
	WrongPassword         *int `mapstructure:"wrong_password" json:"wrongPassword,omitempty" gorm:"column:wrongpassword" bson:"wrongPassword,omitempty" dynamodbav:"wrongPassword,omitempty" firestore:"wrongPassword,omitempty"`
	PasswordExpired       *int `mapstructure:"password_expired" json:"passwordExpired,omitempty" gorm:"column:passwordexpired" bson:"passwordExpired,omitempty" dynamodbav:"passwordExpired,omitempty" firestore:"passwordExpired,omitempty"`
	AccessTimeLocked      *int `mapstructure:"access_time_locked" json:"accessTimeLocked,omitempty" gorm:"column:accesstimelocked" bson:"accessTimeLocked,omitempty" dynamodbav:"accessTimeLocked,omitempty" firestore:"accessTimeLocked,omitempty"`
	Locked                *int `mapstructure:"locked" json:"locked,omitempty" gorm:"column:locked" bson:"locked,omitempty" dynamodbav:"locked,omitempty" firestore:"locked,omitempty"`
	Suspended             *int `mapstructure:"suspended" json:"suspended,omitempty" gorm:"column:suspended" bson:"suspended,omitempty" dynamodbav:"suspended,omitempty" firestore:"suspended,omitempty"`
	Disabled              *int `mapstructure:"disabled" json:"disabled,omitempty" gorm:"column:disabled" bson:"disabled,omitempty" dynamodbav:"disabled,omitempty" firestore:"disabled,omitempty"`
	Error                 *int `mapstructure:"Error" json:"error,omitempty" gorm:"column:error" bson:"error,omitempty" dynamodbav:"error,omitempty" firestore:"error,omitempty"`
}
type Status struct {
	Timeout               int `mapstructure:"timeout" json:"timeout,omitempty" gorm:"column:timeout" bson:"timeout,omitempty" dynamodbav:"timeout,omitempty" firestore:"timeout,omitempty"`
	NotFound              int `mapstructure:"not_found" json:"notFound,omitempty" gorm:"column:notfound" bson:"notFound,omitempty" dynamodbav:"notFound,omitempty" firestore:"notFound,omitempty"`
	Fail                  int `mapstructure:"fail" json:"fail,omitempty" gorm:"column:fail" bson:"fail,omitempty" dynamodbav:"fail,omitempty" firestore:"fail,omitempty"`
	Success               int `mapstructure:"success" json:"success,omitempty" gorm:"column:success" bson:"success,omitempty" dynamodbav:"success,omitempty" firestore:"success,omitempty"`
	SuccessAndReactivated int `mapstructure:"success_and_reactivated" json:"successAndReactivated,omitempty" gorm:"column:successandreactivated" bson:"successAndReactivated,omitempty" dynamodbav:"successAndReactivated,omitempty" firestore:"successAndReactivated,omitempty"`
	TwoFactorRequired     int `mapstructure:"two_factor_required" json:"twoFactorRequired,omitempty" gorm:"column:twofactorrequired" bson:"twoFactorRequired,omitempty" dynamodbav:"twoFactorRequired,omitempty" firestore:"twoFactorRequired,omitempty"`
	WrongPassword         int `mapstructure:"wrong_password" json:"wrongPassword,omitempty" gorm:"column:wrongpassword" bson:"wrongPassword,omitempty" dynamodbav:"wrongPassword,omitempty" firestore:"wrongPassword,omitempty"`
	PasswordExpired       int `mapstructure:"password_expired" json:"passwordExpired,omitempty" gorm:"column:passwordexpired" bson:"passwordExpired,omitempty" dynamodbav:"passwordExpired,omitempty" firestore:"passwordExpired,omitempty"`
	AccessTimeLocked      int `mapstructure:"access_time_locked" json:"accessTimeLocked,omitempty" gorm:"column:accesstimelocked" bson:"accessTimeLocked,omitempty" dynamodbav:"accessTimeLocked,omitempty" firestore:"accessTimeLocked,omitempty"`
	Locked                int `mapstructure:"locked" json:"locked,omitempty" gorm:"column:locked" bson:"locked,omitempty" dynamodbav:"locked,omitempty" firestore:"locked,omitempty"`
	Suspended             int `mapstructure:"suspended" json:"suspended,omitempty" gorm:"column:suspended" bson:"suspended,omitempty" dynamodbav:"suspended,omitempty" firestore:"suspended,omitempty"`
	Disabled              int `mapstructure:"disabled" json:"disabled,omitempty" gorm:"column:disabled" bson:"disabled,omitempty" dynamodbav:"disabled,omitempty" firestore:"disabled,omitempty"`
	Error                 int `mapstructure:"error" json:"error,omitempty" gorm:"column:error" bson:"error,omitempty" dynamodbav:"error,omitempty" firestore:"error,omitempty"`
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
