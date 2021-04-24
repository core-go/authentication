package auth

type TokenConfig struct {
	Secret  string `mapstructure:"secret" json:"secret,omitempty" gorm:"column:secret" bson:"secret,omitempty" dynamodbav:"secret,omitempty" firestore:"secret,omitempty"`
	Expires int64  `mapstructure:"expires" json:"expires,omitempty" gorm:"column:expires" bson:"expires,omitempty" dynamodbav:"expires,omitempty" firestore:"expires,omitempty"`
}
