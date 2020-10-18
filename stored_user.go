package auth

type StoredUser struct {
	UserId     string            `json:"userId,omitempty" gorm:"column:userid" bson:"_id,omitempty" dynamodbav:"userId,omitempty" firestore:"userId,omitempty"`
	Username   string            `json:"username,omitempty" gorm:"column:username" bson:"username,omitempty" dynamodbav:"username,omitempty" firestore:"username,omitempty"`
	Contact    string            `json:"contact,omitempty" gorm:"column:contact" bson:"contact,omitempty" dynamodbav:"contact,omitempty" firestore:"contact,omitempty"`
	UserType   string            `json:"userType,omitempty" gorm:"column:usertype" bson:"userType,omitempty" dynamodbav:"userType,omitempty" firestore:"userType,omitempty"`
	Ip         string            `json:"ip,omitempty" gorm:"column:ip" bson:"ip,omitempty" dynamodbav:"ip,omitempty" firestore:"ip,omitempty"`
	Roles      *[]string         `json:"roles,omitempty" gorm:"column:roles" bson:"roles,omitempty" dynamodbav:"roles,omitempty" firestore:"roles,omitempty"`
	Privileges *[]string         `json:"privileges,omitempty" bson:"privileges,omitempty" gorm:"column:privileges" dynamodbav:"privileges,omitempty" firestore:"privileges,omitempty"`
	Tokens     map[string]string `json:"tokens,omitempty" bson:"tokens,omitempty" gorm:"column:tokens" dynamodbav:"tokens,omitempty" firestore:"tokens,omitempty"`
}
