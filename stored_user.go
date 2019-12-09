package auth

type StoredUser struct {
	UserId     string    `json:"userId,omitempty" bson:"userId,omitempty"`
	UserName   string    `json:"userName,omitempty" bson:"userName,omitempty"`
	Email      string    `json:"email,omitempty" bson:"email,omitempty"`
	UserType   string    `json:"userType,omitempty" bson:"userType,omitempty"`
	Roles      *[]string `json:"roles,omitempty" bson:"roles,omitempty" gorm:"column:roles"`
	Privileges *[]string `json:"privileges,omitempty" bson:"privileges,omitempty" gorm:"column:privileges"`
}
