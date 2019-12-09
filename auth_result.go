package auth

type AuthResult struct {
	Status  AuthStatus   `json:"status,omitempty" bson:"status,omitempty" gorm:"column:status"`
	User    *UserAccount `json:"user,omitempty" bson:"user,omitempty" gorm:"column:user"`
	Message string       `json:"message,omitempty" bson:"message,omitempty" gorm:"column:message"`
}
