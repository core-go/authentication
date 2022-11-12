package oauth2

type OAuth2Config struct {
	Services string             `mapstructure:"services" json:"services,omitempty" gorm:"column:services" bson:"services,omitempty" dynamodbav:"services,omitempty" firestore:"services,omitempty"`
	Schema   OAuth2SchemaConfig `mapstructure:"schema" json:"schema,omitempty" gorm:"column:schema" bson:"schema,omitempty" dynamodbav:"schema,omitempty" firestore:"schema,omitempty"`
}

type CallbackURL struct {
	Facebook  string `mapstructure:"facebook" json:"facebook,omitempty" gorm:"column:facebook" bson:"facebook,omitempty" dynamodbav:"facebook,omitempty" firestore:"facebook,omitempty"`
	Google    string `mapstructure:"google" json:"google,omitempty" gorm:"column:google" bson:"google,omitempty" dynamodbav:"google,omitempty" firestore:"google,omitempty"`
	LinkedIn  string `mapstructure:"linked_in" json:"linkedIn,omitempty" gorm:"column:linkedin" bson:"linkedIn,omitempty" dynamodbav:"linkedIn,omitempty" firestore:"linkedIn,omitempty"`
	Twitter   string `mapstructure:"twitter" json:"twitter,omitempty" gorm:"column:twitter" bson:"twitter,omitempty" dynamodbav:"twitter,omitempty" firestore:"twitter,omitempty"`
	Microsoft string `mapstructure:"microsoft" json:"microsoft,omitempty" gorm:"column:microsoft" bson:"microsoft,omitempty" dynamodbav:"microsoft,omitempty" firestore:"microsoft,omitempty"`
	Amazon    string `mapstructure:"amazon" json:"amazon,omitempty" gorm:"column:amazon" bson:"amazon,omitempty" dynamodbav:"amazon,omitempty" firestore:"amazon,omitempty"`
	Apple     string `mapstructure:"apple" json:"apple,omitempty" gorm:"column:apple" bson:"apple,omitempty" dynamodbav:"apple,omitempty" firestore:"apple,omitempty"`
	Dropbox   string `mapstructure:"dropbox" json:"dropbox,omitempty" gorm:"column:dropbox" bson:"dropbox,omitempty" dynamodbav:"dropbox,omitempty" firestore:"dropbox,omitempty"`
	Github    string `mapstructure:"github" json:"github,omitempty" gorm:"column:github" bson:"github,omitempty" dynamodbav:"github,omitempty" firestore:"github,omitempty"`
	Gitlab    string `mapstructure:"gitlab" json:"gitlab,omitempty" gorm:"column:gitlab" bson:"gitlab,omitempty" dynamodbav:"gitlab,omitempty" firestore:"gitlab,omitempty"`
	Paypal    string `mapstructure:"paypal" json:"paypal,omitempty" gorm:"column:paypal" bson:"paypal,omitempty" dynamodbav:"paypal,omitempty" firestore:"paypal,omitempty"`
	Instagram string `mapstructure:"instagram" json:"instagram,omitempty" gorm:"column:instagram" bson:"instagram,omitempty" dynamodbav:"instagram,omitempty" firestore:"instagram,omitempty"`
	Kakao     string `mapstructure:"kakao" json:"kakao,omitempty" gorm:"column:kakao" bson:"kakao,omitempty" dynamodbav:"kakao,omitempty" firestore:"kakao,omitempty"`
	Slack     string `mapstructure:"slack" json:"slack,omitempty" gorm:"column:slack" bson:"slack,omitempty" dynamodbav:"slack,omitempty" firestore:"slack,omitempty"`
	Spotify   string `mapstructure:"Spotify" json:"Spotify,omitempty" gorm:"column:Spotify" bson:"Spotify,omitempty" dynamodbav:"Spotify,omitempty" firestore:"Spotify,omitempty"`
	Uber      string `mapstructure:"uber" json:"uber,omitempty" gorm:"column:uber" bson:"uber,omitempty" dynamodbav:"uber,omitempty" firestore:"uber,omitempty"`
	Heroku    string `mapstructure:"heroku" json:"heroku,omitempty" gorm:"column:heroku" bson:"heroku,omitempty" dynamodbav:"heroku,omitempty" firestore:"heroku,omitempty"`
	Asana     string `mapstructure:"asana" json:"asana,omitempty" gorm:"column:asana" bson:"asana,omitempty" dynamodbav:"asana,omitempty" firestore:"asana,omitempty"`
}
