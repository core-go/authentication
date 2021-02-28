package auth

import "time"

type UserInfo struct {
	Id                  string     `json:"id,omitempty" gorm:"column:id" bson:"_id,omitempty" dynamodbav:"id,omitempty" firestore:"id,omitempty"`
	Username            string     `json:"username,omitempty" gorm:"column:username" bson:"username,omitempty" dynamodbav:"username,omitempty" firestore:"username,omitempty"`
	Contact             string     `json:"contact,omitempty" gorm:"column:contact" bson:"contact,omitempty" dynamodbav:"contact,omitempty" firestore:"contact,omitempty"`
	Email               string     `json:"email,omitempty" gorm:"column:email" bson:"email,omitempty" dynamodbav:"email,omitempty" firestore:"email,omitempty"`
	Phone               string     `json:"phone,omitempty" gorm:"column:phone" bson:"phone,omitempty" dynamodbav:"phone,omitempty" firestore:"phone,omitempty"`
	DisplayName         string     `json:"displayName,omitempty" gorm:"column:displayname" bson:"displayName,omitempty" dynamodbav:"displayName,omitempty" firestore:"displayName,omitempty"`
	Password            string     `json:"password,omitempty" gorm:"column:password" bson:"password,omitempty" dynamodbav:"password,omitempty" firestore:"password,omitempty"`
	Disable             bool       `json:"disable,omitempty" gorm:"column:disable" bson:"disable,omitempty" dynamodbav:"disable,omitempty" firestore:"disable,omitempty"`
	Deactivated         bool       `json:"deactivated,omitempty" gorm:"column:deactivated" bson:"deactivated,omitempty" dynamodbav:"deactivated,omitempty" firestore:"deactivated,omitempty"`
	Suspended           bool       `json:"suspended,omitempty" gorm:"column:suspended" bson:"suspended,omitempty" dynamodbav:"suspended,omitempty" firestore:"suspended,omitempty"`
	LockedUntilTime     *time.Time `json:"lockedUntilTime,omitempty" gorm:"column:lockeduntiltime" bson:"lockedUntilTime,omitempty" dynamodbav:"lockedUntilTime,omitempty" firestore:"lockedUntilTime,omitempty"`
	SuccessTime         *time.Time `json:"successTime,omitempty" gorm:"column:successtime" bson:"successTime,omitempty" dynamodbav:"successTime,omitempty" firestore:"successTime,omitempty"`
	FailTime            *time.Time `json:"failTime,omitempty" gorm:"column:failtime" bson:"failTime,omitempty" dynamodbav:"failTime,omitempty" firestore:"failTime,omitempty"`
	FailCount           int        `json:"failCount,omitempty" gorm:"column:failcount" bson:"failCount,omitempty" dynamodbav:"failCount,omitempty" firestore:"failCount,omitempty"`
	PasswordChangedTime *time.Time `json:"passwordChangedTime,omitempty" gorm:"column:passwordchangedtime" bson:"passwordChangedTime,omitempty" dynamodbav:"passwordChangedTime,omitempty" firestore:"passwordChangedTime,omitempty"`
	MaxPasswordAge      int        `json:"maxPasswordAge,omitempty" gorm:"column:maxpasswordage" bson:"maxPasswordAge,omitempty" dynamodbav:"maxPasswordAge,omitempty" firestore:"maxPasswordAge,omitempty"`
	UserType            string     `json:"userType,omitempty" gorm:"column:usertype" bson:"userType,omitempty" dynamodbav:"userType,omitempty" firestore:"userType,omitempty"`
	Roles               *[]string  `json:"roles,omitempty" gorm:"column:roles" bson:"roles,omitempty" dynamodbav:"roles,omitempty" firestore:"roles,omitempty"`
	Privileges          *[]string  `json:"privileges,omitempty" gorm:"column:privileges" bson:"privileges,omitempty" dynamodbav:"privileges,omitempty" firestore:"privileges,omitempty"`
	AccessDateFrom      *time.Time `json:"accessDateFrom,omitempty" gorm:"column:accessdatefrom" bson:"accessDateFrom,omitempty" dynamodbav:"accessDateFrom,omitempty" firestore:"accessDateFrom,omitempty"`
	AccessDateTo        *time.Time `json:"accessDateTo,omitempty" gorm:"column:accessDateTo" bson:"accessDateTo,omitempty" dynamodbav:"accessDateTo,omitempty" firestore:"accessDateTo,omitempty"`
	AccessTimeFrom      *time.Time `json:"accessTimeFrom,omitempty" gorm:"column:accesstimefrom" bson:"accessTimeFrom,omitempty" dynamodbav:"accessTimeFrom,omitempty" firestore:"accessTimeFrom,omitempty"`
	AccessTimeTo        *time.Time `json:"accessTimeTo,omitempty" gorm:"column:accesstimeto" bson:"accessTimeTo,omitempty" dynamodbav:"accessTimeTo,omitempty" firestore:"accessTimeTo,omitempty"`
	TwoFactors          bool       `json:"twoFactors,omitempty" gorm:"column:twofactors" bson:"twoFactors,omitempty" dynamodbav:"twoFactors,omitempty" firestore:"twoFactors,omitempty"`
	Language            string     `json:"language,omitempty" gorm:"column:language" bson:"language,omitempty" dynamodbav:"language,omitempty" firestore:"language,omitempty"`
	Gender              string     `json:"gender,omitempty" gorm:"column:gender" bson:"gender,omitempty" dynamodbav:"gender,omitempty" firestore:"gender,omitempty"`
	DateFormat          string     `json:"dateFormat,omitempty" gorm:"column:dateformat" bson:"dateFormat,omitempty" dynamodbav:"dateFormat,omitempty" firestore:"dateFormat,omitempty"`
	TimeFormat          string     `json:"timeFormat,omitempty" gorm:"column:timeformat" bson:"timeFormat,omitempty" dynamodbav:"timeFormat,omitempty" firestore:"timeFormat,omitempty"`
	ImageURL            string     `json:"imageURL,omitempty" gorm:"column:imageurl" bson:"imageURL,omitempty" dynamodbav:"imageURL,omitempty" firestore:"imageURL,omitempty"`
	Status              string     `json:"status,omitempty" gorm:"column:status" bson:"status,omitempty" dynamodbav:"status,omitempty" firestore:"status,omitempty"`
}
