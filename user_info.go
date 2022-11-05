package auth

import "time"

type UserInfo struct {
	Id                  string     `yaml:"id" mapstructure:"id" json:"id,omitempty" gorm:"column:id" bson:"_id,omitempty" dynamodbav:"id,omitempty" firestore:"id,omitempty"`
	Username            string     `yaml:"username" mapstructure:"username" json:"username,omitempty" gorm:"column:username" bson:"username,omitempty" dynamodbav:"username,omitempty" firestore:"username,omitempty"`
	Contact             *string    `yaml:"contact" mapstructure:"contact" json:"contact,omitempty" gorm:"column:contact" bson:"contact,omitempty" dynamodbav:"contact,omitempty" firestore:"contact,omitempty"`
	Email               *string    `yaml:"email" mapstructure:"email" json:"email,omitempty" gorm:"column:email" bson:"email,omitempty" dynamodbav:"email,omitempty" firestore:"email,omitempty"`
	Phone               *string    `yaml:"phone" mapstructure:"phone" json:"phone,omitempty" gorm:"column:phone" bson:"phone,omitempty" dynamodbav:"phone,omitempty" firestore:"phone,omitempty"`
	DisplayName         *string    `yaml:"display_name" mapstructure:"display_name" json:"displayName,omitempty" gorm:"column:displayname" bson:"displayName,omitempty" dynamodbav:"displayName,omitempty" firestore:"displayName,omitempty"`
	Password            string     `yaml:"password" mapstructure:"password" json:"password,omitempty" gorm:"column:password" bson:"password,omitempty" dynamodbav:"password,omitempty" firestore:"password,omitempty"`
	Disable             bool       `yaml:"disable" mapstructure:"disable" json:"disable,omitempty" gorm:"column:disable" bson:"disable,omitempty" dynamodbav:"disable,omitempty" firestore:"disable,omitempty" true:"D"`
	Deactivated         *bool      `yaml:"deactivated" mapstructure:"deactivated" json:"deactivated,omitempty" gorm:"column:deactivated" bson:"deactivated,omitempty" dynamodbav:"deactivated,omitempty" firestore:"deactivated,omitempty" true:"I"`
	Suspended           bool       `yaml:"suspended" mapstructure:"suspended" json:"suspended,omitempty" gorm:"column:suspended" bson:"suspended,omitempty" dynamodbav:"suspended,omitempty" firestore:"suspended,omitempty" true:"S"`
	LockedUntilTime     *time.Time `yaml:"locked_until_time" mapstructure:"locked_until_time" json:"lockedUntilTime,omitempty" gorm:"column:lockeduntiltime" bson:"lockedUntilTime,omitempty" dynamodbav:"lockedUntilTime,omitempty" firestore:"lockedUntilTime,omitempty"`
	SuccessTime         *time.Time `yaml:"success_time" mapstructure:"success_time" json:"successTime,omitempty" gorm:"column:successtime" bson:"successTime,omitempty" dynamodbav:"successTime,omitempty" firestore:"successTime,omitempty"`
	FailTime            *time.Time `yaml:"fail_time" mapstructure:"fail_time" json:"failTime,omitempty" gorm:"column:failtime" bson:"failTime,omitempty" dynamodbav:"failTime,omitempty" firestore:"failTime,omitempty"`
	FailCount           *int       `yaml:"fail_count" mapstructure:"fail_count" json:"failCount,omitempty" gorm:"column:failcount" bson:"failCount,omitempty" dynamodbav:"failCount,omitempty" firestore:"failCount,omitempty"`
	PasswordChangedTime *time.Time `yaml:"password_changed_time" mapstructure:"password_changed_time" json:"passwordChangedTime,omitempty" gorm:"column:passwordchangedtime" bson:"passwordChangedTime,omitempty" dynamodbav:"passwordChangedTime,omitempty" firestore:"passwordChangedTime,omitempty"`
	MaxPasswordAge      *int32     `yaml:"max_password_age" mapstructure:"max_password_age" json:"maxPasswordAge,omitempty" gorm:"column:maxpasswordage" bson:"maxPasswordAge,omitempty" dynamodbav:"maxPasswordAge,omitempty" firestore:"maxPasswordAge,omitempty"`
	UserType            *string    `yaml:"user_type" mapstructure:"user_type" json:"userType,omitempty" gorm:"column:usertype" bson:"userType,omitempty" dynamodbav:"userType,omitempty" firestore:"userType,omitempty"`
	Roles               []string   `yaml:"roles" mapstructure:"roles" json:"roles,omitempty" gorm:"column:roles" bson:"roles,omitempty" dynamodbav:"roles,omitempty" firestore:"roles,omitempty"`
	Privileges          []string   `yaml:"privileges" mapstructure:"privileges" json:"privileges,omitempty" gorm:"column:privileges" bson:"privileges,omitempty" dynamodbav:"privileges,omitempty" firestore:"privileges,omitempty"`
	AccessDateFrom      *time.Time `yaml:"access_date_from" mapstructure:"access_date_from" json:"accessDateFrom,omitempty" gorm:"column:accessdatefrom" bson:"accessDateFrom,omitempty" dynamodbav:"accessDateFrom,omitempty" firestore:"accessDateFrom,omitempty"`
	AccessDateTo        *time.Time `yaml:"access_date_to" mapstructure:"access_date_to" json:"accessDateTo,omitempty" gorm:"column:accessdateto" bson:"accessDateTo,omitempty" dynamodbav:"accessDateTo,omitempty" firestore:"accessDateTo,omitempty"`
	AccessTimeFrom      *time.Time `yaml:"access_time_from" mapstructure:"access_time_from" json:"accessTimeFrom,omitempty" gorm:"column:accesstimefrom" bson:"accessTimeFrom,omitempty" dynamodbav:"accessTimeFrom,omitempty" firestore:"accessTimeFrom,omitempty"`
	AccessTimeTo        *time.Time `yaml:"access_time_to" mapstructure:"access_time_to" json:"accessTimeTo,omitempty" gorm:"column:accesstimeto" bson:"accessTimeTo,omitempty" dynamodbav:"accessTimeTo,omitempty" firestore:"accessTimeTo,omitempty"`
	TwoFactors          bool       `yaml:"two_factors" mapstructure:"two_factors" json:"twoFactors,omitempty" gorm:"column:twofactors" bson:"twoFactors,omitempty" dynamodbav:"twoFactors,omitempty" firestore:"twoFactors,omitempty" true:"A"`
	Language            *string    `yaml:"language" mapstructure:"language" json:"language,omitempty" gorm:"column:language" bson:"language,omitempty" dynamodbav:"language,omitempty" firestore:"language,omitempty"`
	Gender              *string    `yaml:"gender" mapstructure:"gender" json:"gender,omitempty" gorm:"column:gender" bson:"gender,omitempty" dynamodbav:"gender,omitempty" firestore:"gender,omitempty"`
	DateFormat          *string    `yaml:"date_format" mapstructure:"date_format" json:"dateFormat,omitempty" gorm:"column:dateformat" bson:"dateFormat,omitempty" dynamodbav:"dateFormat,omitempty" firestore:"dateFormat,omitempty"`
	TimeFormat          *string    `yaml:"time_format" mapstructure:"time_format" json:"timeFormat,omitempty" gorm:"column:timeformat" bson:"timeFormat,omitempty" dynamodbav:"timeFormat,omitempty" firestore:"timeFormat,omitempty"`
	ImageURL            *string    `yaml:"image_url" mapstructure:"image_url" json:"imageURL,omitempty" gorm:"column:imageurl" bson:"imageURL,omitempty" dynamodbav:"imageURL,omitempty" firestore:"imageURL,omitempty"`
	Status              *string    `yaml:"status" mapstructure:"status" json:"status,omitempty" gorm:"column:status" bson:"status,omitempty" dynamodbav:"status,omitempty" firestore:"status,omitempty"`
}
