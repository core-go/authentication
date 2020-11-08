package auth

type Privilege struct {
	Id          string       `json:"id,omitempty" gorm:"column:id" bson:"_id,omitempty" dynamodbav:"id,omitempty" firestore:"id,omitempty"`
	Name        string       `json:"name,omitempty" gorm:"column:name" bson:"name,omitempty" dynamodbav:"name,omitempty" firestore:"name,omitempty"`
	Resource    string       `json:"resource,omitempty" gorm:"column:resource" bson:"resource,omitempty" dynamodbav:"resource,omitempty" firestore:"resource,omitempty"`
	Path        string       `json:"path,omitempty" gorm:"column:path" bson:"path,omitempty" dynamodbav:"path,omitempty" firestore:"path,omitempty"`
	Icon        string       `json:"icon,omitempty" gorm:"column:icon" bson:"icon,omitempty" dynamodbav:"icon,omitempty" firestore:"icon,omitempty"`
	Permissions int32        `json:"permissions,omitempty" gorm:"column:permissions" bson:"permissions,omitempty" dynamodbav:"permissions,omitempty" firestore:"permissions,omitempty"`
	Sequence    int          `json:"sequence,omitempty" gorm:"column:sequence" bson:"sequence" dynamodbav:"sequence,omitempty" firestore:"sequence,omitempty"`
	Children    *[]Privilege `json:"children,omitempty" gorm:"column:children" bson:"children,omitempty" dynamodbav:"children,omitempty" firestore:"children,omitempty"`
}
