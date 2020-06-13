package auth

type AuthResult struct {
	Status  AuthStatus   `json:"status" gorm:"column:status" bson:"status" dynamodbav:"status" firestore:"status"`
	User    *UserAccount `json:"user,omitempty" gorm:"column:user" bson:"user,omitempty" dynamodbav:"user,omitempty" firestore:"user,omitempty"`
	Message string       `json:"message,omitempty" gorm:"column:message" bson:"message,omitempty" dynamodbav:"message,omitempty" firestore:"message,omitempty"`
}
