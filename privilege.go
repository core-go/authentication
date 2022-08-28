package auth

type Privilege struct {
	Id          string       `yaml:"id" mapstructure:"id" json:"id,omitempty" gorm:"column:id" bson:"_id,omitempty" dynamodbav:"id,omitempty" firestore:"id,omitempty"`
	Name        string       `yaml:"name" mapstructure:"name" json:"name,omitempty" gorm:"column:name" bson:"name,omitempty" dynamodbav:"name,omitempty" firestore:"name,omitempty"`
	Resource    string       `yaml:"resource" mapstructure:"resource" json:"resource,omitempty" gorm:"column:resource" bson:"resource,omitempty" dynamodbav:"resource,omitempty" firestore:"resource,omitempty"`
	Path        string       `yaml:"path" mapstructure:"path" json:"path,omitempty" gorm:"column:path" bson:"path,omitempty" dynamodbav:"path,omitempty" firestore:"path,omitempty"`
	Icon        string       `yaml:"icon" mapstructure:"icon" json:"icon,omitempty" gorm:"column:icon" bson:"icon,omitempty" dynamodbav:"icon,omitempty" firestore:"icon,omitempty"`
	Permissions int32        `yaml:"permissions" mapstructure:"permissions" json:"permissions,omitempty" gorm:"column:permissions" bson:"permissions,omitempty" dynamodbav:"permissions,omitempty" firestore:"permissions,omitempty"`
	Sequence    int          `yaml:"sequence" mapstructure:"sequence" json:"sequence,omitempty" gorm:"column:sequence" bson:"sequence" dynamodbav:"sequence,omitempty" firestore:"sequence,omitempty"`
	Children    *[]Privilege `yaml:"children" mapstructure:"children" json:"children,omitempty" gorm:"column:children" bson:"children,omitempty" dynamodbav:"children,omitempty" firestore:"children,omitempty"`
}
