package oauth2

type OAuth2Info struct {
	Id             string `mapstructure:"id" json:"id,omitempty" gorm:"column:id" bson:"_id,omitempty" dynamodbav:"id,omitempty" firestore:"id,omitempty"`
	Code           string `mapstructure:"code" json:"code,omitempty" gorm:"column:code" bson:"code,omitempty" dynamodbav:"code,omitempty" firestore:"code,omitempty"`
	RedirectUri    string `mapstructure:"redirect_uri" json:"redirectUri,omitempty" gorm:"column:redirecturi" bson:"redirectUri,omitempty" dynamodbav:"redirectUri,omitempty" firestore:"redirectUri,omitempty"`
	InvitationMail string `mapstructure:"invitation_mail" json:"invitationMail,omitempty" gorm:"column:invitationmail" bson:"invitationMail,omitempty" dynamodbav:"invitationMail,omitempty" firestore:"invitationMail,omitempty"`
	Link           bool   `mapstructure:"link" json:"link,omitempty" gorm:"column:link" bson:"link,omitempty" dynamodbav:"link,omitempty" firestore:"link,omitempty"`
}
