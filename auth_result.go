package auth

type AuthResult struct {
	Status  int          `yaml:"status" mapstructure:"status" json:"status" gorm:"column:status" bson:"status" dynamodbav:"status" firestore:"status"`
	User    *UserAccount `yaml:"user" mapstructure:"user" json:"user,omitempty" gorm:"column:user" bson:"user,omitempty" dynamodbav:"user,omitempty" firestore:"user,omitempty"`
	Token   string       `yaml:"token" mapstructure:"token" json:"token,omitempty" gorm:"column:token" bson:"token,omitempty" dynamodbav:"token,omitempty" firestore:"token,omitempty"`
	Message string       `yaml:"message" mapstructure:"message" json:"message,omitempty" gorm:"column:message" bson:"message,omitempty" dynamodbav:"message,omitempty" firestore:"message,omitempty"`
}
