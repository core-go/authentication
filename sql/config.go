package sql

import auth "github.com/core-go/authentication"

type DBConfig struct {
	Id              string `yaml:"id" mapstructure:"id" json:"id,omitempty" gorm:"column:id" bson:"id,omitempty" dynamodbav:"id,omitempty" firestore:"id,omitempty"`
	User            string `yaml:"user" mapstructure:"user" json:"user,omitempty" gorm:"column:user" bson:"user,omitempty" dynamodbav:"user,omitempty" firestore:"user,omitempty"`
	Password        string `yaml:"password" mapstructure:"password" json:"password,omitempty" gorm:"column:password" bson:"password,omitempty" dynamodbav:"password,omitempty" firestore:"password,omitempty"`
	Username        string `yaml:"username" mapstructure:"username" json:"username,omitempty" gorm:"column:username" bson:"username,omitempty" dynamodbav:"username,omitempty" firestore:"username,omitempty"`
	SuccessTime     string `yaml:"success_time" mapstructure:"success_time" json:"successTime,omitempty" gorm:"column:successTime" bson:"successTime,omitempty" dynamodbav:"successTime,omitempty" firestore:"successTime,omitempty"`
	FailTime        string `yaml:"fail_time" mapstructure:"fail_time" json:"failTime,omitempty" gorm:"column:failtime" bson:"failTime,omitempty" dynamodbav:"failTime,omitempty" firestore:"failTime,omitempty"`
	FailCount       string `yaml:"fail_count" mapstructure:"fail_count" json:"failCount,omitempty" gorm:"column:failcount" bson:"failCount,omitempty" dynamodbav:"failCount,omitempty" firestore:"failCount,omitempty"`
	LockedUntilTime string `yaml:"locked_until_time" mapstructure:"locked_until_time" json:"lockedUntilTime,omitempty" gorm:"column:lockeduntiltime" bson:"lockedUntilTime,omitempty" dynamodbav:"lockedUntilTime,omitempty" firestore:"lockedUntilTime,omitempty"`
	Status          string `yaml:"status" mapstructure:"status" json:"status,omitempty" gorm:"column:status" bson:"status,omitempty" dynamodbav:"status,omitempty" firestore:"status,omitempty"`
	MaxPasswordAge  string `yaml:"max_password_age" mapstructure:"max_password_age" json:"maxPasswordAge,omitempty" gorm:"column:maxpasswordage" bson:"maxPasswordAge,omitempty" dynamodbav:"maxPasswordAge,omitempty" firestore:"maxPasswordAge,omitempty"`
}
type TemplateConfig struct {
	Subject string `yaml:"subject" mapstructure:"subject" json:"subject,omitempty" gorm:"column:subject" bson:"subject,omitempty" dynamodbav:"subject,omitempty" firestore:"subject,omitempty"`
	Body    string `yaml:"body" mapstructure:"body" json:"body,omitempty" gorm:"column:body" bson:"body,omitempty" dynamodbav:"body,omitempty" firestore:"body,omitempty"`
}
type SqlAuthConfig struct {
	Query             string                `yaml:"query" mapstructure:"query" json:"query,omitempty" gorm:"column:query" bson:"query,omitempty" dynamodbav:"query,omitempty" firestore:"query,omitempty"`
	Token             auth.TokenConfig      `yaml:"token" mapstructure:"token" json:"token,omitempty" gorm:"column:token" bson:"token,omitempty" dynamodbav:"token,omitempty" firestore:"token,omitempty"`
	Status            *auth.StatusConfig    `yaml:"status" mapstructure:"payload" json:"payload,omitempty" gorm:"column:payload" bson:"payload,omitempty" dynamodbav:"payload,omitempty" firestore:"payload,omitempty"`
	LockedMinutes     int                   `yaml:"lockedMinutes" mapstructure:"lockedMinutes" json:"lockedMinutes,omitempty" gorm:"column:lockedMinutes" bson:"lockedMinutes,omitempty" dynamodbav:"lockedMinutes,omitempty" firestore:"lockedMinutes,omitempty"`
	MaxPasswordFailed int                   `yaml:"maxPasswordFailed" mapstructure:"maxPasswordFailed" json:"maxPasswordFailed,omitempty" gorm:"column:maxPasswordFailed" bson:"maxPasswordFailed,omitempty" dynamodbav:"maxPasswordFailed,omitempty" firestore:"maxPasswordFailed,omitempty"`
	Payload           auth.PayloadConfig    `yaml:"payload" mapstructure:"payload" json:"payload,omitempty" gorm:"column:payload" bson:"payload,omitempty" dynamodbav:"payload,omitempty" firestore:"payload,omitempty"`
	UserStatus        auth.UserStatusConfig `yaml:"user_status" mapstructure:"user_status" json:"userStatus,omitempty" gorm:"column:userstatus" bson:"userStatus,omitempty" dynamodbav:"userStatus,omitempty" firestore:"userStatus,omitempty"`
	DB                DBConfig              `yaml:"db" mapstructure:"db" json:"db,omitempty" gorm:"column:db" bson:"db,omitempty" dynamodbav:"db,omitempty" firestore:"db,omitempty"`
	Expires           int64                 `yaml:"expires" mapstructure:"expires" json:"expires,omitempty" gorm:"column:expires" bson:"expires,omitempty" dynamodbav:"expires,omitempty" firestore:"expires,omitempty"`
	Template          *TemplateConfig       `yaml:"template" mapstructure:"template" json:"template,omitempty" gorm:"column:template" bson:"template,omitempty" dynamodbav:"template,omitempty" firestore:"template,omitempty"`
}
