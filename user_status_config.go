package auth

type UserStatusConfig struct {
	Deactivated string `mapstructure:"deactivated" json:"deactivated,omitempty" gorm:"column:deactivated" bson:"deactivated,omitempty" dynamodbav:"deactivated,omitempty" firestore:"deactivated,omitempty"`
	Disable     string `mapstructure:"disable" json:"disable,omitempty" gorm:"column:disable" bson:"disable,omitempty" dynamodbav:"disable,omitempty" firestore:"disable,omitempty"`
	Suspended   string `mapstructure:"suspended" json:"suspended,omitempty" gorm:"column:suspended" bson:"suspended,omitempty" dynamodbav:"suspended,omitempty" firestore:"suspended,omitempty"`
}
