package auth

type AuthResult struct {
	Status  int          `mapstructure:"status" json:"status" gorm:"column:status" bson:"status" dynamodbav:"status" firestore:"status"`
	User    *UserAccount `mapstructure:"user" json:"user,omitempty" gorm:"column:user" bson:"user,omitempty" dynamodbav:"user,omitempty" firestore:"user,omitempty"`
	Message string       `mapstructure:"message" json:"message,omitempty" gorm:"column:message" bson:"message,omitempty" dynamodbav:"message,omitempty" firestore:"message,omitempty"`
}
