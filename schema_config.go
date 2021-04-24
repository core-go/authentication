package auth

type SchemaConfig struct {
	Id         string `mapstructure:"id" json:"id,omitempty" gorm:"column:id" bson:"id,omitempty" dynamodbav:"id,omitempty" firestore:"id,omitempty"`
	UserId     string `mapstructure:"user_id" json:"userId,omitempty" gorm:"column:userid" bson:"userId,omitempty" dynamodbav:"userId,omitempty" firestore:"userId,omitempty"`
	Username   string `mapstructure:"user_name" json:"username,omitempty" gorm:"column:username" bson:"username,omitempty" dynamodbav:"username,omitempty" firestore:"username,omitempty"`
	Password   string `mapstructure:"password" json:"password,omitempty" gorm:"column:password" bson:"password,omitempty" dynamodbav:"password,omitempty" firestore:"password,omitempty"`
	TwoFactors string `mapstructure:"two_factors" json:"twoFactors,omitempty" gorm:"column:twofactors" bson:"twoFactors,omitempty" dynamodbav:"twoFactors,omitempty" firestore:"twoFactors,omitempty"`

	SuccessTime         string `mapstructure:"success_time" json:"successTime,omitempty" gorm:"column:successTime" bson:"successTime,omitempty" dynamodbav:"successTime,omitempty" firestore:"successTime,omitempty"`
	FailTime            string `mapstructure:"fail_time" json:"failTime,omitempty" gorm:"column:failtime" bson:"failTime,omitempty" dynamodbav:"failTime,omitempty" firestore:"failTime,omitempty"`
	FailCount           string `mapstructure:"fail_count" json:"failCount,omitempty" gorm:"column:failcount" bson:"failCount,omitempty" dynamodbav:"failCount,omitempty" firestore:"failCount,omitempty"`
	LockedUntilTime     string `mapstructure:"locked_until_time" json:"lockedUntilTime,omitempty" gorm:"column:lockeduntiltime" bson:"lockedUntilTime,omitempty" dynamodbav:"lockedUntilTime,omitempty" firestore:"lockedUntilTime,omitempty"`
	PasswordChangedTime string `mapstructure:"password_changed_time" json:"passwordChangedTime,omitempty" gorm:"column:passwordchangedtime" bson:"passwordChangedTime,omitempty" dynamodbav:"passwordChangedTime,omitempty" firestore:"passwordChangedTime,omitempty"`
	Status              string `mapstructure:"status" json:"status,omitempty" gorm:"column:status" bson:"status,omitempty" dynamodbav:"status,omitempty" firestore:"status,omitempty"`

	Contact        string `mapstructure:"contact" json:"contact,omitempty" gorm:"column:contact" bson:"contact,omitempty" dynamodbav:"contact,omitempty" firestore:"contact,omitempty"`
	Email          string `mapstructure:"email" json:"email,omitempty" gorm:"column:email" bson:"email,omitempty" dynamodbav:"email,omitempty" firestore:"email,omitempty"`
	Phone          string `mapstructure:"phone" json:"phone,omitempty" gorm:"column:phone" bson:"phone,omitempty" dynamodbav:"phone,omitempty" firestore:"phone,omitempty"`
	DisplayName    string `mapstructure:"display_name" json:"displayName,omitempty" gorm:"column:displayname" bson:"displayName,omitempty" dynamodbav:"displayName,omitempty" firestore:"displayName,omitempty"`
	MaxPasswordAge string `mapstructure:"max_password_age" json:"maxPasswordAge,omitempty" gorm:"column:maxpasswordage" bson:"maxPasswordAge,omitempty" dynamodbav:"maxPasswordAge,omitempty" firestore:"maxPasswordAge,omitempty"`
	UserType       string `mapstructure:"user_type" json:"userType,omitempty" gorm:"column:usertype" bson:"userType,omitempty" dynamodbav:"userType,omitempty" firestore:"userType,omitempty"`
	Roles          string `mapstructure:"roles" json:"roles,omitempty" gorm:"column:roles" bson:"roles,omitempty" dynamodbav:"roles,omitempty" firestore:"roles,omitempty"`
	AccessDateFrom string `mapstructure:"access_date_from" json:"accessDateFrom,omitempty" gorm:"column:accessdatefrom" bson:"accessDateFrom,omitempty" dynamodbav:"accessDateFrom,omitempty" firestore:"accessDateFrom,omitempty"`
	AccessDateTo   string `mapstructure:"access_date_to" json:"accessDateTo,omitempty" gorm:"column:accessdateto" bson:"accessDateTo,omitempty" dynamodbav:"accessDateTo,omitempty" firestore:"accessDateTo,omitempty"`
	AccessTimeFrom string `mapstructure:"access_time_from" json:"accessTimeFrom,omitempty" gorm:"column:accesstimefrom" bson:"accessTimeFrom,omitempty" dynamodbav:"accessTimeFrom,omitempty" firestore:"accessTimeFrom,omitempty"`
	AccessTimeTo   string `mapstructure:"access_time_to" json:"accessTimeTo,omitempty" gorm:"column:accesstimeto" bson:"accessTimeTo,omitempty" dynamodbav:"accessTimeTo,omitempty" firestore:"accessTimeTo,omitempty"`

	Language   string `mapstructure:"language" json:"language,omitempty" gorm:"column:language" bson:"language,omitempty" dynamodbav:"language,omitempty" firestore:"language,omitempty"`
	Gender     string `mapstructure:"gender" json:"gender,omitempty" gorm:"column:gender" bson:"gender,omitempty" dynamodbav:"gender,omitempty" firestore:"gender,omitempty"`
	DateFormat string `mapstructure:"date_format" json:"dateFormat,omitempty" gorm:"column:dateformat" bson:"dateFormat,omitempty" dynamodbav:"dateFormat,omitempty" firestore:"dateFormat,omitempty"`
	TimeFormat string `mapstructure:"time_format" json:"timeFormat,omitempty" gorm:"column:timeformat" bson:"timeFormat,omitempty" dynamodbav:"timeFormat,omitempty" firestore:"timeFormat,omitempty"`
	ImageURL   string `mapstructure:"image_url" json:"imageURL,omitempty" gorm:"column:imageurl" bson:"imageURL,omitempty" dynamodbav:"imageURL,omitempty" firestore:"imageURL,omitempty"`
}
