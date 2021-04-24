package auth

import "time"

type UserAccount struct {
	Id                  string      `mapstructure:"id" json:"id,omitempty" gorm:"column:id" bson:"_id,omitempty" dynamodbav:"id,omitempty" firestore:"id,omitempty"`
	Username            string      `mapstructure:"username" json:"username,omitempty" gorm:"column:username" bson:"username,omitempty" dynamodbav:"username,omitempty" firestore:"username,omitempty"`
	Contact             string      `mapstructure:"contact" json:"contact,omitempty" gorm:"column:contact" bson:"contact,omitempty" dynamodbav:"contact,omitempty" firestore:"contact,omitempty"`
	Email               string      `mapstructure:"email" json:"email,omitempty" gorm:"column:email" bson:"email,omitempty" dynamodbav:"email,omitempty" firestore:"email,omitempty"`
	Phone               string      `mapstructure:"phone" json:"phone,omitempty" gorm:"column:phone" bson:"phone,omitempty" dynamodbav:"phone,omitempty" firestore:"phone,omitempty"`
	DisplayName         string      `mapstructure:"display_name" json:"displayName,omitempty" gorm:"column:displayname" bson:"displayName,omitempty" dynamodbav:"displayName,omitempty" firestore:"displayName,omitempty"`
	PasswordExpiredTime *time.Time  `mapstructure:"password_expired_time" json:"passwordExpiredTime,omitempty" gorm:"column:passwordexpiredtime" bson:"passwordExpiredTime,omitempty" dynamodbav:"passwordExpiredTime,omitempty" firestore:"passwordExpiredTime,omitempty"`
	Token               string      `mapstructure:"token" json:"token,omitempty" gorm:"column:token" bson:"token,omitempty" dynamodbav:"token,omitempty" firestore:"token,omitempty"`
	TokenExpiredTime    *time.Time  `mapstructure:"token_expired_time" json:"tokenExpiredTime,omitempty" gorm:"column:tokenexpiredtime" bson:"tokenExpiredTime,omitempty" dynamodbav:"tokenExpiredTime,omitempty" firestore:"tokenExpiredTime,omitempty"`
	NewUser             *bool       `mapstructure:"new_user" json:"newUser,omitempty" gorm:"column:newuser" bson:"newUser,omitempty" dynamodbav:"newUser,omitempty" firestore:"newUser,omitempty"`
	Language            string      `mapstructure:"language" json:"language,omitempty" gorm:"column:language" bson:"language,omitempty" dynamodbav:"language,omitempty" firestore:"language,omitempty"`
	Gender              string      `mapstructure:"gender" json:"gender,omitempty" gorm:"column:gender" bson:"gender,omitempty" dynamodbav:"gender,omitempty" firestore:"gender,omitempty"`
	DateFormat          string      `mapstructure:"date_format" json:"dateFormat,omitempty" gorm:"column:dateformat" bson:"dateFormat,omitempty" dynamodbav:"dateFormat,omitempty" firestore:"dateFormat,omitempty"`
	TimeFormat          string      `mapstructure:"time_format" json:"timeFormat,omitempty" gorm:"column:timeformat" bson:"timeFormat,omitempty" dynamodbav:"timeFormat,omitempty" firestore:"timeFormat,omitempty"`
	ImageURL            string      `mapstructure:"image_url" json:"imageURL,omitempty" gorm:"column:imageurl" bson:"imageURL,omitempty" dynamodbav:"imageURL,omitempty" firestore:"imageURL,omitempty"`
	Type                string      `mapstructure:"type" json:"type,omitempty" gorm:"column:type" bson:"type,omitempty" dynamodbav:"type,omitempty" firestore:"type,omitempty"`
	Roles               []string    `mapstructure:"roles" json:"roles,omitempty" gorm:"column:roles" bson:"roles,omitempty" dynamodbav:"roles,omitempty" firestore:"roles,omitempty"`
	Privileges          []Privilege `mapstructure:"privileges" json:"privileges,omitempty" gorm:"column:privileges" bson:"privileges,omitempty" dynamodbav:"privileges,omitempty" firestore:"privileges,omitempty"`
}
