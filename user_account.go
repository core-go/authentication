package auth

import "time"

type UserAccount struct {
	UserId              string       `json:"userId,omitempty" gorm:"column:userid" bson:"_id,omitempty" dynamodbav:"userId,omitempty" firestore:"userId,omitempty"`
	Username            string       `json:"username,omitempty" gorm:"column:username" bson:"username,omitempty" dynamodbav:"username,omitempty" firestore:"username,omitempty"`
	Contact             string       `json:"contact,omitempty" gorm:"column:contact" bson:"contact,omitempty" dynamodbav:"contact,omitempty" firestore:"contact,omitempty"`
	DisplayName         string       `json:"displayName,omitempty" gorm:"column:displayname" bson:"displayName,omitempty" dynamodbav:"displayName,omitempty" firestore:"displayName,omitempty"`
	PasswordExpiredTime *time.Time   `json:"passwordExpiredTime,omitempty" gorm:"column:passwordexpiredtime" bson:"passwordExpiredTime,omitempty" dynamodbav:"passwordExpiredTime,omitempty" firestore:"passwordExpiredTime,omitempty"`
	Token               string       `json:"token,omitempty" gorm:"column:token" bson:"token,omitempty" dynamodbav:"token,omitempty" firestore:"token,omitempty"`
	TokenExpiredTime    *time.Time   `json:"tokenExpiredTime,omitempty" gorm:"column:tokenexpiredtime" bson:"tokenExpiredTime,omitempty" dynamodbav:"tokenExpiredTime,omitempty" firestore:"tokenExpiredTime,omitempty"`
	NewUser             bool         `json:"newUser,omitempty" gorm:"column:newuser" bson:"newUser,omitempty" dynamodbav:"newUser,omitempty" firestore:"newUser,omitempty"`
	UserType            string       `json:"userType,omitempty" gorm:"column:usertype" bson:"userType,omitempty" dynamodbav:"userType,omitempty" firestore:"userType,omitempty"`
	Roles               *[]string    `json:"roles,omitempty" gorm:"column:roles" bson:"roles,omitempty" dynamodbav:"roles,omitempty" firestore:"roles,omitempty"`
	Privileges          *[]Privilege `json:"privileges,omitempty" gorm:"column:privileges" bson:"privileges,omitempty" dynamodbav:"privileges,omitempty" firestore:"privileges,omitempty"`
}
