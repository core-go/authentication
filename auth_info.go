package auth

type AuthInfo struct {
	Step       int    `json:"step,omitempty" gorm:"column:step" bson:"step,omitempty" dynamodbav:"step,omitempty" firestore:"step,omitempty"`
	Username   string `json:"username,omitempty" gorm:"column:username;primary_key" bson:"username,omitempty" dynamodbav:"username,omitempty" firestore:"username,omitempty"`
	Password   string `json:"password,omitempty" gorm:"column:password" bson:"password,omitempty" dynamodbav:"password,omitempty" firestore:"password,omitempty"`
	Passcode   string `json:"passcode,omitempty" gorm:"column:passcode" bson:"passcode,omitempty" dynamodbav:"passcode,omitempty" firestore:"passcode,omitempty"`
	Ip         string `json:"ip,omitempty" gorm:"column:ip" bson:"ip,omitempty" dynamodbav:"ip,omitempty" firestore:"ip,omitempty"`
	Device     string `json:"device,omitempty" gorm:"column:device" bson:"device,omitempty" dynamodbav:"device,omitempty" firestore:"device,omitempty"`
	SenderType string `json:"senderType,omitempty" gorm:"column:sendertype" bson:"senderType,omitempty" dynamodbav:"senderType,omitempty" firestore:"senderType,omitempty"`
}
