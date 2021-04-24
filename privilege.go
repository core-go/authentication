package auth

type Privilege struct {
	Id          string       `mapstructure:id"" json:"id,omitempty" gorm:"column:id" bson:"_id,omitempty" dynamodbav:"id,omitempty" firestore:"id,omitempty"`
	Name        string       `mapstructure:"name" json:"name,omitempty" gorm:"column:name" bson:"name,omitempty" dynamodbav:"name,omitempty" firestore:"name,omitempty"`
	Resource    string       `mapstructure:"resource" json:"resource,omitempty" gorm:"column:resource" bson:"resource,omitempty" dynamodbav:"resource,omitempty" firestore:"resource,omitempty"`
	Path        string       `mapstructure:"path" json:"path,omitempty" gorm:"column:path" bson:"path,omitempty" dynamodbav:"path,omitempty" firestore:"path,omitempty"`
	Icon        string       `mapstructure:"icon" json:"icon,omitempty" gorm:"column:icon" bson:"icon,omitempty" dynamodbav:"icon,omitempty" firestore:"icon,omitempty"`
	Permissions int32        `mapstructure:"permissions" json:"permissions,omitempty" gorm:"column:permissions" bson:"permissions,omitempty" dynamodbav:"permissions,omitempty" firestore:"permissions,omitempty"`
	Sequence    int          `mapstructure:"sequence" json:"sequence,omitempty" gorm:"column:sequence" bson:"sequence" dynamodbav:"sequence,omitempty" firestore:"sequence,omitempty"`
	Children    *[]Privilege `mapstructure:"children" json:"children,omitempty" gorm:"column:children" bson:"children,omitempty" dynamodbav:"children,omitempty" firestore:"children,omitempty"`
}
