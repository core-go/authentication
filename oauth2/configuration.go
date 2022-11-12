package oauth2

type Configuration struct {
	Id              string `json:"id,omitempty" gorm:"column:id" bson:"_id,omitempty" dynamodbav:"id,omitempty" firestore:"id,omitempty"`
	Link            string `json:"link,omitempty" gorm:"column:link" bson:"link,omitempty" dynamodbav:"link,omitempty" firestore:"link,omitempty"`
	ClientId        string `json:"clientId,omitempty" gorm:"column:clientid" bson:"clientId,omitempty" dynamodbav:"clientId,omitempty" firestore:"clientId,omitempty"`
	Scope           string `json:"scope,omitempty" gorm:"column:scope" bson:"scope,omitempty" dynamodbav:"scope,omitempty" firestore:"scope,omitempty"`
	RedirectUri     string `json:"redirectUri,omitempty" gorm:"column:redirecturi" bson:"redirectUri,omitempty" dynamodbav:"redirectUri,omitempty" firestore:"redirectUri,omitempty"`
	AccessTokenLink string `json:"accessTokenLink,omitempty" gorm:"column:accesstokenlink" bson:"accessTokenLink,omitempty" dynamodbav:"accessTokenLink,omitempty" firestore:"accessTokenLink,omitempty"`
	ClientSecret    string `json:"clientSecret,omitempty" gorm:"column:clientsecret" bson:"clientSecret,omitempty" dynamodbav:"clientSecret,omitempty" firestore:"clientSecret,omitempty"`
}
