package sql

type SchemaConfig struct {
	Id            string `yaml:"id" mapstructure:"id"`
	Username      string `yaml:"username" mapstructure:"username"`
	PrincipalName string `yaml:"principal_name" mapstructure:"principal_name"`
	Status        string `yaml:"status" mapstructure:"status"`

	DisplayName string `yaml:"display_name" mapstructure:"display_name"`
	GivenName   string `yaml:"given_name" mapstructure:"given_name"`
	Surname     string `yaml:"surname" mapstructure:"surname"`

	JobTitle    string `yaml:"job_title" mapstructure:"job_title" json:"jobTitle,omitempty" gorm:"column:jobTitle" bson:"jobTitle,omitempty" dynamodbav:"jobTitle,omitempty" firestore:"jobTitle,omitempty"`
	Language    string `yaml:"language" mapstructure:"language" json:"language,omitempty" gorm:"column:language" bson:"language,omitempty" dynamodbav:"language,omitempty" firestore:"language,omitempty"`

	CreatedTime string `yaml:"created_time" mapstructure:"created_time"`
	CreatedBy   string `yaml:"created_by" mapstructure:"created_by"`
	UpdatedTime string `yaml:"updated_time" mapstructure:"updated_time"`
	UpdatedBy   string `yaml:"updated_by" mapstructure:"updated_by"`
	Version     string `yaml:"version" mapstructure:"version"`
}
