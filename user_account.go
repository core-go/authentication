package auth

import "time"

type UserAccount struct {
	Id                  string      `yaml:"id" mapstructure:"id" json:"id,omitempty" gorm:"column:id" bson:"_id,omitempty" dynamodbav:"id,omitempty" firestore:"id,omitempty"`
	Username            string      `yaml:"username" mapstructure:"username" json:"username,omitempty" gorm:"column:username" bson:"username,omitempty" dynamodbav:"username,omitempty" firestore:"username,omitempty"`
	Contact             string      `yaml:"contact" mapstructure:"contact" json:"contact,omitempty" gorm:"column:contact" bson:"contact,omitempty" dynamodbav:"contact,omitempty" firestore:"contact,omitempty"`
	Email               string      `yaml:"email" mapstructure:"email" json:"email,omitempty" gorm:"column:email" bson:"email,omitempty" dynamodbav:"email,omitempty" firestore:"email,omitempty"`
	Phone               string      `yaml:"phone" mapstructure:"phone" json:"phone,omitempty" gorm:"column:phone" bson:"phone,omitempty" dynamodbav:"phone,omitempty" firestore:"phone,omitempty"`
	DisplayName         string      `yaml:"display_name" mapstructure:"display_name" json:"displayName,omitempty" gorm:"column:displayname" bson:"displayName,omitempty" dynamodbav:"displayName,omitempty" firestore:"displayName,omitempty"`
	PasswordExpiredTime *time.Time  `yaml:"password_expired_time" mapstructure:"password_expired_time" json:"passwordExpiredTime,omitempty" gorm:"column:passwordexpiredtime" bson:"passwordExpiredTime,omitempty" dynamodbav:"passwordExpiredTime,omitempty" firestore:"passwordExpiredTime,omitempty"`
	Token               string      `yaml:"token" mapstructure:"token" json:"token,omitempty" gorm:"column:token" bson:"token,omitempty" dynamodbav:"token,omitempty" firestore:"token,omitempty"`
	TokenExpiredTime    *time.Time  `yaml:"token_expired_time" mapstructure:"token_expired_time" json:"tokenExpiredTime,omitempty" gorm:"column:tokenexpiredtime" bson:"tokenExpiredTime,omitempty" dynamodbav:"tokenExpiredTime,omitempty" firestore:"tokenExpiredTime,omitempty"`
	NewUser             *bool       `yaml:"new_user" mapstructure:"new_user" json:"newUser,omitempty" gorm:"column:newuser" bson:"newUser,omitempty" dynamodbav:"newUser,omitempty" firestore:"newUser,omitempty"`
	Language            string      `yaml:"language" mapstructure:"language" json:"language,omitempty" gorm:"column:language" bson:"language,omitempty" dynamodbav:"language,omitempty" firestore:"language,omitempty"`
	Gender              string      `yaml:"gender" mapstructure:"gender" json:"gender,omitempty" gorm:"column:gender" bson:"gender,omitempty" dynamodbav:"gender,omitempty" firestore:"gender,omitempty"`
	DateFormat          string      `yaml:"date_format" mapstructure:"date_format" json:"dateFormat,omitempty" gorm:"column:dateformat" bson:"dateFormat,omitempty" dynamodbav:"dateFormat,omitempty" firestore:"dateFormat,omitempty"`
	TimeFormat          string      `yaml:"time_format" mapstructure:"time_format" json:"timeFormat,omitempty" gorm:"column:timeformat" bson:"timeFormat,omitempty" dynamodbav:"timeFormat,omitempty" firestore:"timeFormat,omitempty"`
	ImageURL            string      `yaml:"image_url" mapstructure:"image_url" json:"imageURL,omitempty" gorm:"column:imageurl" bson:"imageURL,omitempty" dynamodbav:"imageURL,omitempty" firestore:"imageURL,omitempty"`
	Type                string      `yaml:"type" mapstructure:"type" json:"type,omitempty" gorm:"column:type" bson:"type,omitempty" dynamodbav:"type,omitempty" firestore:"type,omitempty"`
	Roles               []string    `yaml:"roles" mapstructure:"roles" json:"roles,omitempty" gorm:"column:roles" bson:"roles,omitempty" dynamodbav:"roles,omitempty" firestore:"roles,omitempty"`
	Privileges          []Privilege `yaml:"privileges" mapstructure:"privileges" json:"privileges,omitempty" gorm:"column:privileges" bson:"privileges,omitempty" dynamodbav:"privileges,omitempty" firestore:"privileges,omitempty"`
}
