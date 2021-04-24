package auth

import "time"

type AccessTime struct {
	AccessDateFrom  *time.Time `mapstructure:"access_date_from" json:"accessDateFrom,omitempty" gorm:"column:accessdatefrom" bson:"accessDateFrom,omitempty" dynamodbav:"accessDateFrom,omitempty" firestore:"accessDateFrom,omitempty"`
	AccessDateTo    *time.Time `mapstructure:"access_date_to" json:"accessDateTo,omitempty" gorm:"column:accessDateTo" bson:"accessDateTo,omitempty" dynamodbav:"accessDateTo,omitempty" firestore:"accessDateTo,omitempty"`
	AccessTimeFrom  *time.Time `mapstructure:"access_time_from" json:"accessTimeFrom,omitempty" gorm:"column:accesstimefrom" bson:"accessTimeFrom,omitempty" dynamodbav:"accessTimeFrom,omitempty" firestore:"accessTimeFrom,omitempty"`
	AccessTimeTo    *time.Time `mapstructure:"access_time_to" json:"accessTimeTo,omitempty" gorm:"column:accesstimeto" bson:"accessTimeTo,omitempty" dynamodbav:"accessTimeTo,omitempty" firestore:"accessTimeTo,omitempty"`
}
