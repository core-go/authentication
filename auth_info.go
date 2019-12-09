package auth

type AuthInfo struct {
	UserName string `json:"userName,omitempty" bson:"userName,omitempty" gorm:"column:username"`
	Password string `json:"password,omitempty" bson:"password,omitempty" gorm:"column:password"`
}
