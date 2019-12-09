package auth

import "time"

type UserInfo struct {
	UserId               string     `json:"userId,omitempty" bson:"userId,omitempty" gorm:"column:userid"`
	UserName             string     `json:"userName,omitempty" bson:"userName,omitempty" gorm:"column:username"`
	Email                string     `json:"email,omitempty" bson:"email,omitempty" gorm:"column:email"`
	DisplayName          string     `json:"displayName,omitempty" bson:"displayName,omitempty" gorm:"column:displayname"`
	Password             string     `json:"password,omitempty" bson:"password,omitempty" gorm:"column:password"`
	Disable              bool       `json:"disable,omitempty" bson:"disable,omitempty" gorm:"column:disable"`
	Deactivated          bool       `json:"deactivated,omitempty" bson:"deactivated,omitempty" gorm:"column:deactivated"`
	Suspended            bool       `json:"blocked,omitempty" bson:"blocked,omitempty" gorm:"column:suspended"`
	LockedUntilTime      *time.Time `json:"lockedUntilTime,omitempty" bson:"lockedUntilTime,omitempty" gorm:"column:lockeduntiltime"`
	SuccessTime          *time.Time `json:"successTime,omitempty" bson:"successTime,omitempty" gorm:"column:successtime"`
	FailTime             *time.Time `json:"failTime,omitempty" bson:"failTime,omitempty" gorm:"column:failtime"`
	FailCount            int        `json:"failCount,omitempty" bson:"failCount,omitempty" gorm:"column:failcount"`
	PasswordModifiedTime *time.Time `json:"passwordModifiedTime,omitempty" bson:"passwordModifiedTime,omitempty" gorm:"column:passwordmodifiedtime"`
	MaxPasswordAge       int        `json:"maxPasswordAge,omitempty" bson:"maxPasswordAge,omitempty" gorm:"column:maxpasswordage"`
	UserType             string     `json:"userType,omitempty" bson:"userType,omitempty" gorm:"column:usertype"`
	Roles                *[]string  `json:"roles,omitempty" bson:"roles,omitempty" gorm:"column:roles"`
	AccessDateFrom       *time.Time `json:"accessDateFrom,omitempty" bson:"accessDateFrom,omitempty" gorm:"type:date;column:accessdatefrom"`
	AccessDateTo         *time.Time `json:"accessDateTo,omitempty" bson:"accessDateTo,omitempty" gorm:"column:accessDateTo"`
	AccessTimeFrom       *time.Time `json:"accessTimeFrom,omitempty" bson:"accessTimeFrom,omitempty" gorm:"type:date;column:accesstimefrom"`
	AccessTimeTo         *time.Time `json:"accessTimeTo,omitempty" bson:"accessTimeTo,omitempty" gorm:"column:accesstimeto"`
}
