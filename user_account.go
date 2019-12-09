package auth

import "time"

type UserAccount struct {
	UserId              string       `json:"userId,omitempty" bson:"userId,omitempty" gorm:"column:userid"`
	UserName            string       `json:"userName,omitempty" bson:"userName,omitempty" gorm:"column:username"`
	Email               string       `json:"email,omitempty" bson:"email,omitempty" gorm:"column:email"`
	DisplayName         string       `json:"displayName,omitempty" bson:"displayName,omitempty" gorm:"column:displayname"`
	PasswordExpiredTime *time.Time   `json:"passwordExpiredTime,omitempty" bson:"passwordExpiredTime,omitempty" gorm:"column:passwordexpiredtime"`
	Token               string       `json:"token,omitempty" bson:"token,omitempty" gorm:"column:token"`
	TokenExpiredTime    *time.Time   `json:"tokenExpiredTime,omitempty" bson:"tokenExpiredTime,omitempty" gorm:"column:tokenexpiredtime"`
	NewUser             bool         `json:"newUser,omitempty" bson:"newUser,omitempty" gorm:"column:newuser"`
	UserType            string       `json:"userType,omitempty" bson:"userType,omitempty" gorm:"column:usertype"`
	Roles               *[]string    `json:"roles,omitempty" bson:"roles,omitempty" gorm:"column:roles"`
	Privileges          *[]Privilege `json:"privileges,omitempty" bson:"privileges,omitempty" gorm:"column:privileges"`
}
