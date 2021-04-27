package auth

type AuthConfig struct {
	Secret            string       `mapstructure:"secret" json:"secret,omitempty" gorm:"column:secret" bson:"secret,omitempty" dynamodbav:"secret,omitempty" firestore:"secret,omitempty"`
	Expires           int64        `mapstructure:"expires" json:"expires,omitempty" gorm:"column:expires" bson:"expires,omitempty" dynamodbav:"expires,omitempty" firestore:"expires,omitempty"`
	MaxPasswordFailed int          `mapstructure:"max_password_failed" json:"maxPasswordFailed,omitempty" gorm:"column:maxpasswordfailed" bson:"maxPasswordFailed,omitempty" dynamodbav:"maxPasswordFailed,omitempty" firestore:"maxPasswordFailed,omitempty"`
	LockedMinutes     int          `mapstructure:"locked_minutes" json:"lockedMinutes,omitempty" gorm:"column:lockedminutes" bson:"lockedMinutes,omitempty" dynamodbav:"lockedMinutes,omitempty" firestore:"lockedMinutes,omitempty"`
	MaxPasswordAge    int32        `mapstructure:"max_password_age" json:"maxPasswordAge,omitempty" gorm:"column:maxpasswordage" bson:"maxPasswordAge,omitempty" dynamodbav:"maxPasswordAge,omitempty" firestore:"maxPasswordAge,omitempty"`
	Schema            SchemaConfig `mapstructure:"schema" json:"schema,omitempty" gorm:"column:schema" bson:"schema,omitempty" dynamodbav:"schema,omitempty" firestore:"schema,omitempty"`
}
