package auth

type AuthInfo struct {
	Step     int    `mapstructure:"step" json:"step,omitempty" gorm:"column:step" bson:"step,omitempty" dynamodbav:"step,omitempty" firestore:"step,omitempty"`
	Username string `mapstructure:"username" son:"username,omitempty" gorm:"column:username;primary_key" bson:"username,omitempty" dynamodbav:"username,omitempty" firestore:"username,omitempty"`
	Password string `mapstructure:"password" json:"password,omitempty" gorm:"column:password" bson:"password,omitempty" dynamodbav:"password,omitempty" firestore:"password,omitempty"`
	Passcode string `mapstructure:"passcode" json:"passcode,omitempty" gorm:"column:passcode" bson:"passcode,omitempty" dynamodbav:"passcode,omitempty" firestore:"passcode,omitempty"`
	Ip       string `mapstructure:"ip" json:"ip,omitempty" gorm:"column:ip" bson:"ip,omitempty" dynamodbav:"ip,omitempty" firestore:"ip,omitempty"`
	Device   string `mapstructure:"device" json:"device,omitempty" gorm:"column:device" bson:"device,omitempty" dynamodbav:"device,omitempty" firestore:"device,omitempty"`
	Sender   string `mapstructure:"sender" json:"sender,omitempty" gorm:"column:sender" bson:"sender,omitempty" dynamodbav:"sender,omitempty" firestore:"sender,omitempty"`
}
